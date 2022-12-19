package core

import "core/plugins"

// Plugin 插件
type Plugin interface {
	Invok(meta plugins.Meta)
}