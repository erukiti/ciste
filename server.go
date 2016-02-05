package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func cisteServer(home string, conf Conf, args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	port := fs.Int("port", conf.Port, "httpd port")
	// domain := fs.String("domain", conf.Domain, "FQDN")

	fs.Parse(args)

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

				go ci(conf, box)
			}(fd)
		}
	}()

	cisteHttpServer(*port)

	for {
		time.Sleep(1 * time.Second)
	}
	log.Println("?????")
}
