package handler

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zerozwt/octant/server/session"
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

	handlerMap[fullPath] = append([]swe.HandlerFunc{handler, swe.InitLogID, setLogRenderer}, middlewares...)
}

func InitAPIServer(server *swe.APIServer) {
	for path, handlers := range handlerMap {
		server.RegisterHandler(path, handlers[0], handlers[1:]...)
	}
}

// -----------------------------------------------------------------

type apiLogRenderer struct{}

var apiLog apiLogRenderer

func setLogRenderer(ctx *swe.Context) {
	swe.CtxLogger(ctx).SetRenderer(apiLog)
	ctx.Next()
}

func (r apiLogRenderer) RenderLog(ctx *swe.Context, level swe.LogLevel, ts time.Time, file string, line int, content string) string {
	builder := strings.Builder{}
	builder.WriteByte('[')
	builder.WriteString(level.String())
	builder.WriteByte(']')

	builder.WriteByte('[')
	builder.WriteString(swe.RenderTime(ts))
	builder.WriteByte(']')

	builder.WriteByte('[')
	builder.WriteString(filepath.Base(file))
	builder.WriteByte(':')
	builder.WriteString(strconv.Itoa(line))
	builder.WriteByte(']')

	builder.WriteByte('[')
	builder.WriteString(swe.CtxLogID(ctx))
	builder.WriteByte(']')

	if session.IsAdmin(ctx) {
		builder.Write([]byte(`[ADMIN]`))
	}
	if info, ok := session.GetStreamerSession(ctx); ok {
		builder.Write([]byte(`[USER:`))
		builder.WriteString(info.AccountName)
		builder.WriteByte(']')
	}
	if info, ok := session.GetDDSession(ctx); ok {
		builder.Write([]byte(`[DD:`))
		builder.WriteString(fmt.Sprint(info.UID))
		builder.WriteByte(']')
	}

	builder.WriteByte(' ')
	builder.WriteString(content)
	builder.WriteByte('\n')
	return builder.String()
}
