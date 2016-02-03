// http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type routeVhost struct {
	vhost   string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
	vhosts []*routeVhost
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) HandleFuncVHost(vhost string, handler func(http.ResponseWriter, *http.Request)) {
	h.vhosts = append(h.vhosts, &routeVhost{vhost, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.Host)

	a := strings.Split(r.Host, ".")
	if len(a) > 1 {
		for _, vhost := range h.vhosts {
			if a[0] == vhost.vhost {
				vhost.handler.ServeHTTP(w, r)
				return
			}
		}
	}

	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func CreateRegexpHandler() *RegexpHandler {
	return &RegexpHandler{}
}
