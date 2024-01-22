package main

import (
	"fmt"
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"git.mrcyjanek.net/p3pch4t/p3pgroup/extra"
	"git.mrcyjanek.net/p3pch4t/p3pgroup/structs"
	"github.com/google/uuid"
	"log"
	"os"
	"path"
)

func createGroup(ownerKeyID string, groupName string, groupEmail string) string {
	uniqueKey := uuid.NewString()
	if groupEmail == "" {
		groupEmail = uniqueKey + "@" + botPi.Endpoint.GetUriHostDomain()
	}
	structs.Groups[uniqueKey] = &structs.P3PGROUP_GroupInfo{
		GroupName:    groupName,
		GroupEmail:   groupEmail,
		EndpointPath: uniqueKey,
		UniqueKey:    uniqueKey,
		OwnerKeyID:   ownerKeyID,
	}

	structs.Groups[uniqueKey].PI = core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), groupName, uniqueKey, true)
	structs.Groups[uniqueKey].PI.Endpoint = core.Endpoint(string(botPi.Endpoint) + "/" + uniqueKey)
	if !structs.Groups[uniqueKey].PI.IsAccountReady() {
		structs.Groups[uniqueKey].PI.Create(groupName, groupEmail, 4096)
	}
	botPi.DB.Save(structs.Groups[uniqueKey])
	loadGroups()
	return uniqueKey
}

func loadGroups() {
	var p3pgi []structs.P3PGROUP_GroupInfo
	botPi.DB.Find(&p3pgi)
	for i := range p3pgi {
		_, ok := structs.Groups[p3pgi[i].UniqueKey]
		if !ok {
			log.Println("Loading group:", p3pgi[i].GroupName, p3pgi[i].UniqueKey)
			structs.Groups[p3pgi[i].UniqueKey] = &p3pgi[i]
			structs.Groups[p3pgi[i].UniqueKey].PI = core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), p3pgi[i].GroupName, p3pgi[i].UniqueKey, true)
			structs.DbAutoMigrateGroup(structs.Groups[p3pgi[i].UniqueKey].PI)
			structs.Groups[p3pgi[i].UniqueKey].PI.MessageCallback = append(structs.Groups[p3pgi[i].UniqueKey].PI.MessageCallback, groupMsgHandler)
			structs.Groups[p3pgi[i].UniqueKey].PI.IntroduceCallback = append(structs.Groups[p3pgi[i].UniqueKey].PI.IntroduceCallback, groupIntroduceHandler)
			structs.Groups[p3pgi[i].UniqueKey].PI.FileStoreElementCallback = append(structs.Groups[p3pgi[i].UniqueKey].PI.FileStoreElementCallback, groupFseCallback)
			go structs.Groups[p3pgi[i].UniqueKey].PI.EventQueueRunner()
		}
	}
}

func groupFseCallback(pi *core.PrivateInfoS, ui *core.UserInfo, fse *core.FileStoreElement, completed bool) {
	structs.GroupSendEventToAll(pi, &core.Event{
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

func groupIntroduceHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event) {
	mui := structs.GetMemberMetadata(pi, ui)
	gUi := structs.GetGroupInfo(pi)
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
	if !mui.IsUserWelcomed {
		structs.GroupSendToAllMessage(pi, "service", fmt.Sprintf("User **%s** `[%s]` have joined this room", ui.Username, ui.GetKeyID()), "")
		mui.IsUserWelcomed = true
		pi.DB.Save(mui)
	}
	if gUi.OwnerKeyID == mui.KeyID {
		if !mui.IsUserAdmin {
			pi.SendMessage(ui, core.MessageTypeText, "**Admin permissions granted. Welcome aboard.**")
		}
		mui.IsUserAdmin = true
		pi.DB.Save(mui)
	}
}

func groupMsgHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event, msg *core.Message) {
	uis := pi.GetAllUserInfo()
	for _, targetUi := range uis {
		mui := structs.GetMemberMetadata(pi, targetUi)
		if mui.IsUserBanned {
			continue
		}
	}
	if extra.ExtraMessageHandler(pi, ui, evt, msg) {
		go structs.GroupSendToAllMessage(pi, ui.Username, msg.Body, ui.KeyID)
	}
}
