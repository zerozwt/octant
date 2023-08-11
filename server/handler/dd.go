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
	registerHandler(GET, "/dd/access", dd.access, session.CheckDDInAccess)
	registerHandler(GET, "/dd/login", dd.checkLogin, session.CheckDD)
	registerHandler(POST, "/dd/login", dd.login)
	registerHandler(GET, "/dd/logout", dd.logout, session.CheckDD)
	registerHandler(POST, "/dd/password", dd.setPassword, session.CheckDDInAccess)
}

type ddHandler struct{}

var dd ddHandler

func (ins ddHandler) access(ctx *swe.Context, req *bs.DDACLoginReq) (*bs.DDACLoginRsp, swe.SweError) {
	if user, ok := session.GetDDSession(ctx); ok {
		return &bs.DDACLoginRsp{UID: user.UID, Name: user.UserName}, nil
	}

	info, err := db.GetDDInfoDAL().GetByAccessCode(ctx, req.AccessCode)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if info == nil {
		swe.CtxLogger(ctx).Error("query user by access code %s failed: not found", req.AccessCode)
		return nil, swe.Error(EC_DD_CODE_NOT_FOUND, fmt.Errorf("access code not found"))
	}

	ret := bs.DDACLoginRsp{
		UID:  info.UID,
		Name: info.UserName,
	}
	if len(info.PublicKey) > 0 {
		ret.NeedPassword = true
	} else {
		ret.PasswordNotSet = true
	}

	return &ret, nil
}

func (ins ddHandler) checkLogin(ctx *swe.Context, req *bs.Nothing) (*bs.DDACLoginRsp, swe.SweError) {
	user, _ := session.GetDDSession(ctx)
	return &bs.DDACLoginRsp{UID: user.UID, Name: user.UserName}, nil
}

func (ins ddHandler) login(ctx *swe.Context, req *bs.DDLoginReq) (*bs.Nothing, swe.SweError) {
	logger := swe.CtxLogger(ctx)

	info, err := db.GetDDInfoDAL().Get(ctx, req.UID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if info == nil {
		logger.Error("query dd by uid %d failed: not found", req.UID)
		return nil, swe.Error(EC_DD_CODE_NOT_FOUND, fmt.Errorf("not found"))
	}

	if len(info.PrivateKey) == 0 {
		logger.Error("decode private key for uid %d failed: password not set", req.UID)
		return nil, swe.Error(EC_DD_PASSWORD_INCORRECT, fmt.Errorf("not set"))
	}

	sess := session.DDSession{
		UID:      info.UID,
		UserName: info.UserName,
	}

	sess.PrivateKey, err = utils.DecryptByPass(req.Password, info.PrivateKey)
	if err != nil {
		logger.Error("decode private key for uid %d failed: %v", req.UID, err)
		return nil, swe.Error(EC_DD_PASSWORD_INCORRECT, err)
	}

	sess.PublicKey, err = utils.Base64Decode(info.PublicKey)
	if err != nil {
		logger.Error("decode public key for uid %d failed: %v", req.UID, err)
		return nil, swe.Error(EC_DD_PASSWORD_INCORRECT, err)
	}

	session.GrantDD(ctx, &sess)
	return &bs.Nothing{}, nil
}

func (ins ddHandler) logout(ctx *swe.Context, req *bs.Nothing) (*bs.Nothing, swe.SweError) {
	session.RevokeDD(ctx)
	return &bs.Nothing{}, nil
}

func (ins ddHandler) setPassword(ctx *swe.Context, req *bs.DDSetPasswordReq) (*bs.Nothing, swe.SweError) {
	if _, ok := session.GetDDSession(ctx); ok {
		return &bs.Nothing{}, ins.setPasswordByPassword(ctx, req)
	}
	return &bs.Nothing{}, ins.setPasswordByAccessCode(ctx, req)
}

func (ins ddHandler) setPasswordByPassword(ctx *swe.Context, req *bs.DDSetPasswordReq) swe.SweError {
	user, _ := session.GetDDSession(ctx)
	logger := swe.CtxLogger(ctx)

	info, err := db.GetDDInfoDAL().Get(ctx, user.UID)
	if err != nil {
		logger.Error("query db error %v", err)
		return swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if info == nil {
		swe.CtxLogger(ctx).Error("query user by uid %d failed: not found", info.UID)
		return swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("user not found"))
	}

	if len(info.PrivateKey) == 0 {
		// should never be here
		logger.Error("update password for uid %d failed: password not set", info.UID)
		return swe.Error(EC_DD_PASSWORD_INCORRECT, fmt.Errorf("password not set"))
	}

	priKey, err := utils.DecryptByPass(req.Old, info.PrivateKey)
	if err != nil {
		logger.Error("decode private key for uid %d failed: %v", info.UID, err)
		return swe.Error(EC_DD_PASSWORD_INCORRECT, fmt.Errorf("password incorrect"))
	}

	info.PrivateKey = utils.EncryptByPass(req.New, priKey)

	if err = db.GetDDInfoDAL().SetKeyPair(ctx, info); err != nil {
		logger.Error("update db error %v", err)
		return swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	return nil
}

func (ins ddHandler) setPasswordByAccessCode(ctx *swe.Context, req *bs.DDSetPasswordReq) swe.SweError {
	logger := swe.CtxLogger(ctx)
	info, err := db.GetDDInfoDAL().GetByAccessCode(ctx, req.AccessCode)
	if err != nil {
		logger.Error("query db error %v", err)
		return swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if info == nil {
		logger.Error("query user by access code %s failed: not found", req.AccessCode)
		return swe.Error(EC_DD_CODE_NOT_FOUND, fmt.Errorf("access code not found"))
	}

	if len(info.PrivateKey) > 0 {
		logger.Error("query user by access code %s failed: password already set", req.AccessCode)
		return swe.Error(EC_DD_PASSWORD_INCORRECT, fmt.Errorf("access code cannot be used in changing password"))
	}

	priKey, pubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		logger.Error("generate key pair for user %d failed: %v", info.UID, err)
		return swe.Error(EC_DD_KEYGEN_FAIL, err)
	}

	info.PrivateKey = utils.EncryptByPass(req.New, priKey)
	info.PublicKey = utils.Base64Encode(pubKey)

	if err = db.GetDDInfoDAL().SetKeyPair(ctx, info); err != nil {
		logger.Error("update db error %v", err)
		return swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	return nil
}
