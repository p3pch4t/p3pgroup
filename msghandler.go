package main

import (
	"fmt"
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"github.com/google/shlex"
)

func botMsgHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event, msg *core.Message) {
	// - `!create "Group Name" "Group Email [optional]"` - Create a group
	// - `!list` - Show list of all groups
	// - `!details group_id` - Show details about a group.
	// - `!delete group_id` - Delete a group
	// - `!help` - Show this very message
	cmd, err := shlex.Split(msg.Body)
	if err != nil {
		pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("Unable to process message. Is your syntax valid?.\n\n%s", err.Error()))
		return
	}
	if len(cmd) == 0 {
		pi.SendMessage(ui, core.MessageTypeText, "It looks like your message is empty. Please provide a command. Check `!help` for instructions.")
		return
	}
	switch cmd[0] {
	case "!help":
		pi.SendMessage(ui, core.MessageTypeText, welcomeMsg)
	case "!create":
		if len(cmd) != 2 && len(cmd) != 3 {
			pi.SendMessage(ui, core.MessageTypeText, "Usage: `!create 'Group Name' 'Group Email [optional]'`")
			return
		}
		email := ""
		if len(cmd) == 3 {
			email = cmd[2]
		}
		groupId := createGroup(ui.GetKeyID(), cmd[1], email)
		pi.SendMessage(ui, core.MessageTypeText, "Group Created. ID: `"+groupId+"`")
	case "!list":
		var text = "Groups:\n"
		for _, info := range groups {
			text += fmt.Sprintf("- %s (`%s`)\n", info.GroupName, info.pi.Endpoint)
		}
		pi.SendMessage(ui, core.MessageTypeText, text)
	case "!details":
		if len(cmd) != 2 {
			pi.SendMessage(ui, core.MessageTypeText, "Usage: `!details 'groupId'`")
			return
		}
	default:
		pi.SendMessage(ui, core.MessageTypeText, "Command not found. Check `!help`")

	}
}
