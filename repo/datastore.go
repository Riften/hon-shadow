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

type Botstore interface {
	Queryable
	AddOrUpdate(key string, value []byte) error
	Get(key string) *pb.BotKV
	Delete(key string) error
}

type FileStore interface {
	Queryable
	Add(file *pb.FileIndex) error
	Get(hash string) *pb.FileIndex
	GetByPrimary(mill string, checksum string) *pb.FileIndex
	GetBySource(mill string, source string, opts string) *pb.FileIndex
	AddTarget(hash string, target string) error
	RemoveTarget(hash string, target string) error
	Count() int
	Delete(hash string) error
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

