package bs

import (
	"fmt"

	"github.com/zerozwt/swe"
)

type Nothing struct{}

type PageReq struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (req PageReq) Validate(*swe.Context) error {
	if req.Page < 1 {
		return fmt.Errorf("invalid page %d", req.Page)
	}
	if req.Size < 1 || req.Size > 100 {
		return fmt.Errorf("invalid size %d", req.Size)
	}
	return nil
}

type PageRsp struct {
	Count int   `json:"count"`
	List  []any `json:"list"`
}

type IDReq struct {
	ID int64 `json:"id"`
}
