package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/erukiti/go-util"
	"log"
	"net"
	// "net/http"
	"os"
	"time"
)

func cisteServer(home string, args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	logFile := fs.String("log", "~/.ciste/log.txt", "log file")
	httpdPort := fs.Int("port", 3000, "httpd port")

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

	cisteHttpServer(*httpdPort)

	for {
		time.Sleep(1 * time.Second)
	}
	log.Println("?????")
}
