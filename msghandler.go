package main

import (
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"strings"
)

func botMsgHandler(pi *core.PrivateInfoS, ui core.UserInfo, evt *core.Event, msg *core.Message) {
	if strings.HasPrefix(msg.Body, "/help") {
		pi.SendMessage(ui, core.MessageTypeText, "Hello!")
	}
}
