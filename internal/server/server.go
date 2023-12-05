package server

import "net/http"

type httpServer struct {
	addr    string
	handler http.Handler
}

func New(addr string, handler http.Handler) *httpServer {
	return &httpServer{
		addr:    addr,
		handler: handler,
	}
}

func (s *httpServer) Run() error {
	return http.ListenAndServe(s.addr, s.handler)
}
