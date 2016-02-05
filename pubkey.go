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
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/erukiti/go-util"
	"log"
	"os"
)

func cistePubkey(home string, conf Conf, args []string) {
	if len(args) != 4 {
		fmt.Printf("usage: %s pubkey <user> <public key>\n", os.Args[0])
		os.Exit(1)
	}

	os.Mkdir(util.PathResolv(home, "~/.ssh"), 0700)
	writer, err := os.OpenFile(util.PathResolv(home, "~/.ssh/authorized_keys"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Println("authorize_keys wrote failed. %s\n", err)
		return
	}
	_ = writer

	data, err := base64.StdEncoding.DecodeString(args[2])
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	md5sum := fmt.Sprintf("%x", md5.Sum(data))
	fingerPrint := ""
	for len(md5sum) > 2 {
		fingerPrint += md5sum[0:2] + ":"
		md5sum = md5sum[2:]
	}
	fingerPrint += md5sum

	fmt.Fprintf(writer, "command=\"%s ssh %s %s\",no-agent-forwarding,no-pty,no-user-rc,no-X11-forwarding,no-port-forwarding %s %s %s\n", os.Args[0], args[0], fingerPrint, args[1], args[2], args[3])
}
