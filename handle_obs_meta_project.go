package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"math/rand"
	"time"
)

func (bot *robot) createOBSMetaProject(repo string, log *logrus.Entry) {
	if !bot.cfg.EnableCreatingOBSMetaProject {
		return
	}

	project := &bot.cfg.OBSMetaProject
	path := project.genProjectFilePath(repo)
	b := &project.Branch

	// get pid by org and repo field
	gid, err := bot.getGroupID(b.Org)
	if err != nil {
		return
	}

	pid, _, err := bot.getProjectID(gid, b.Repo)
	if err != nil {
		return
	}

	// file exists
	if _, err := bot.cli.GetPathContent(pid, path, b.Branch); err == nil {
		return
	}

	content, err := project.genProjectFileContent(repo)
	if err != nil {
		log.Errorf("generate file of project:%s, err:%s", repo, err.Error())
		return
	}

	w := &bot.cfg.WatchingFiles
	msg := fmt.Sprintf(
		"add project according to the file: %s/%s/%s:%s",
		w.Org, w.Repo, w.Branch, w.RepoOrg,
	)

	err = bot.cli.CreateFile(pid, path, gitlab.CreateFileOptions{
		Branch:        &b.Branch,
		Content:       &content,
		CommitMessage: &msg,
	})
	if err != nil {
		for i := 0; i < 2; i++ {

			rand.Seed(time.Now().UnixNano())
			number := rand.Intn(10000)
			time.Sleep(time.Duration(number) * time.Millisecond)
			err = bot.cli.CreateFile(pid, path, gitlab.CreateFileOptions{
				Branch:        &b.Branch,
				Content:       &content,
				CommitMessage: &msg,
			})
			if err == nil {
				break
			}

			log.Errorf("create file: %s, err:%s", path, err.Error())
		}
	}
}
