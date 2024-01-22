package structs

import (
	"fmt"
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"strings"
)

func GroupSendEventToAll(pi *core.PrivateInfoS, evt *core.Event, exceptUserKeyID string) {
	uis := pi.GetAllUserInfo()
	for _, targetUi := range uis {
		if targetUi.GetKeyID() == exceptUserKeyID {
			continue
		}
		mui := GetMemberMetadata(pi, targetUi)
		if mui.IsUserBanned {
			continue
		}
		core.QueueEvent(pi, *evt, targetUi)
	}
}

func GetGroupInfo(pi *core.PrivateInfoS) *P3PGROUP_GroupInfo {
	a := strings.Split(string(pi.Endpoint), "/")
	uid := a[len(a)-1]
	return Groups[uid]
}

func GetMemberMetadata(pi *core.PrivateInfoS, ui *core.UserInfo) *P3PGROUP_MemberUserInfo {
	mui := &P3PGROUP_MemberUserInfo{}
	pi.DB.First(mui, "key_id = ?", ui.KeyID)
	mui.KeyID = ui.GetKeyID()
	return mui
}

func GroupSendToAllMessage(pi *core.PrivateInfoS, from string, body string, exceptUserKeyID string) {
	uis := pi.GetAllUserInfo()
	for _, ui := range uis {
		if ui.KeyID == exceptUserKeyID {
			continue
		}
		mui := GetMemberMetadata(pi, ui)
		if mui.IsUserBanned {
			continue
		}
		pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("`%s`: %s", from, body))
	}
}
