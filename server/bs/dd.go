package bs

import (
	"fmt"

	"github.com/zerozwt/swe"
)

type DDACLoginReq struct {
	AccessCode string `form:"access_code"`
}

func (req DDACLoginReq) Validate(ctx *swe.Context) error {
	if len(req.AccessCode) == 0 || len(req.AccessCode) > 256 {
		return fmt.Errorf("invalid access code size")
	}
	return nil
}

type DDACLoginRsp struct {
	UID            int64  `json:"uid"`
	Name           string `json:"name"`
	NeedPassword   bool   `json:"need_password"`
	PasswordNotSet bool   `json:"password_not_set"`
}

type DDLoginReq struct {
	UID      int64  `json:"uid"`
	Password string `json:"password"`
}

func (req DDLoginReq) Validate(ctx *swe.Context) error {
	if len(req.Password) > 256 {
		return fmt.Errorf("password too long")
	}
	return nil
}

type DDSetPasswordReq struct {
	AccessCode string `json:"access_code"`
	Old        string `json:"old_password"`
	New        string `json:"new_password"`
}

func (req DDSetPasswordReq) Validate(ctx *swe.Context) error {
	if len(req.AccessCode) == 0 && len(req.Old) == 0 {
		return fmt.Errorf("empty old pass and access code")
	}
	if len(req.New) == 0 {
		return fmt.Errorf("empty new password")
	}
	if len(req.AccessCode) > 256 {
		return fmt.Errorf("access code too long")
	}
	if len(req.Old) > 256 {
		return fmt.Errorf("old password too long")
	}
	if len(req.New) > 256 {
		return fmt.Errorf("new password too long")
	}
	return nil
}

type DDEventItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Reward   string `json:"reward"`
	Addr     bool   `json:"addr"`
	Streamer struct {
		RoomID int64  `json:"room_id"`
		Name   string `json:"name"`
	} `json:"streamer"`
}

type DDAddrInfo struct {
	EventID int64  `json:"event_id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Addr    string `json:"address"`
}

func (req DDAddrInfo) Validate(ctx *swe.Context) error {
	if len(req.Name) > 256 {
		return fmt.Errorf("name too long")
	}
	if len(req.Phone) > 256 {
		return fmt.Errorf("phonw too long")
	}
	if len(req.Addr) > 1024 {
		return fmt.Errorf("addr too long")
	}
	return nil
}
