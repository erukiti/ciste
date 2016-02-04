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
	"flag"
	"github.com/erukiti/go-util"
	"log"
	"os"
)

func cisteSetup(home string, args []string) {
	fs := flag.NewFlagSet("setup", flag.ExitOnError)

	domain := fs.String("domain", "localhost", "server domain")
	port := fs.Int("port", 3000, "ciste web server port")
	logFile := fs.String("log", "~/.ciste/log.txt", "log file")
	confFile := fs.String("conf", "~/.ciste/conf.json", "configutaion file")

	fs.Parse(args)

	if domain == nil || port == nil || logFile == nil || confFile == nil {
		fs.PrintDefaults()
		os.Exit(1)
	}

	log.Println("setup start.")

	os.MkdirAll(util.PathResolv(home, "~/.ciste"), 0755)
	os.MkdirAll(util.PathResolv(home, "~/.ciste/app"), 0755)
	os.MkdirAll(util.PathResolv(home, "~/.ciste/repository"), 0755)

	conf := Conf{*logFile, *domain, *port}
	writeConf(home, conf, *confFile)
}
