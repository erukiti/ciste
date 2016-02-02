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

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logWriter, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("log file error: %s\n", err)
	} else {
		log.SetOutput(logWriter)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		cwd = "/home/git"
	}

	original := os.Getenv("SSH_ORIGINAL_COMMAND")
	sp := strings.Split(original, " ")
	originalCommand := sp[0]
	repo := strings.Trim(sp[1], "'")

	switch os.Args[1] {
	case "ssh":
		// len

		repoPath := fmt.Sprintf("%s/repository/%s", cwd, repo)

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

			code := fmt.Sprintf("#! /bin/sh\ncat | %s receive %s %s %s\n", os.Args[0], os.Args[2], os.Args[3], repo)
			ioutil.WriteFile(hookPath, []byte(code), 0755)
		}

		receiveCommand := fmt.Sprintf("%s 'repository/%s'", originalCommand, repo)
		log.Println(receiveCommand)
		cmd := exec.Command("git-shell", "-c", receiveCommand)
		stdin, err := cmd.StdinPipe()
		defer stdin.Close()
		stdout, err := cmd.StdoutPipe()
		defer stdout.Close()

		go io.Copy(stdin, os.Stdin)
		go io.Copy(os.Stdout, stdout)
		err = cmd.Run()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

	case "receive":
		log.Println("receive start")

		fmt.Println("hoge")

	}

	_ = cwd

	os.Exit(1)
}
