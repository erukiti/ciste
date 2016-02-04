package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func ci(box *Box) {
	// _, err := os.Stat(fmt.Sprintf("%s/package.json", appPath))

	dockerfile := []byte(`
FROM erukiti/ndenv:base-wheezy

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
