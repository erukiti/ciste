package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

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
		fmt.Fprintln(w, `{"error": "not found"}`)
	} else {
		w.Header().Set("Content-Type", "text/json")
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
		fmt.Fprintln(w, `{"error": "resource not found"}`)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, string(outputText))
	}
}

func apiStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r)

	fmt.Fprintln(w, "HOGE!")
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	fn := fmt.Sprintf("ciste-web-content/dist/%s", strings.Trim(r.URL.Path, "/"))
	data, err := Asset(fn)
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

func cisteHttpServer() {
	regexpHandler := CreateRegexpHandler()
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/box/[0-9a-f]+/status$"), boxStatus)
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/box/[0-9a-f]+/output$"), boxOutput)
	regexpHandler.HandleFunc(regexp.MustCompile("^/api/v1/status$"), apiStatus)

	regexpHandler.HandleFunc(regexp.MustCompile("^/[^/]+$"), staticFiles)
	err := http.ListenAndServe(":3000", regexpHandler)
	if err != nil {
		log.Println(err)
	}

}
