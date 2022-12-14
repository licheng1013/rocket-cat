package core


// Plugin 插件
type Plugin interface {
	Invok(app *App)
}