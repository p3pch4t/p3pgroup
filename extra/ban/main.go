package ban

import (
	"fmt"
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"github.com/google/shlex"
	"strings"
)

func Handle(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event, msg *core.Message) bool {
	if !strings.HasPrefix(msg.Body, "!ban") {
		return true
	}
	cmd, err := shlex.Split(msg.Body)
	if err != nil {
		pi.SendMessage(ui, core.MessageTypeText, fmt.Sprintf("Unable to process message. Is your syntax valid?.\n\n%s", err.Error()))
		return false
	}
	pi.SendMessage(ui, core.MessageTypeText, "this is test: ['"+strings.Join(cmd, "', '")+"']")
	return false
}
