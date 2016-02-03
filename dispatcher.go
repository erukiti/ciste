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
	"github.com/erukiti/go-util"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

func readPidFile() (int, error) {
	content, err := ioutil.ReadFile(util.PathResolv("/", PidFile))
	if err != nil {
		return -1, err
	} else {
		pid, err := strconv.Atoi(string(content))
		if err != nil {
			return -1, err
		} else {
			return pid, nil
		}
	}
}

func writePidFile() error {
	pidFile := util.PathResolvWithMkdirAll("/", PidFile)
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

func removePidFile() error {
	log.Println("PID file remove.")
	return os.Remove(util.PathResolv("/", PidFile))
}

func getSocketName(pid int) string {
	return fmt.Sprintf("/tmp/ciste.%d.sock", pid)
}

func dispatch(appPath string) {
	var err error

	log.Println("dispatch")

	pid, err := readPidFile()
	log.Printf("[debug] PID: %v", pid)

	if err != nil {
		// think hospitality
		log.Println(err)
		return
	}

	socketFile := getSocketName(pid)

	c, err := net.Dial("unix", socketFile)
	log.Printf("[debug] socket %v\n", c)
	if err != nil {
		log.Printf("%s: %v\n", socketFile, err)
		return
	}

	defer c.Close()

	_, err = c.Write([]byte(fmt.Sprintf("%s\n", appPath)))
	if err != nil {
		log.Printf("socket write error: %v\n", err)
	} else {
		log.Println("sent")
	}
}
