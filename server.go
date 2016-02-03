package main

import (
	"flag"
	"fmt"
	"github.com/erukiti/go-util"
	"log"
	// "net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
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

				var box *Box
				json.Unmarshal(buf[:nr], &box)

				log.Printf("%v\n", box)

				go ci(box)
			}(fd)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
	log.Println("?????")
}

func ci(box *Box) {
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

	appPath := box.GetAppDir()
	ioutil.WriteFile(fmt.Sprintf("%s/Dockerfile", appPath), dockerfile, 0644)

	var success bool

	imageId := box.GetImageId()
	success = execCommand(box, "docker", "build", "-t", imageId, ".")
	if !success {
		return
	}

	box.Status.Success = execCommand(box, "docker", "run", "--rm", imageId, "bash", "-l", "npm", "test")

	jsonData, err := json.Marshal(box.Status)
	if err != nil {
		log.Printf("json failed %v\n", err)
		return
	}

	log.Printf("[debug] %s", string(jsonData))

	ioutil.WriteFile(box.GetResultStatusPath(), jsonData, 0644)

}

func execCommand(box *Box, args ...string) bool {
	var err error

	appPath := box.GetAppDir()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = appPath

	output, err := os.OpenFile(box.GetResultOutputPath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return false
	}
	defer output.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	defer stderr.Close()

	go io.Copy(output, stdout)
	go io.Copy(output, stderr)

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return false
	}
	st, err := cmd.Process.Wait()
	if err != nil {
		log.Println(err)
		return false
	}

	return st.Success()
}
