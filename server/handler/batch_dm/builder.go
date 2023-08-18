package batch_dm

import (
	"fmt"
	"strings"

	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/swe"
)

var ErrInfoNotFount error = fmt.Errorf("info not found")

type BuildCtx struct {
	StreamerName string
	InviteLink   string
	Event        *db.RewardEvent
	InfoMap      map[int64]*db.DDInfo
}

type Builder interface {
	BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error)
}

type builderList []Builder

func (b builderList) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	ret := make([]string, 0, len(b))

	for _, item := range b {
		segment, err := item.BuildContent(ctx, record, buildCtx)
		if err != nil {
			return "", err
		}
		ret = append(ret, segment)
	}

	return strings.Join(ret, ""), nil
}

type plainBuilder string

func (b plainBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	return string(b), nil
}

type streamerNameBuilder struct{}

func (b streamerNameBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	return buildCtx.StreamerName, nil
}

type roomIDBuilder struct{}

func (b roomIDBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	return fmt.Sprint(buildCtx.Event.RoomID), nil
}

type eventNameBuilder struct{}

func (b eventNameBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	return buildCtx.Event.EventName, nil
}

type uidBuilder struct{}

func (b uidBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	info, ok := buildCtx.InfoMap[record.UID]
	if !ok {
		return "", ErrInfoNotFount
	}
	return fmt.Sprint(info.UID), nil
}

type nameBuilder struct{}

func (b nameBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	info, ok := buildCtx.InfoMap[record.UID]
	if !ok {
		return "", ErrInfoNotFount
	}
	return info.UserName, nil
}

type inviteBuilder struct{}

func (b inviteBuilder) BuildContent(ctx *swe.Context, record *db.RewardUser, buildCtx *BuildCtx) (string, error) {
	info, ok := buildCtx.InfoMap[record.UID]
	if !ok {
		return "", ErrInfoNotFount
	}
	return buildCtx.InviteLink + info.AccessCode, nil
}

func MakeBuilder(template string) Builder {
	state := 0
	b := strings.Builder{}

	ret := builderList{}

	for _, ch := range template {
		switch state {
		case 0:
			if ch == '{' {
				state = 1
				ret = append(ret, plainBuilder(b.String()))
				b = strings.Builder{}
			} else {
				b.WriteRune(ch)
			}
		case 1:
			if ch == '}' {
				state = 0
				ret = append(ret, templateBuilder(b.String()))
				b = strings.Builder{}
			} else {
				b.WriteRune(ch)
			}
		}
	}

	if b.Len() > 0 {
		if state == 0 {
			ret = append(ret, plainBuilder(b.String()))
		} else {
			ret = append(ret, templateBuilder(b.String()))
		}
	}
	return ret
}

func templateBuilder(value string) Builder {
	tmp := strings.ToLower(value)
	switch tmp {
	case "streamer":
		return streamerNameBuilder{}
	case "room":
		return roomIDBuilder{}
	case "event":
		return eventNameBuilder{}
	case "uid":
		return uidBuilder{}
	case "name":
		return nameBuilder{}
	case "invite_link":
		return inviteBuilder{}
	}
	return plainBuilder("{" + value + "}")
}
