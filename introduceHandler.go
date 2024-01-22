package main

import (
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"git.mrcyjanek.net/p3pch4t/p3pgroup/structs"
)
import _ "embed"

//go:embed "assets/welcomeMsg.md"
var welcomeMsg string

func botIntroduceHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event) {
	var gUi structs.P3PGROUP_UserInfo
	keyID := ui.GetKeyID()
	pi.DB.First(&gUi, "key_id = ?", keyID)
	gUi.KeyID = keyID // in case of new insert
	if !gUi.IsIntroduced {
		gUi.IsIntroduced = true
		pi.SendMessage(ui, core.MessageTypeText, welcomeMsg)
	}
	pi.DB.Save(&gUi)
}
