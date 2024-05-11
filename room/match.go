package room

import "sync"

type Match struct {
	// 匹配人数多少回调
	maxMatch int64
	// 匹配成功回调,参数是匹配成功的玩家id
	callback func(players []IPlayer)
	// userId : player
	players sync.Map
}

func NewMatchQueue(maxMatch int64, callback func(players []IPlayer)) *Match {
	return &Match{
		maxMatch: maxMatch,
		callback: callback,
	}
}

func (m *Match) AddMatch(player IPlayer) {
	m.players.Store(player.UserId(), player)
	size := 0
	list := make([]IPlayer, 0)
	m.players.Range(func(key, value any) bool {
		if size > 1 {
			return false
		}
		list = append(list, value.(IPlayer))
		size++
		return true
	})
	if int64(size) >= m.maxMatch { // 匹配成功
		m.callback(list)
		for _, v := range list {
			m.players.Delete(v.UserId())
		}
	}
}

func (m *Match) RemoveMatch(player IPlayer) {
	m.players.Delete(player.UserId())
}

func (m *Match) GetPlayer() []IPlayer {
	list := make([]IPlayer, 0)
	m.players.Range(func(key, value any) bool {
		list = append(list, value.(IPlayer))
		return true
	})
	return list
}
