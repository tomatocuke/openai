package bootstrap

import (
	"fmt"
	"net/http"
	"openai/internal/config"
	"strings"
)

var (
	_ http.Handler = (*Engine)(nil)
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodGet, pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPost, pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	prefixMatched := strings.Index(req.URL.Path, config.Http.Prefix) == 0
	path := strings.Replace(req.URL.Path, config.Http.Prefix, "/", 1)
	key := req.Method + "-" + path
	handler, ok := engine.router[key]

	if prefixMatched && ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
