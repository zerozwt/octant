package session

import (
	"net/http"
	"strings"
	"time"

	"github.com/zerozwt/swe"
)

type AdminSession struct{}

type StreamerSession struct{}

type DDSession struct{}

const (
	ctxAdminKey    = "ctx_o_admin"
	ctxStreamerKey = "ctx_o_streamer"
	ctxDDKey       = "ctx_o_dd"

	cookieAdminKey    = "octant_a"
	cookieStreamerKey = "octant_s"
	cookieDDKey       = "octant_d"
)

var adminCheckFail []byte = []byte(`{"code":114514,"msg":"","data":{}}`)
var streamerCheckFail []byte = []byte(`{"code":1919,"msg":"","data":{}}`)
var ddCheckFail []byte = []byte(`{"code":810,"msg":"","data":{}}`)

// -----------------------------------------------------------------

func CheckAdmin(ctx *swe.Context) {
	checkPermission[*AdminSession](ctx, cookieAdminKey, ctxAdminKey, adminCheckFail)
}

func GrantAdmin(ctx *swe.Context) {
	grantCtxPermission(ctx, cookieAdminKey, &AdminSession{})
}

func RevokeAdmin(ctx *swe.Context) {
	revokeCtxPermission(ctx, cookieAdminKey)
}

func IsAdmin(ctx *swe.Context) bool {
	_, ok := swe.CtxValue[*AdminSession](ctx, ctxAdminKey)
	return ok
}

// -----------------------------------------------------------------

func CheckStreamer(ctx *swe.Context) {
	checkPermission[*StreamerSession](ctx, cookieStreamerKey, ctxStreamerKey, streamerCheckFail)
}

func GrantStreamer(ctx *swe.Context, data *StreamerSession) {
	grantCtxPermission(ctx, cookieStreamerKey, data)
}

func RevokeStreamer(ctx *swe.Context) {
	revokeCtxPermission(ctx, cookieStreamerKey)
}

func GetStreamerSession(ctx *swe.Context) (*StreamerSession, bool) {
	return swe.CtxValue[*StreamerSession](ctx, ctxStreamerKey)
}

// -----------------------------------------------------------------

func CheckDD(ctx *swe.Context) {
	checkPermission[*DDSession](ctx, cookieDDKey, ctxDDKey, ddCheckFail)
}

func GrantDD(ctx *swe.Context, data *DDSession) {
	grantCtxPermission(ctx, cookieDDKey, data)
}

func RevokeDD(ctx *swe.Context) {
	revokeCtxPermission(ctx, cookieDDKey)
}

func GetDDSession(ctx *swe.Context) (*DDSession, bool) {
	return swe.CtxValue[*DDSession](ctx, ctxDDKey)
}

// -----------------------------------------------------------------

func checkPermission[T any](ctx *swe.Context, cookieKey, ctxKey string, fail []byte) {
	sessKey, ok := getCtxCookie(ctx, cookieKey)
	if !ok {
		ctx.Response.Write(fail)
		return
	}

	sessData, ok := getSessionData[T](sessKey)
	if !ok {
		ctx.Response.Write(fail)
		return
	}

	ctx.Put(ctxKey, sessData)
	ctx.Next()
}

func grantCtxPermission(ctx *swe.Context, cookieKey string, sessData any) {
	sessKey := GetManager().GenerateSessionKey()
	setCtxCookie(ctx, cookieKey, sessKey)
	GetManager().Set(sessKey, sessData, 3600*24*7)
}

func revokeCtxPermission(ctx *swe.Context, cookieKey string) {
	sessKey, ok := getCtxCookie(ctx, cookieKey)
	if ok {
		GetManager().Del(sessKey)
	}
	setCtxCookie(ctx, cookieKey, "")
}

func getCtxCookie(ctx *swe.Context, key string) (string, bool) {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func setCtxCookie(ctx *swe.Context, key, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   trimPort(ctx.Request.Host),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if len(value) == 0 {
		cookie.Value = "-"
		cookie.Expires = time.Now().Add(-time.Hour)
	}
	ctx.Response.Header().Add("Set-Cookie", cookie.String())
}

func trimPort(host string) string {
	if idx := strings.LastIndex(host, ":"); idx >= 0 {
		return host[:idx]
	}
	return host
}
