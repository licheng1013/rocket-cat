package room

// IRoom 房间接口, 其默认实现为DefaultRoom, 所以继承DefaultRoom使用即可
type IRoom interface {
	// 获取房间id
	GetId() int64
	//  获取房间状态
	GetState() State
	//  获取房间内所有玩家Id
	GetUserIds() []int64
	//  获取房间内所有玩家
	GetPlayers() []IPlayer
	// JoinRoom 加入房间，
	JoinRoom(player IPlayer) bool
	// QuitRoom 退出房间，
	QuitRoom(player IPlayer)
	// GetPlayer 获取某个玩家
	GetPlayer(userId int64) IPlayer
	// 获取玩家数量
	GetPlayerTotal() int
	// 房间更新时间
	GetUpdateTime() int64
}
