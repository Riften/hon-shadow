package core

import (
	"fmt"
	"github.com/Riften/hon-shadow/repo"
	"github.com/Riften/hon-shadow/repo/db"
	logging "github.com/ipfs/go-log"
	"strings"
	"sync"
)

var log = logging.Logger("tex-core")

// InitConfig is used to setup a textile node
type InitConfig struct {
	RepoPath        string
	//SwarmPorts      string
}


// RunConfig is used to define run options for a textile node
type RunConfig struct {
	RepoPath          string
}

type Variables struct {
	SwarmAddress    string
	FailedAddresses []string
	lock            sync.Mutex
	//StreamFileChannels map[string]chan *pb.StreamFile
	//streamBlockIndex map[string]uint64
}


// common errors
var ErrAccountRequired = fmt.Errorf("account required")
var ErrStarted = fmt.Errorf("node is started")
var ErrStopped = fmt.Errorf("node is stopped")
var ErrOffline = fmt.Errorf("node is offline")
var ErrMissingRepoConfig = fmt.Errorf("you must specify InitConfig.RepoPath or InitConfig.BaseRepoPath and InitConfig.Account")

// Repo returns the actual location of the configured repo
func (conf InitConfig) Repo() (string, error) {
	if len(conf.RepoPath) > 0 {
		return conf.RepoPath, nil
	} else {
		return "", ErrMissingRepoConfig
	}
}

// InitRepo initializes a new node repo.
// It does the following things:
//		- Create repo directory.
//		- Create datastore and save it to directory.
func InitRepo(conf InitConfig) error {
	repoPath, err := conf.Repo()
	if err != nil {
		return err
	}

	// init repo
	err = repo.Init(repoPath)

	if err != nil {
		return err
	}

	log.Debug("create db")
	sqliteDb, err := db.Create(repoPath, "")
	if err != nil {
		return err
	}
	log.Debug("init db")
	//err = sqliteDb.Config().Init("")
	if err != nil {
		log.Error(err)
		return err
	}


	log.Debug("finish")
	//return applyTextileConfigOptions(conf)
	return nil
}


type loggingWaitGroup struct {
	n  string
	wg sync.WaitGroup
}

func (lwg *loggingWaitGroup) Add(delta int, src string) {
	log.Debugf("%s wait added delta %d (src=%s)", lwg.n, delta, src)
	lwg.wg.Add(delta)
}

func (lwg *loggingWaitGroup) Done(src string) {
	log.Debugf("%s wait done (src=%s)", lwg.n, src)
	lwg.wg.Done()
}

func (lwg *loggingWaitGroup) Wait(src string) {
	log.Debugf("%s waiting (src=%s)", lwg.n, src)
	lwg.wg.Wait()
}

func getTarget(output string) string {
	str := strings.Split(output, " ")[1];
	return str;
}

func getStatus(output string) bool {
	str := strings.Split(output, " ")[2];
	if str=="success"{
		return true
	}else {
		return false
	}
}

func getId(output string) string {
	list := strings.Split(output, "/")
	return list[len(list)-1]
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}



