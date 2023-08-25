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

	registerHandler(GET, "/dd/events", dd.events, session.CheckDD)
	registerHandler(GET, "/dd/address", dd.getAddress, session.CheckDD)
	registerHandler(POST, "/dd/address", dd.setAddress, session.CheckDD)
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

func (ins ddHandler) events(ctx *swe.Context, req *bs.PageReq) (*bs.PageRsp, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	user, _ := session.GetDDSession(ctx)

	// get user records
	count, list, err := db.GetRewardEventDAL().UserRecords(ctx, user.UID, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	ret := bs.PageRsp{
		Count: count,
		List:  []any{},
	}

	if count == 0 {
		return &ret, nil
	}

	// get events
	eids := []int64{}
	for _, item := range list {
		eids = append(eids, item.EventID)
	}
	events, err := db.GetRewardEventDAL().GetByIDs(ctx, eids)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	rids := []int64{}
	evtMap := map[int64]*db.RewardEvent{}
	for idx := range events {
		if events[idx].Hidden != 0 {
			continue
		}
		rids = append(rids, events[idx].RoomID)
		evtMap[events[idx].ID] = &events[idx]
	}

	// get streamers
	streamers, err := db.GetStreamerDAL().FindByRoomIDs(ctx, rids)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	stMap := map[int64]*db.Streamer{}
	for idx := range streamers {
		stMap[streamers[idx].RoomID] = &streamers[idx]
	}

	for _, item := range list {
		data := bs.DDEventItem{
			ID:   item.EventID,
			Addr: len(item.AddressInfo) > 0,
		}

		event, ok := evtMap[data.ID]
		if !ok {
			continue
		}
		data.Name = event.EventName
		data.Reward = event.RewardContent

		st, ok := stMap[event.RoomID]
		if !ok {
			continue
		}
		data.Streamer.Name = st.StreamerName
		data.Streamer.RoomID = st.RoomID

		ret.List = append(ret.List, data)
	}

	return &ret, nil
}

func (ddHandler) getAddress(ctx *swe.Context, req *bs.IDReq) (*bs.DDAddrInfo, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	user, _ := session.GetDDSession(ctx)

	// query records
	records, err := db.GetRewardEventDAL().UserInfoByUIDAndEventID(ctx, user.UID, req.ID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if len(records) == 0 {
		logger.Error("uid %d event id %d no records", user.UID, req.ID)
		return &bs.DDAddrInfo{}, nil
	}
	if len(records[0].AddressInfo) == 0 {
		return &bs.DDAddrInfo{EventID: req.ID}, nil
	}

	// query event
	event, err := db.GetRewardEventDAL().Get(ctx, req.ID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if event == nil {
		logger.Error("event %d not found", req.ID)
		return &bs.DDAddrInfo{EventID: req.ID}, nil
	}

	// query streamer
	st, err := db.GetStreamerDAL().Find(ctx, event.RoomID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if st == nil {
		logger.Error("streamer room %d not found", event.RoomID)
		return &bs.DDAddrInfo{EventID: req.ID}, nil
	}

	// decrypt
	pubKey, err := utils.Base64Decode(st.PublicKey)
	if err != nil {
		logger.Error("decode streamer %d public key failed: %v", st.RoomID, err)
		return &bs.DDAddrInfo{EventID: req.ID}, nil
	}

	info, err := utils.DecryptUserAddress(ctx, pubKey, records[0].AddressInfo)
	if err != nil {
		logger.Error("decrypt address info failed: %v", err)
		return &bs.DDAddrInfo{EventID: req.ID}, nil
	}

	ret := &bs.DDAddrInfo{
		EventID: req.ID,
		Name:    info.Name,
		Phone:   info.Phone,
		Addr:    info.Addr,
	}

	return ret, nil
}

func (ins ddHandler) setAddress(ctx *swe.Context, req *bs.DDAddrInfo) (*bs.Nothing, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	user, _ := session.GetDDSession(ctx)

	// query event
	event, err := db.GetRewardEventDAL().Get(ctx, req.EventID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if event == nil {
		logger.Error("event %d not found", req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}

	// query streamer
	st, err := db.GetStreamerDAL().Find(ctx, event.RoomID)
	if err != nil {
		logger.Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if st == nil {
		logger.Error("streamer room %d not found", event.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("streamer not found"))
	}

	pubKey, err := utils.Base64Decode(st.PublicKey)
	if err != nil {
		logger.Error("decode streamer %d public key failed: %v", st.RoomID, err)
		return nil, swe.Error(EC_ST_DECODE_PUB_FAIL, fmt.Errorf("streamer pubkey decode failed"))
	}

	addrData, err := utils.EncryptUserAddress(ctx, pubKey, &utils.RewardUserAddress{
		Name:  req.Name,
		Phone: req.Phone,
		Addr:  req.Addr,
	})

	if err != nil {
		logger.Error("encrypt addr data failed: %v", err)
		return nil, swe.Error(EC_DD_ADDR_ENC_FAIL, err)
	}

	rows, err := db.GetRewardEventDAL().UpdateAddrInfo(ctx, user.UID, req.EventID, addrData)
	if err != nil {
		logger.Error("update address info failed: uid %d event %d : %v", user.UID, req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if rows < 1 {
		logger.Error("update address info failed: uid %d event %d : no rows affected", user.UID, req.EventID)
		return nil, swe.Error(EC_DD_SET_ADDR_FAIL, fmt.Errorf("no records updated"))
	}
	return &bs.Nothing{}, nil
}
