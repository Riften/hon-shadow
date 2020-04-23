package repo

import (
	"database/sql"
	"time"

	//"github.com/Riften/hon-shadow/keypair"
	"github.com/Riften/hon-shadow/pb"
)

type Datastore interface {
	Config() ConfigStore
	Peers() PeerStore
	StreamMetas() StreamMetaStore
	StreamBlocks() StreamBlockStore
	Ping() error
	Close()
}

type Queryable interface {
	BeginTransaction() (*sql.Tx, error)
	PrepareQuery(string) (*sql.Stmt, error)
	PrepareAndExecuteQuery(string, ...interface{}) (*sql.Rows, error)
	ExecuteQuery(string, ...interface{}) (sql.Result, error)
}

type ConfigStore interface {
	Init(pin string) error
	Configure(accnt *keypair.Full, created time.Time) error
	GetAccount() (*keypair.Full, error)
	GetCreationDate() (time.Time, error)
	IsEncrypted() bool
	GetLastDaily() (time.Time, error)
	SetLastDaily() error
}

type PeerStore interface {
	Queryable
	Add(peer *pb.Peer) error
	AddOrUpdate(peer *pb.Peer) error
	Get(id string) *pb.Peer
	GetBestUser(id string) *pb.User
	List(query string) []*pb.Peer
	Find(address string, name string, exclude []string) []*pb.Peer
	Count(query string) int
	UpdateName(id string, name string) error
	UpdateAvatar(id string, avatar string) error
	UpdateInboxes(id string, inboxes []*pb.Cafe) error
	Delete(id string) error
	DeleteByAddress(address string) error
}



type StreamBlockStore interface {
	Queryable
	Add(streamblock *pb.StreamBlock) error
	ListByStream(streamid string, startindex int, maxnum int) []*pb.StreamBlock
	Delete(streamid string) error
	GetByCid(cid string) *pb.StreamBlock
    BlockCount(streamid string) uint64
    LastIndex(streamid string) uint64
}

type StreamMetaStore interface {
	Queryable
	Add(stream *pb.StreamMeta) error
    UpdateNblocks(id string, nblocks uint64) error
	Get(streamId string) *pb.StreamMeta
	Delete(streamId string) error
	List() *pb.StreamMetaList
}

