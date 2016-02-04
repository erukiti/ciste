package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Output struct {
	Output string
}

func boxStatus(w http.ResponseWriter, r *http.Request) {
	a := strings.Split(r.URL.Path, "/")
	rev := a[4]

	filePath := getResultStatusPath(rev)
	statusText, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		// 404 とか適当なの返す
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, `{"error": "not found"}`)
	} else {
		w.Header().Set("Content-Type", "text/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, string(statusText))
	}
}

func boxOutput(w http.ResponseWriter, r *http.Request) {
	a := strings.Split(r.URL.Path, "/")
	rev := a[4]

	filePath := getResultOutputPath(rev)
	outputText, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		// 404 とか適当なの返す
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, `{"error": "resource not found"}`)
	} else {
		w.Header().Set("Content-Type", "text/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		output := Output{string(outputText)}
		jsonData, err := json.Marshal(output)
		if err != nil {
			log.Printf("json failed %v\n", err)
			return
		}

		w.Write(jsonData)
	}
}

func apiStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, "HOGE!")
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	var fn string
	if r.URL.Path == "/" {
		fn = "index.html"
	} else {
		fn = strings.Trim(r.URL.Path, "/")
	}
	path := fmt.Sprintf("ciste-web-content/dist/%s", fn)
	data, err := Asset(path)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `file not found`)
	} else {
		//Won(*3*)chu FixMe!
		// MIME file type?
		w.Write(data)
	}
}

func testvd(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
}

func cisteHttpServer(port int) {
	regexpHandler := CreateRegexpHandler()
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/box/[0-9a-f]{40}/status$"), boxStatus)
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/box/[0-9a-f]{40}/output$"), boxOutput)
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/status$"), apiStatus)

	regexpHandler.HandleFunc(regexp.MustCompile("^/[^/]*$"), staticFiles)

	// regexpHandler.HandleFuncVHost("hoge", testvd)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), regexpHandler)
	if err != nil {
		log.Println(err)
	}

}
