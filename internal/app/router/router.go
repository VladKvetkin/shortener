package router

import "net/http"

type Router struct {
	Routes map[string]http.Handler
}
