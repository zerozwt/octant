package event_calc

import (
	"fmt"
	"strings"

	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

const (
	CTX_KEY_ADDR = "oct_user_addr"
)

type AddrMap map[int64]*db.RewardUserAddress

type Picker interface {
	Pick(ctx *swe.Context, user *UserData) string
	Header(ctx *swe.Context) string
}

// -----------------------------------------------------------------

type uidPicker struct{}

func (p uidPicker) Pick(ctx *swe.Context, user *UserData) string { return "'" + fmt.Sprint(user.UID) }
func (p uidPicker) Header(ctx *swe.Context) string               { return "B站UID" }

type namePicker struct{}

func (p namePicker) Pick(ctx *swe.Context, user *UserData) string { return user.Name }
func (p namePicker) Header(ctx *swe.Context) string               { return "用户昵称" }

// -----------------------------------------------------------------

type giftNamePicker struct{}

func (p giftNamePicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Gift {
		tmp = append(tmp, item.GiftName)
	}

	return strings.Join(tmp, "\n")
}

func (p giftNamePicker) Header(ctx *swe.Context) string { return "礼物名称" }

type giftNumPicker struct{}

func (p giftNumPicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Gift {
		tmp = append(tmp, fmt.Sprint(item.GiftCount))
	}

	return strings.Join(tmp, "\n")
}

func (p giftNumPicker) Header(ctx *swe.Context) string { return "礼物数量" }

type giftTimePicker struct{}

func (p giftTimePicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Gift {
		tmp = append(tmp, utils.TimeToCSTString(item.SendTime))
	}

	return strings.Join(tmp, "\n")
}

func (p giftTimePicker) Header(ctx *swe.Context) string { return "送礼时间" }

// -----------------------------------------------------------------

type scPicker struct{}

func (p scPicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.SC {
		tmp = append(tmp, item.Content)
	}

	return strings.Join(tmp, "\n")
}

func (p scPicker) Header(ctx *swe.Context) string { return "SC内容" }

type scPricePicker struct{}

func (p scPricePicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.SC {
		tmp = append(tmp, fmt.Sprint(item.Price))
	}

	return strings.Join(tmp, "\n")
}

func (p scPricePicker) Header(ctx *swe.Context) string { return "SC金额" }

type scTimePicker struct{}

func (p scTimePicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.SC {
		tmp = append(tmp, utils.TimeToCSTString(item.SendTime))
	}

	return strings.Join(tmp, "\n")
}

func (p scTimePicker) Header(ctx *swe.Context) string { return "发言时间" }

// -----------------------------------------------------------------

type memberPicker struct{}

func (p memberPicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Member {
		switch item.GuardLevel {
		case 1:
			tmp = append(tmp, "总督")
		case 2:
			tmp = append(tmp, "提督")
		case 3:
			tmp = append(tmp, "舰长")
		}
	}

	return strings.Join(tmp, "\n")
}

func (p memberPicker) Header(ctx *swe.Context) string { return "大航海类型" }

type memberCountPicker struct{}

func (p memberCountPicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Member {
		tmp = append(tmp, fmt.Sprint(item.Count))
	}

	return strings.Join(tmp, "\n")
}

func (p memberCountPicker) Header(ctx *swe.Context) string { return "大航海月数" }

type memberTimePicker struct{}

func (p memberTimePicker) Pick(ctx *swe.Context, user *UserData) string {
	tmp := []string{}

	for _, item := range user.Member {
		tmp = append(tmp, utils.TimeToCSTString(item.SendTime))
	}

	return strings.Join(tmp, "\n")
}

func (p memberTimePicker) Header(ctx *swe.Context) string { return "上舰时间" }

// -----------------------------------------------------------------

type recvNamePicker struct{}

func (p recvNamePicker) Pick(ctx *swe.Context, user *UserData) string {
	addr, ok := swe.CtxValue[AddrMap](ctx, CTX_KEY_ADDR)
	if !ok {
		return ""
	}
	if item, ok := addr[user.UID]; ok {
		return item.Name
	}
	return ""
}

func (p recvNamePicker) Header(ctx *swe.Context) string { return "收件人姓名" }

type recvPhonePicker struct{}

func (p recvPhonePicker) Pick(ctx *swe.Context, user *UserData) string {
	addr, ok := swe.CtxValue[AddrMap](ctx, CTX_KEY_ADDR)
	if !ok {
		return ""
	}
	if item, ok := addr[user.UID]; ok {
		return item.Phone
	}
	return ""
}

func (p recvPhonePicker) Header(ctx *swe.Context) string { return "收件人电话" }

type recvAddrPicker struct{}

func (p recvAddrPicker) Pick(ctx *swe.Context, user *UserData) string {
	addr, ok := swe.CtxValue[AddrMap](ctx, CTX_KEY_ADDR)
	if !ok {
		return ""
	}
	if item, ok := addr[user.UID]; ok {
		return item.Addr
	}
	return ""
}

func (p recvAddrPicker) Header(ctx *swe.Context) string { return "收件地址" }

// -----------------------------------------------------------------

func BuildPickers(ctx *swe.Context, validatedCondition *bs.EventCondition) []Picker {
	ret := []Picker{uidPicker{}, namePicker{}}

	tr := bs.ConditionTimeRange{}
	validatedCondition.CalculateRange(&tr)

	if _, ok := tr.Range["gift"]; ok {
		ret = append(ret, giftNamePicker{}, giftNumPicker{}, giftTimePicker{})
	}
	if _, ok := tr.Range["sc"]; ok {
		ret = append(ret, scPicker{}, scPricePicker{}, scTimePicker{})
	}
	if _, ok := tr.Range["member"]; ok {
		ret = append(ret, memberPicker{}, memberCountPicker{}, memberTimePicker{})
	}

	ret = append(ret, recvNamePicker{}, recvPhonePicker{}, recvAddrPicker{})

	return ret
}

func Table(ctx *swe.Context, users []*UserData, pickers []Picker) [][]string {
	ret := [][]string{}

	// header
	header := make([]string, 0, len(pickers))
	for _, picker := range pickers {
		header = append(header, picker.Header(ctx))
	}
	ret = append(ret, header)

	// data
	for _, user := range users {
		line := make([]string, 0, len(pickers))
		for _, picker := range pickers {
			line = append(line, picker.Pick(ctx, user))
		}
		ret = append(ret, line)
	}

	return ret
}
