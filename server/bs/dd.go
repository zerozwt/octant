package bs

import (
	"fmt"

	"github.com/zerozwt/swe"
)

type DDACLoginReq struct {
	AccessCode string `form:"access_code"`
}

func (req DDACLoginReq) Validate(ctx *swe.Context) error {
	if len(req.AccessCode) == 0 {
		return fmt.Errorf("empty access code")
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
	return nil
}
