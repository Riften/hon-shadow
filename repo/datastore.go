package repo

import (
	"database/sql"
	//"time"

	//"github.com/Riften/hon-shadow/keypair"
	"github.com/Riften/hon-shadow/pb"
)

type Datastore interface {
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

