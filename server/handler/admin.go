package handler

import (
	"fmt"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/octant/server/bridge"
	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(GET, "/admin/login", admin.checkLogin, session.CheckAdmin)
	registerHandler(POST, "/admin/login", admin.login)
	registerHandler(GET, "/admin/logout", admin.logout, session.CheckAdmin)
	registerHandler(POST, "/admin/password", admin.changePass, session.CheckAdmin)

	registerHandler(GET, "/admin/streamer/list", admin.streamerList, session.CheckAdmin)
	registerHandler(POST, "/admin/streamer/add", admin.createStreamer, session.CheckAdmin)
	registerHandler(POST, "/admin/streamer/delete", admin.deleteStreamer, session.CheckAdmin)
	registerHandler(POST, "/admin/streamer/reset", admin.resetStreamerPassword, session.CheckAdmin)
}

type adminHandler struct{}

var admin adminHandler

func (ins adminHandler) checkLogin(ctx *swe.Context, unused *bs.Nothing) (*bs.Nothing, swe.SweError) {
	return &bs.Nothing{}, nil
}

func (ins adminHandler) login(ctx *swe.Context, req *bs.AdminLoginReq) (*bs.Nothing, swe.SweError) {
	passInDB, err := db.GetSysConfigDAL().GetConfig(ctx, db.DB_SYSCONF_ADMIN_PASS)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	passReq := db.GetSysConfigDAL().EncodeAdminPassword(req.Password)

	if passInDB != passReq {
		return nil, swe.Error(EC_ADMIN_PASSWORD_INCORRECT, fmt.Errorf("admin password incorrect"))
	}

	session.GrantAdmin(ctx)

	return &bs.Nothing{}, nil
}

func (ins adminHandler) logout(ctx *swe.Context, req *bs.Nothing) (*bs.Nothing, swe.SweError) {
	session.RevokeAdmin(ctx)
	return &bs.Nothing{}, nil
}

func (ins adminHandler) changePass(ctx *swe.Context, req *bs.AdminChangePassReq) (*bs.Nothing, swe.SweError) {
	passInDB, err := db.GetSysConfigDAL().GetConfig(ctx, db.DB_SYSCONF_ADMIN_PASS)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	passOld := db.GetSysConfigDAL().EncodeAdminPassword(req.Old)

	if passInDB != passOld {
		return nil, swe.Error(EC_ADMIN_PASSWORD_INCORRECT, fmt.Errorf("admin password incorrect"))
	}

	passNew := db.GetSysConfigDAL().EncodeAdminPassword(req.New)
	err = db.GetSysConfigDAL().SetConfig(ctx, db.DB_SYSCONF_ADMIN_PASS, passNew)
	if err != nil {
		swe.CtxLogger(ctx).Error("write to db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	return &bs.Nothing{}, nil
}

func (ins adminHandler) streamerList(ctx *swe.Context, req *bs.PageReq) (*bs.PageRsp, swe.SweError) {
	count, streamers, err := db.GetStreamerDAL().Page(ctx, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	ret := &bs.PageRsp{
		Count: count,
		List:  []any{},
	}

	for _, item := range streamers {
		ret.List = append(ret.List, bs.AdminLiveRoomListItem{
			ID:           item.RoomID,
			StreamerName: item.StreamerName,
			Account:      item.AccountName,
		})
	}

	return ret, nil
}

func (ins adminHandler) createStreamer(ctx *swe.Context, req *bs.AdminCreateStreamerReq) (*bs.Nothing, swe.SweError) {
	// get live room info
	info, err := dm.GetRoomInfo(req.ID)
	if err != nil {
		swe.CtxLogger(ctx).Error("get room info for live room %d failed: %v", req.ID, err)
		return nil, swe.Error(EC_ADMIN_ROOM_INFO_FAIL, err)
	}

	item := db.Streamer{
		RoomID:       req.ID,
		StreamerName: info.Liver.Base.Name,
		AccountName:  req.Account,
	}

	priKey, pubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		swe.CtxLogger(ctx).Error("generate key pair for live room %d failed: %v", req.ID, err)
		return nil, swe.Error(EC_ADMIN_KEYGEN_FAIL, err)
	}

	item.PrivateKey = utils.EncryptByPass(req.Password, priKey)
	item.PublicKey = utils.Base64Encode(pubKey)

	// try insert into table
	rows, err := db.GetStreamerDAL().Insert(ctx, &item, false)
	if err != nil {
		swe.CtxLogger(ctx).Error("insert into db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if rows <= 0 {
		swe.CtxLogger(ctx).Error("add live room %d failed: room id or account duplicated", req.ID)
		return nil, swe.Error(EC_ADMIN_DUPLICATED_STREAMER, fmt.Errorf("duplicated streamer"))
	}

	// start tracking live room
	err = bridge.GetBridge().AddRoom(req.ID)
	if err != nil {
		swe.CtxLogger(ctx).Error("add room %d to bridge failed: %v", req.ID, err)
	}

	return &bs.Nothing{}, nil
}

func (ins adminHandler) deleteStreamer(ctx *swe.Context, req *bs.IDReq) (*bs.Nothing, swe.SweError) {
	// delete from db
	rows, err := db.GetStreamerDAL().Delete(ctx, req.ID)
	if err != nil {
		swe.CtxLogger(ctx).Error("delete from db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if rows > 0 {
		// stop tracking live room
		err = bridge.GetBridge().DelRoom(req.ID)
		if err != nil {
			swe.CtxLogger(ctx).Error("delete room %d from bridge failed: %v", req.ID, err)
		}
	}

	return &bs.Nothing{}, nil
}

func (ins adminHandler) resetStreamerPassword(ctx *swe.Context, req *bs.AdminResetStreamerReq) (*bs.Nothing, swe.SweError) {
	// find streamer from db
	streamer, err := db.GetStreamerDAL().Find(ctx, req.ID)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if streamer == nil {
		swe.CtxLogger(ctx).Info("live room %d not exist", req.ID)
		return &bs.Nothing{}, nil
	}

	// generate new keypair
	priKey, pubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		swe.CtxLogger(ctx).Error("generate key pair for live room %d failed: %v", req.ID, err)
		return nil, swe.Error(EC_ADMIN_KEYGEN_FAIL, err)
	}

	streamer.PrivateKey = utils.EncryptByPass(req.Password, priKey)
	streamer.PublicKey = utils.Base64Encode(pubKey)

	// update keypair
	row, err := db.GetStreamerDAL().Insert(ctx, streamer, true)
	if err != nil {
		swe.CtxLogger(ctx).Error("insert into db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if row <= 0 {
		swe.CtxLogger(ctx).Error("update keypair for live room %d affected 0 row", req.ID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("record not affected"))
	}

	return &bs.Nothing{}, nil
}
