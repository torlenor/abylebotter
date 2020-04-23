package quotesplugin

import (
	"strings"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

// OnPost implements the hook from the Bot
func (p *QuotesPlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")
	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if strings.HasPrefix(msg, "!quote ") || msg == "!quote" {
			p.onCommandQuote(post)
			return
		} else if strings.HasPrefix(msg, "!quoteadd ") {
			p.onCommandQuoteAdd(post)
			return
		} else if strings.HasPrefix(msg, "!quoteremove ") {
			p.onCommandQuoteRemove(post)
			return
		} else if msg == "!quoteadd" || msg == "!quotehelp" {
			p.returnHelp(post.ChannelID)
			return
		} else if msg == "!quoteremove" {
			p.returnHelpRemove(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
