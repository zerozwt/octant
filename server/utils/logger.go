package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/swe"
)

type apiLogRenderer struct{}

func LogRenderer() swe.LogRenderer { return apiLogRenderer{} }

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
