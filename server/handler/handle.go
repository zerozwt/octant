package handler

import (
	"sync"

	"github.com/zerozwt/swe"
)

var handlerMap map[string][]swe.HandlerFunc
var handlerLock sync.Mutex

const (
	API_PREFIX = "/api"
)

func registerHandler[InType, OutType any](path string, handler func(*swe.Context, *InType) (*OutType, swe.SweError), middlewares ...swe.HandlerFunc) {
	registerRawHandler(path, swe.MakeAPIHandler(handler), middlewares...)
}

func registerRawHandler(path string, handler swe.HandlerFunc, middlewares ...swe.HandlerFunc) {
	handlerLock.Lock()
	defer handlerLock.Unlock()

	fullPath := API_PREFIX + path

	if _, ok := handlerMap[fullPath]; ok {
		return
	}

	handlerMap[fullPath] = append([]swe.HandlerFunc{handler, swe.InitLogID}, middlewares...)
}

func InitAPIServer(server *swe.APIServer) {
	for path, handlers := range handlerMap {
		server.RegisterHandler(path, handlers[0], handlers[1:]...)
	}
}
