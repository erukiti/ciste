package main

import (
	"fmt"
	"github.com/erukiti/go-util"
	"os"
)

type CIStatus struct {
	Success bool
}

type Box struct {
	Rev        string
	User       string
	Repository string
	Status     CIStatus
}

func (b *Box) GetAppDir() string {
	return util.PathResolv("", fmt.Sprintf("~/app/%s", b.Rev))
}

func (b *Box) GetImageId() string {
	return fmt.Sprintf("app:%s", b.Rev)
}

func (b *Box) GetResultDir() string {
	path := util.PathResolvWithMkdirAll("/", fmt.Sprintf("~/.ciste/result/%s", b.Rev))
	os.MkdirAll(path, 0755)

	return path
}

func (b *Box) GetResultOutputPath() string {
	return fmt.Sprintf("%s/output.txt", b.GetResultDir())
}

func (b *Box) GetResultStatusPath() string {
	return fmt.Sprintf("%s/status.json", b.GetResultDir())
}

func CreateBox(rev string, user string, repository string) *Box {
	return &Box{rev, user, repository, CIStatus{}}
}
