/*
Copyright 2016 SASAKI, Shunsuke. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/erukiti/go-util"
	"io/ioutil"
	"log"
	"os"
)

const PidFile = "~/.ciste/pid"

type Conf struct {
	LogFile string
	Domain  string
	Port    int
	home    string
}

func (c *Conf) GetRepositoryPath(home string, repo string) string {
	return util.PathResolvWithMkdirAll(home, fmt.Sprintf("~/.ciste/repository/%s", repo))
}

func (c *Conf) GetAppPath(home string, ref string) string {
	return util.PathResolvWithMkdirAll(home, fmt.Sprintf("~/.ciste/app/%s", ref))
}

func writeConf(home string, conf Conf, confPath string) {
	jsonData, err := json.Marshal(conf)
	if err != nil {
		log.Printf("json failed %v\n", err)
		return
	}

	log.Printf("conf: %s\n", string(jsonData))

	ioutil.WriteFile(util.PathResolv(home, confPath), jsonData, 0644)
}

func readConf(home string, confPath string) Conf {
	jsonData, err := ioutil.ReadFile(util.PathResolv(home, confPath))
	if err != nil {
		// default flag?
		return Conf{"~/.ciste/log.txt", "localhost", 3000, home}
	}

	var conf Conf
	json.Unmarshal(jsonData, &conf)
	conf.home = home

	return conf

}

func printUsage() {
	log.Println(os.Args)
	fmt.Printf("usage: %s <sub command>...\n", os.Args[0])
	fmt.Println("sub command:")
	fmt.Println("  setup")
	fmt.Println("  pubkey")
}

func main() {
	// var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	home := util.GetMyHome()
	if home == "" {
		home = "/home/git"
	}

	if len(os.Args) < 1 {
		printUsage()
		os.Exit(1)
	}

	confPath := flag.String("conf", "~/.ciste/conf.json", "configuration file")
	logFileOverWritten := flag.String("log", "", "log file")

	flag.Parse()

	var logFile string
	conf := readConf(home, *confPath)
	if logFileOverWritten != nil && *logFileOverWritten != "" {
		logFile = *logFileOverWritten
	} else {
		logFile = conf.LogFile
	}

	logWriter, err := os.OpenFile(util.PathResolvWithMkdirAll(home, logFile), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("log file error: %s\n", err)
	} else {
		log.SetOutput(logWriter)
	}

	log.Println(os.Args)
	args := flag.Args()
	switch args[0] {
	case "ssh":
		cisteSSH(home, conf, args[1:])

	case "receive":
		cisteReceive(home, conf, args[1:])

	case "server":
		cisteServer(home, conf, args[1:])

	case "setup":
		cisteSetup(home, conf, args[1:])

	case "pubkey":
		cistePubkey(home, conf, args[1:])
	default:
		printUsage()
		os.Exit(0)
	}
}
