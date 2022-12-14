package core


// PluginFunc 每次链接时的插件
type PluginFunc interface {
	Invok(app *App)
}