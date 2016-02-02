package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
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

	current, err := user.Current()
	var home string
	if err != nil {
		log.Printf("failed user.Current: %s", err)
		home = "/home/git"
	} else {
		home = current.HomeDir
	}

	original := os.Getenv("SSH_ORIGINAL_COMMAND")
	sp := strings.Split(original, " ")
	originalCommand := sp[0]
	repo := strings.Trim(sp[1], "'")

	switch os.Args[1] {
	case "ssh":
		// len

		repoPath := fmt.Sprintf("%s/repository/%s", home, repo)

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

		fmt.Printf("%v\n", os.Args)

		buf := make([]byte, 1024)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n == 0 {
			fmt.Println("failed: get refs")
			os.Exit(1)
		}

		ar := strings.Split(string(buf), " ")
		if len(ar) < 3 {
			fmt.Println("illegal refs")
			fmt.Println(ar)
			os.Exit(1)
		}

		newrefs := ar[1]
		// ar = strings.Split(ar[2], "/")
		fmt.Println(newrefs)
		appPath := fmt.Sprintf("%s/app/%s", home, newrefs)
		fmt.Println(appPath)
		os.MkdirAll(appPath, 0755)
		repoPath := fmt.Sprintf("%s/repository/%s", home, os.Args[4])
		copyCommand := fmt.Sprintf("cd %s ; (cd %s ; git archive %s) | tar xvf -", appPath, repoPath, newrefs)
		fmt.Println(copyCommand)
		out, err := exec.Command("/bin/sh", "-c", copyCommand).Output()
		if err != nil {
			fmt.Println(out)
			fmt.Println(err)
			os.Exit(1)
		}

		// _, err := os.Stat(fmt.Sprintf("%s/package.json", appPath))

		dockerfile := []byte(`
FROM ndenv:base-wheezy

RUN mkdir /app
ADD * /app/
WORKDIR /app

RUN bash -l -c "ndenv install"
RUN bash -l -c "npm install"
CMD bash -l -c "cd /app ; npm start"
`)

		ioutil.WriteFile(fmt.Sprintf("%s/Dockerfile", appPath), dockerfile, 0644)

		var success bool

		success = execCommand(appPath, "docker", "build", "-t", "node:local", ".")
		if !success {
			os.Exit(1)
		}
		success = execCommand(appPath, "docker", "run", "--rm", "node:local", "bash", "-l", "npm", "test")
		fmt.Println(success)

	}

	_ = home

	os.Exit(1)
}

func execCommand(dir string, args ...string) bool {
	var err error

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stderr.Close()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	st, err := cmd.Process.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return st.Success()
}
