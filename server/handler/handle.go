package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

var handlerMap map[string]map[string][]swe.HandlerFunc = make(map[string]map[string][]swe.HandlerFunc)
var handlerLock sync.Mutex

const (
	API_PREFIX = "/api"
	GET        = http.MethodGet
	POST       = http.MethodPost
)

func registerHandler[InType, OutType any](method, path string, handler func(*swe.Context, *InType) (*OutType, swe.SweError), middlewares ...swe.HandlerFunc) {
	registerRawHandler(method, path, swe.MakeAPIHandler(handler), middlewares...)
}

func registerRawHandler(method, path string, handler swe.HandlerFunc, middlewares ...swe.HandlerFunc) {
	handlerLock.Lock()
	defer handlerLock.Unlock()

	fullPath := API_PREFIX + path

	if _, ok := handlerMap[method]; !ok {
		handlerMap[method] = make(map[string][]swe.HandlerFunc)
	}

	if _, ok := handlerMap[method][fullPath]; ok {
		return
	}

	handlerMap[method][fullPath] = append([]swe.HandlerFunc{handler, swe.InitLogID, setLogRenderer, requestBill}, middlewares...)
}

func InitAPIServer(server *swe.APIServer) {
	for method, set := range handlerMap {
		for path, handlers := range set {
			server.RegisterHandler(method, path, handlers[0], handlers[1:]...)
		}
	}
}

func requestBill(ctx *swe.Context) {
	now := time.Now()
	logger := swe.CtxLogger(ctx)
	defer func() {
		logger.Info("request path: %s process time %v", ctx.Request.URL.Path, time.Since(now))
	}()
	ctx.Next()
}

func setLogRenderer(ctx *swe.Context) {
	swe.CtxLogger(ctx).SetRenderer(utils.LogRenderer())
	ctx.Next()
}
