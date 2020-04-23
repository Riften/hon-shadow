package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/protocol"
	"sync"
	"time"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/helpers"
	inet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-msgio"
	"github.com/Riften/hon-shadow/pb"
	ggio "github.com/gogo/protobuf/io"
	p2phost "github.com/libp2p/go-libp2p-core/host"
)

func messageSenderForPeer(ctx context.Context, p peer.ID, srv Service, h p2phost.Host) (*messageSender, error) {
	//srv.smlk.Lock()
	//ms, ok := srv.strmap[p]
	ms, ok := srv.GetSender(p)
	if ok {
		//srv.smlk.Unlock()
		return ms, nil
	}

	// Create a new messageSender
	ms = &messageSender{
		p: p,
		protocol: srv.Protocol(),
		h:h,
	}
	srv.AddSender(ms)

	if err := ms.prepOrInvalidate(ctx); err != nil {
		// ms has been invalidated.
		//srv.smlk.Lock()
		//defer srv.smlk.Unlock()

		//if msCur, ok := srv.strmap[p]; ok {
			// Changed. Use the new one, old one is invalid and
			// not in the map so we can just throw it away.
		//	if ms != msCur {
		//		return msCur, nil
		//	}
			// Not changed, remove the now invalid stream from the
			// map.
		//	delete(srv.strmap, p)
		//}
		// Invalid but not in map. Must have been removed by a disconnect.
		srv.RemoveSender(p)
		return nil, err
	}
	// All ready to go.
	return ms, nil
}

type messageSender struct {
	s  inet.Stream
	r  msgio.ReadCloser
	lk sync.Mutex
	p  peer.ID

	//srv *Service
	protocol protocol.ID
	h p2phost.Host

	invalid   bool
	singleMes int
}

// invalidate is called before this messageSender is removed from the strmap.
// It prevents the messageSender from being reused/reinitialized and then
// forgotten (leaving the stream open).
func (ms *messageSender) invalidate() {
	ms.invalid = true
	if ms.s != nil {
		_ = ms.s.Reset()
		ms.s = nil
	}
}

func (ms *messageSender) prepOrInvalidate(ctx context.Context) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	if err := ms.prep(ctx); err != nil {
		ms.invalidate()
		return err
	}
	return nil
}

func (ms *messageSender) prep(ctx context.Context) error {
	if ms.invalid {
		return fmt.Errorf("message sender has been invalidated")
	}
	if ms.s != nil {
		return nil
	}

	//nstr, err := ms.srv.Node().PeerHost.NewStream(ctx, ms.p, ms.srv.handler.Protocol())
	nstr, err := ms.h.NewStream(ctx, ms.p, ms.protocol)
	if err != nil {
		return err
	}

	ms.r = msgio.NewVarintReaderSize(nstr, inet.MessageSizeMax)
	ms.s = nstr

	return nil
}

// streamReuseTries is the number of times we will try to reuse a stream to a
// given peer before giving up and reverting to the old one-message-per-stream
// behaviour.
const streamReuseTries = 3

func (ms *messageSender) SendMessage(ctx context.Context, pmes *pb.Envelope) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(ctx); err != nil {
			return err
		}

		if err := ms.writeMsg(pmes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				//log.Info("error writing message, bailing: ", err)
				return err
			}
			//log.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return nil
	}
}

func (ms *messageSender) SendRequest(ctx context.Context, pmes *pb.Envelope) (*pb.Envelope, error) {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(ctx); err != nil {
			return nil, err
		}

		if err := ms.writeMsg(pmes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				//log.Info("error writing message, bailing: ", err)
				return nil, err
			}
			//log.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		mes := new(pb.Envelope)
		if err := ms.ctxReadMsg(ctx, mes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				//log.Info("error reading message, bailing: ", err)
				return nil, err
			}
			//log.Info("error reading message, trying again: ", err)
			retry = true
			continue
		}

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return mes, nil
	}
}

func (ms *messageSender) writeMsg(pmes *pb.Envelope) error {
	return writeMsg(ms.s, pmes)
}

func (ms *messageSender) ctxReadMsg(ctx context.Context, mes *pb.Envelope) error {
	errc := make(chan error, 1)
	go func(r msgio.ReadCloser) {
		bytes, err := r.ReadMsg()
		defer r.ReleaseMsg(bytes)
		if err != nil {
			errc <- err
			return
		}
		errc <- proto.Unmarshal(bytes, mes)
	}(ms.r)

	t := time.NewTimer(dhtReadMessageTimeout)
	defer t.Stop()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return ErrReadTimeout
	}
}

type bufferedDelimitedWriter struct {
	*bufio.Writer
	ggio.WriteCloser
}

var writerPool = sync.Pool{
	New: func() interface{} {
		w := bufio.NewWriter(nil)
		return &bufferedDelimitedWriter{
			Writer:      w,
			WriteCloser: ggio.NewDelimitedWriter(w),
		}
	},
}

func writeMsg(w io.Writer, mes *pb.Envelope) error {
	bw := writerPool.Get().(*bufferedDelimitedWriter)
	bw.Reset(w)
	err := bw.WriteMsg(mes)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	writerPool.Put(bw)
	return err
}
