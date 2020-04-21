package service

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	inet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/Riften/hon-shadow/pb"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"io"
	//"time"
)
// Shadow protocal service
// Make sure to change it to "shadow" instead of "shadow_test" before connect with normal peer
const shadowServiceProtocol = protocol.ID("/textile/shadow_test/1.0.0")

type ShadowService struct {
	ctx context.Context
	peerHost host.Host
	strmap map[peer.ID]*messageSender
}

func NewShadowService (ctx context.Context, phost host.Host) *ShadowService {
	return &ShadowService{
		ctx: ctx,
		peerHost: phost,
	}
}

func (srv *ShadowService) Start() {
	//srv.Node().PeerHost.SetStreamHandler(srv.handler.Protocol(), srv.handleNewStream)
	//go srv.listen("")
	//go srv.listen(srv.Node().Identity.Pretty())
	srv.peerHost.SetStreamHandler(shadowServiceProtocol, srv.handleNewStream)
}

func (srv *ShadowService) Protocol() protocol.ID {
	return shadowServiceProtocol
}

func (srv *ShadowService) handleNewStream(s inet.Stream) {
	defer s.Reset()
	if srv.handleNewMessage(s) {
		// Gracefully close the stream for writes.
		_ = s.Close()
	}
}

func (srv *ShadowService) handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/shadow_service.go Handler: New message receive from %s.\n", pid.Pretty())
	switch env.Message.Type {
	case pb.Message_SHADOW_INFORM:
		return srv.handleInform(env, pid)
	case pb.Message_SHADOW_STREAM_META:
		return srv.handleStreamMeta(env, pid)
	case pb.Message_SHADOW_INFORM_RES:
		return srv.handleInformRes(env, pid)
	default:
		return nil, nil
	}
}

func (srv *ShadowService) handleInform(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	return nil, nil
}

func (srv *ShadowService) handleStreamMeta(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	return nil, nil
}

func (srv *ShadowService) handleInformRes(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	return nil, nil
}

// Core function of service.
// Give the way to handle a message.
func (srv *ShadowService) handleNewMessage(s inet.Stream) bool {
	ctx := srv.ctx

	r := msgio.NewVarintReaderSize(s, inet.MessageSizeMax)

	mPeer := s.Conn().RemotePeer()

	//timer := time.AfterFunc(dhtStreamIdleTimeout, func() { s.Reset() })
	//defer timer.Stop()

	for {
		var req pb.Envelope
		msgbytes, err := r.ReadMsg()
		if err != nil {
			defer r.ReleaseMsg(msgbytes)
			if err == io.EOF {
				return true
			}
			// This string test is necessary because there isn't a single stream reset error
			// instance	in use.
			if err.Error() != "stream reset" {
				//log.Debugf("error reading message: %#v", err)
			}
			return false
		}
		err = proto.Unmarshal(msgbytes, &req)

		// release buffer
		r.ReleaseMsg(msgbytes)

		if err != nil {
			//log.Debugf("error unmarshalling message: %#v", err)
			fmt.Printf("Error unmarshalling message: %#v", err)
			return false
		}

		//timer.Reset(dhtStreamIdleTimeout)

		if err := srv.VerifyEnvelope(&req, mPeer); err != nil {
			//log.Warningf("error verifying message: %s", err)
			continue
		}

		// try a core handler for this msg type


		//log.Debugf("received %s from %s", req.Message.Type.String(), mPeer.Pretty())
		rpmes, err := srv.handle(&req, mPeer)
		if err != nil {
			//log.Warningf("error handling message %s: %s", req.Message.Type.String(), err)
			return false
		}

		//err = srv.updateFromMessage(ctx, mPeer)
		//if err != nil {
			//log.Warningf("error updating from: %s", err)
		//}

		if rpmes == nil {
			continue
		}

		// send out response msg
		//log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

		// send out response msg
		err = writeMsg(s, rpmes)
		if err != nil {
			//log.Debugf("error writing response: %s", err)
			return false
		}
	}
}

// VerifyEnvelope verifies the authenticity of an envelope
func (srv *ShadowService) VerifyEnvelope(env *pb.Envelope, pid peer.ID) error {
	_, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}

	_, err = pid.ExtractPublicKey()
	if err != nil {
		return err
	}

	//return crypto.Verify(pk, ser, env.Sig)
	return nil
}



