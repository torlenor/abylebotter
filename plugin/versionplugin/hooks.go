package versionplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/abylebotter/model"
)

// OnPost implements the hook from the Bot
func (p *VersionPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")
	if strings.HasPrefix(msg, "!version") {
		versionPost := post
		versionPost.Content = p.API.GetVersion()
		p.API.LogTrace(fmt.Sprintf("Echoing version back to Channel = %s, content = %s", versionPost.Channel, versionPost.Content))
		p.API.CreatePost(versionPost)
	}
}
