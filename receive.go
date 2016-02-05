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
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getNewRefs() string {
	buf := make([]byte, 1024)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if n == 0 {
		log.Println("failed: get refs")
		os.Exit(1)
	}

	ar := strings.Split(string(buf), " ")
	if len(ar) < 3 {
		log.Println("illegal refs")
		log.Println(ar)
		os.Exit(1)
	}

	return ar[1]
}

func cisteReceive(home string, conf Conf, args []string) {
	log.Printf("receive start. %v", args)

	newrefs := getNewRefs()
	log.Printf("[debug] new refs: ", newrefs)

	appPath := conf.GetAppPath(home, newrefs)

	os.MkdirAll(appPath, 0755)

	repoPath := conf.GetRepositoryPath(home, args[2])

	copyCommand := fmt.Sprintf("cd %s ; (cd %s ; git archive %s) | tar xvf -", appPath, repoPath, newrefs)
	log.Println(copyCommand)
	out, err := exec.Command("/bin/sh", "-c", copyCommand).Output()
	if err != nil {
		fmt.Println(out)
		fmt.Println(err)
		log.Println(out)
		log.Println(err)
		os.Exit(1)
	}

	box := CreateBox(newrefs, args[0], args[2])

	dispatch(*box)
}
