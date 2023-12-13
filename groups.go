package main

import (
	"fmt"
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"strings"
)

var groups = make(map[string]*P3PGROUP_GroupInfo)

type P3PGROUP_GroupInfo struct {
	gorm.Model
	GroupName    string
	GroupEmail   string
	EndpointPath string
	UniqueKey    string
	OwnerKeyID   string
	pi           *core.PrivateInfoS `gorm:"-"`
}

func createGroup(ownerKeyID string, groupName string, groupEmail string) string {
	uniqueKey := uuid.NewString()
	if groupEmail == "" {
		groupEmail = uniqueKey + "@" + botPi.Endpoint.GetUriHostDomain()
	}
	groups[uniqueKey] = &P3PGROUP_GroupInfo{
		GroupName:    groupName,
		GroupEmail:   groupEmail,
		EndpointPath: uniqueKey,
		UniqueKey:    uniqueKey,
		OwnerKeyID:   ownerKeyID,
	}

	groups[uniqueKey].pi = core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), groupName, uniqueKey, true)
	groups[uniqueKey].pi.Endpoint = core.Endpoint(string(botPi.Endpoint) + "/" + uniqueKey)
	if !groups[uniqueKey].pi.IsAccountReady() {
		groups[uniqueKey].pi.Create(groupName, groupEmail, 4096)
	}
	botPi.DB.Save(groups[uniqueKey])
	loadGroups()
	return uniqueKey
}

func loadGroups() {
	var p3pgi []P3PGROUP_GroupInfo
	botPi.DB.Find(&p3pgi)
	for i := range p3pgi {
		_, ok := groups[p3pgi[i].UniqueKey]
		if !ok {
			log.Println("Loading group:", p3pgi[i].GroupName, p3pgi[i].UniqueKey)
			groups[p3pgi[i].UniqueKey] = &p3pgi[i]
			groups[p3pgi[i].UniqueKey].pi = core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), p3pgi[i].GroupName, p3pgi[i].UniqueKey, true)
			dbAutoMigrateGroup(groups[p3pgi[i].UniqueKey].pi)
			groups[p3pgi[i].UniqueKey].pi.MessageCallback = append(groups[p3pgi[i].UniqueKey].pi.MessageCallback, groupMsgHandler)
			groups[p3pgi[i].UniqueKey].pi.IntroduceCallback = append(groups[p3pgi[i].UniqueKey].pi.IntroduceCallback, groupIntroduceHandler)
			groups[p3pgi[i].UniqueKey].pi.FileStoreElementCallback = append(groups[p3pgi[i].UniqueKey].pi.FileStoreElementCallback, groupFseCallback)
			go groups[p3pgi[i].UniqueKey].pi.EventQueueRunner()
		}
	}
}

func groupFseCallback(pi *core.PrivateInfoS, ui *core.UserInfo, fse *core.FileStoreElement, completed bool) {
	groupSendEventToAll(pi, &core.Event{
		EventType: core.EventTypeFile,
		Data: core.EventDataMixed{
			EventDataFile: core.EventDataFile{
				Uuid:       fse.Uuid,
				HttpPath:   fse.ExternalHttpPath,
				Path:       fse.Path,
				Sha512sum:  fse.Sha512sum,
				SizeBytes:  fse.SizeBytes,
				IsDeleted:  fse.IsDeleted,
				ModifyTime: fse.ModifyTime,
			},
		},
		Uuid: fse.Uuid,
	}, ui.GetKeyID())
}
func groupSendEventToAll(pi *core.PrivateInfoS, evt *core.Event, exceptUserKeyID string) {
	uis := pi.GetAllUserInfo()
	for _, targetUi := range uis {
		if targetUi.GetKeyID() == exceptUserKeyID {
			continue
		}
		mui := getMemberMetadata(pi, targetUi)
		if mui.IsUserBanned {
			continue
		}
		core.QueueEvent(pi, *evt, targetUi)
	}
}

func getGroupInfo(pi *core.PrivateInfoS) *P3PGROUP_GroupInfo {
	a := strings.Split(string(pi.Endpoint), "/")
	uid := a[len(a)-1]
	return groups[uid]
}

func getMemberMetadata(pi *core.PrivateInfoS, ui *core.UserInfo) *P3PGROUP_MemberUserInfo {
	mui := &P3PGROUP_MemberUserInfo{}
	pi.DB.First(mui, "key_id = ?", ui.KeyID)
	mui.KeyID = ui.GetKeyID()
	return mui
}

func groupSendToAllMessage(pi *core.PrivateInfoS, from string, body string, exceptUserKeyID string) {
	uis := pi.GetAllUserInfo()
	for _, ui := range uis {
		if ui.KeyID == exceptUserKeyID {
			continue
		}
		mui := getMemberMetadata(pi, ui)
		if mui.IsUserBanned {
			continue
		}
		pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("`%s`: %s", from, body))
	}
}

func groupIntroduceHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event) {
	mui := getMemberMetadata(pi, ui)
	gUi := getGroupInfo(pi)
	if mui.IsUserBanned {
		if mui.BanReason == "" {
			pi.SendMessage(ui, core.MessageTypeText, "Unable to join group chat. Reason: User is banned.")
		} else {
			pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("Unable to join group chat. Reason: %s", mui.BanReason))
		}
		// Shall we?
		// pi.PurgeUser(ui)
		return
	}
	pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("Welcome to **%s** [%s]", gUi.GroupName, gUi.GroupEmail))
	groupSendToAllMessage(pi, "service", fmt.Sprintf("User **%s** `[%s]` have joined this room", ui.Username, ui.GetKeyID()), "")
}

func groupMsgHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event, msg *core.Message) {
	uis := pi.GetAllUserInfo()
	for _, targetUi := range uis {
		mui := getMemberMetadata(pi, targetUi)
		if mui.IsUserBanned {
			continue
		}
	}
	groupSendToAllMessage(pi, ui.Username, msg.Body, ui.KeyID)
}
