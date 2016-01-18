package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func execShellScript(script string) ([]byte, error) {
	f, err := ioutil.TempFile(os.TempDir(), "ciste")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	f.Write([]byte(script))

	return exec.Command("sh", f.Name()).CombinedOutput()
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

func extractTar(src []byte, dest string) ([]byte, error) {
	f, err := ioutil.TempFile(os.TempDir(), "ciste")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(src)

	return execShellScript(fmt.Sprintf("cd %s && tar zxvf %s", dest, f.Name()))
}

func install(name string, url string, hash string, dest string) error {
	var err error
	log.Printf("%s installing... ", name)

	if isInstalled(name) {
		log.Println("already installed.")
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		log.Println("failed.")
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println("failed.")
		return err
	}
	if sha256Sum(data) != hash {
		log.Println("failed.")
		return fmt.Errorf("download file sha256 illegal.")
	}

	_, err = extractTar(data, dest)
	if err != nil {
		log.Println("failed.")
		return err
	} else {
		log.Println("succeeded.")
		return nil
	}
}

func main() {
	err := install("docker", "https://get.docker.com/builds/Linux/x86_64/docker-1.9.1.tgz", "6a095ccfd095b1283420563bd315263fa40015f1cee265de023efef144c7e52d", "/")
	if err != nil {
		log.Fatal(err)
	}
}
