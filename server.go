package main

import (
	"flag"
	"fmt"
	"github.com/erukiti/go-util"
	"log"
	// "net/http"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func cisteServer(home string, args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	logFile := fs.String("log", "~/.ciste/log.txt", "log file")

	fs.Parse(args)

	if logFile != nil && *logFile != "" {
		s := util.PathResolvWithMkdirAll(home, *logFile)
		logWriter, err := os.OpenFile(s, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("log file error: %s\n", err)
		} else {
			log.SetOutput(logWriter)
		}
	}

	pid := os.Getpid()

	log.Printf("ciste server start. %d", pid)
	writePidFile()

	go func() {
		socketFile := fmt.Sprintf("/tmp/ciste.%d.sock", pid)
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			log.Fatal("listen error:", err)
		}

		for {
			log.Println("accept.")
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go func(fd net.Conn) {
				// var repository string
				buf := make([]byte, 10240)
				nr, err := fd.Read(buf)
				if err != nil {
					log.Printf("%v\n", err)
				}

				appPath := strings.Trim(string(buf[:nr]), "\n")

				log.Println(appPath)

				go ci(appPath)
			}(fd)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
	log.Println("?????")
}

func ci(appPath string) {
	// _, err := os.Stat(fmt.Sprintf("%s/package.json", appPath))

	a := strings.Split(appPath, "/")
	rev := a[len(a)-1]

	log.Println(rev)

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
		return
	}
	success = execCommand(appPath, "docker", "run", "--rm", "node:local", "bash", "-l", "npm", "test")
	log.Println(success)
}

func execCommand(dir string, args ...string) bool {
	var err error

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stderr.Close()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return false
	}
	st, err := cmd.Process.Wait()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return st.Success()
}
