package handler

import (
	"fmt"

	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(GET, "/streamer/login", streamer.checkLogin, session.CheckStreamer)
	registerHandler(POST, "/streamer/login", streamer.login)
	registerHandler(GET, "/streamer/logout", streamer.logout, session.CheckStreamer)
	registerHandler(POST, "/streamer/password", streamer.changePass, session.CheckStreamer)
}

type streamerHandler struct{}

var streamer streamerHandler

func (ins streamerHandler) checkLogin(ctx *swe.Context, req *bs.Nothing) (*bs.StreamerCheckLoginRsp, swe.SweError) {
	ret, _ := session.GetStreamerSession(ctx)
	return &bs.StreamerCheckLoginRsp{
		ID:      ret.RoomID,
		Name:    ret.StreamerName,
		Account: ret.AccountName,
	}, nil
}

func (ins streamerHandler) login(ctx *swe.Context, req *bs.StreamerLoginReq) (*bs.Nothing, swe.SweError) {
	item, err := db.GetStreamerDAL().FindByAccount(req.Account)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if item == nil {
		swe.CtxLogger(ctx).Error("streamer account [%s] not exist", req.Account)
		return nil, swe.Error(EC_ST_NO_ACCOUNT, fmt.Errorf("account not exist"))
	}

	priKey, err := utils.DecryptByPass(req.Password, item.PrivateKey)
	if err != nil {
		swe.CtxLogger(ctx).Error("streamer account [%s] password incorrect", req.Account)
		return nil, swe.Error(EC_ST_PASSWORD_INCORRECT, fmt.Errorf("password incorrect"))
	}

	pubKey, err := utils.Base64Decode(item.PublicKey)
	if err != nil {
		swe.CtxLogger(ctx).Error("streamer account [%s] decode public key failed: %v", req.Account, err)
		return nil, swe.Error(EC_ST_DECODE_PUB_FAIL, err)
	}

	session.GrantStreamer(ctx, &session.StreamerSession{
		RoomID:       item.RoomID,
		StreamerName: item.StreamerName,
		AccountName:  item.AccountName,
		PrivateKey:   priKey,
		PublicKey:    pubKey,
	})
	return &bs.Nothing{}, nil
}

func (ins streamerHandler) logout(ctx *swe.Context, req *bs.Nothing) (*bs.Nothing, swe.SweError) {
	session.RevokeStreamer(ctx)
	return &bs.Nothing{}, nil
}

func (ins streamerHandler) changePass(ctx *swe.Context, req *bs.AdminChangePassReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)

	item, err := db.GetStreamerDAL().FindByAccount(st.AccountName)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if item == nil {
		swe.CtxLogger(ctx).Error("streamer account [%s] not exist", st.AccountName)
		return nil, swe.Error(EC_ST_NO_ACCOUNT, fmt.Errorf("account not exist"))
	}

	priKey, err := utils.DecryptByPass(req.Old, item.PrivateKey)
	if err != nil {
		swe.CtxLogger(ctx).Error("streamer account [%s] password incorrect", st.AccountName)
		return nil, swe.Error(EC_ST_PASSWORD_INCORRECT, fmt.Errorf("password incorrect"))
	}

	item.PrivateKey = utils.EncryptByPass(req.New, priKey)
	if err := db.GetStreamerDAL().UpdatePrivateKey(item.RoomID, item.PrivateKey); err != nil {
		swe.CtxLogger(ctx).Error("update db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	return &bs.Nothing{}, nil
}
