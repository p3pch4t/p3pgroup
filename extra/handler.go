package extra

import (
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"git.mrcyjanek.net/p3pch4t/p3pgroup/extra/ban"
)

// ExtraMessageHandler - The intention of which is to provide
// a modular and extensible way of handling events.
// Things like bridges, plugins, filters - all of that
// go right here.
// NOTE: We may need to add something more - a special
// handler for all extra HTTP requests to the i2p.
// This is mainly for the webhooks (and possibly for
// matterbridge integration).
// false - do not send message
// true - send message
func ExtraMessageHandler(pi *core.PrivateInfoS, ui *core.UserInfo, evt *core.Event, msg *core.Message) bool {
	var relayVars = []bool{}

	relayVars = append(relayVars, ban.Handle(pi, ui, evt, msg))

	for i := range relayVars {
		if !relayVars[i] {
			return false
		}
	}
	return true
}
