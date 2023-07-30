package handler

import (
	"fmt"

	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(GET, "/admin/login", admin.checkLogin, session.CheckAdmin)
	registerHandler(POST, "/admin/login", admin.login)
	registerHandler(POST, "/admin/logout", admin.logout, session.CheckAdmin)
}

type adminHandler struct{}

var admin adminHandler

func (ins adminHandler) checkLogin(ctx *swe.Context, unused *bs.Nothing) (*bs.Nothing, swe.SweError) {
	return &bs.Nothing{}, nil
}

func (ins adminHandler) login(ctx *swe.Context, req *bs.AdminLoginReq) (*bs.Nothing, swe.SweError) {
	passInDB, err := db.GetSysConfigDAL().GetConfig(db.DB_SYSCONF_ADMIN_PASS)
	if err != nil {
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	passReq := db.GetSysConfigDAL().EncodeAdminPassword(req.Password)

	if passInDB != passReq {
		return nil, swe.Error(1002, fmt.Errorf("admin password incorrect"))
	}

	session.GrantAdmin(ctx)

	return &bs.Nothing{}, nil
}

func (ins adminHandler) logout(ctx *swe.Context, req *bs.Nothing) (*bs.Nothing, swe.SweError) {
	session.RevokeAdmin(ctx)
	return &bs.Nothing{}, nil
}
