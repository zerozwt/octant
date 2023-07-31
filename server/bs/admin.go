package bs

import (
	"fmt"

	"github.com/zerozwt/swe"
)

type AdminLoginReq struct {
	Password string `json:"password"`
}

type AdminChangePassReq struct {
	Old string `json:"old_password"`
	New string `json:"new_password"`
}

func (req AdminChangePassReq) Validate(ctx *swe.Context) error {
	if len(req.Old) == 0 {
		return fmt.Errorf("empty old password")
	}
	if len(req.New) == 0 {
		return fmt.Errorf("empty new password")
	}
	return nil
}

type AdminLiveRoomListItem struct {
	ID           int64  `json:"room_id"`
	StreamerName string `json:"name"`
	Account      string `json:"account_name"`
}

type AdminCreateStreamerReq struct {
	ID       int64  `json:"room_id"`
	Account  string `json:"name"`
	Password string `json:"password"`
}

func (req AdminCreateStreamerReq) Validate(ctx *swe.Context) error {
	if len(req.Password) == 0 {
		return fmt.Errorf("empty password")
	}
	if req.ID < 1 {
		return fmt.Errorf("invalid room id %d", req.ID)
	}
	if len(req.Account) == 0 {
		return fmt.Errorf("empty account")
	}
	return nil
}

type AdminResetStreamerReq struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

func (req AdminResetStreamerReq) Validate(ctx *swe.Context) error {
	if len(req.Password) == 0 {
		return fmt.Errorf("empty password")
	}
	if req.ID < 1 {
		return fmt.Errorf("invalid room id %d", req.ID)
	}
	return nil
}
