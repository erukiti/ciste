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

func getAppDir(rev string) string {
	return util.PathResolv("", fmt.Sprintf("~/.ciste/app/%s", rev))
}

func getImageId(rev string) string {
	return fmt.Sprintf("app:%s", rev)
}

func getResultDir(rev string) string {
	path := util.PathResolv("/", fmt.Sprintf("~/.ciste/result/%s", rev))
	os.MkdirAll(path, 0755)

	return path
}

func getResultOutputPath(rev string) string {
	return fmt.Sprintf("%s/output.txt", getResultDir(rev))
}

func getResultStatusPath(rev string) string {
	return fmt.Sprintf("%s/status.json", getResultDir(rev))
}

func (b *Box) GetAppDir() string {
	return getAppDir(b.Rev)
}

func (b *Box) GetImageId() string {
	return getImageId(b.Rev)
}

func (b *Box) GetResultDir() string {
	return getResultDir(b.Rev)
}

func (b *Box) GetResultOutputPath() string {
	return getResultOutputPath(b.Rev)
}

func (b *Box) GetResultStatusPath() string {
	return getResultStatusPath(b.Rev)
}

func CreateBox(rev string, user string, repository string) *Box {
	return &Box{rev, user, repository, CIStatus{}}
}
