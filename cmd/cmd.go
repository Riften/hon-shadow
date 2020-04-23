package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/Riften/hon-shadow/core"
	"os"
)

// Package cmd define the command line instructions

type cmdsMap map[string]func() error	// Store which func for each cmd

var appCmd = kingpin.New("shadow", "shadow is a simplified libp2p node for hon-textile shadow service running on switch")

// Run() start the software
func Run() error {
	cmds := make(cmdsMap)
	appCmd.UsageTemplate(kingpin.CompactUsageTemplate)

	initCmd :=appCmd.Command("init", "Initialize the shadow peer.")
	initCmdRepopath := initCmd.Arg("repo", "Repo path.").Required().String()
	cmds[initCmd.FullCommand()] = func() error {
		cfg := core.InitConfig{RepoPath:*initCmdRepopath}
		return core.InitRepo(cfg)
	}

	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}

	return nil
}
