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
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func cisteSSH(home string, conf Conf, args []string) {
	log.Printf("ssh start. %v", args)

	original := os.Getenv("SSH_ORIGINAL_COMMAND")
	sp := strings.Split(original, " ")
	originalCommand := sp[0]
	repo := strings.Trim(sp[1], "'")

	repoPath := conf.GetRepositoryPath(home, repo)

	_, err := os.Stat(repoPath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(repoPath, 0755)
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = repoPath
		out, err := cmd.Output()
		if err != nil {
			log.Println(string(out))
			log.Println(err)
			os.Exit(1)
		}
		log.Println(string(out))

		hookPath := fmt.Sprintf("%s/hooks/pre-receive", repoPath)
		log.Println(hookPath)

		code := fmt.Sprintf("#! /bin/sh\n%s receive %s %s %s\n", os.Args[0], args[0], args[1], repo)
		ioutil.WriteFile(hookPath, []byte(code), 0755)
	}

	receiveCommand := fmt.Sprintf("%s '.ciste/repository/%s'", originalCommand, repo)
	log.Println(receiveCommand)
	cmd := exec.Command("git-shell", "-c", receiveCommand)

	stdin, err := cmd.StdinPipe()
	defer stdin.Close()
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	go io.Copy(stdin, os.Stdin)
	go io.Copy(os.Stdout, stdout)
	// FIXME: remove "remote: "

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
