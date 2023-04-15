package common

import "sync"

type MatchQueue struct {
	// 匹配人数多少回调
	maxMatch int64
	// 匹配成功回调,参数是匹配成功的玩家id
	callback func(players []int64)
	// 互斥锁
	mutex sync.Mutex
	// id 和 player
	players map[int64]IPlayer
	// 加入顺序
	matches []int64
}

func NewMatchQueue(maxMatch int64, callback func([]int64)) *MatchQueue {
	return &MatchQueue{
		maxMatch: maxMatch,
		players:  make(map[int64]IPlayer, 0),
		callback: callback,
	}
}

func (mq *MatchQueue) AddMatch(player IPlayer) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	// 如果已经在匹配队列中了，就不再添加
	if _, ok := mq.players[player.UserId()]; ok {
		return
	}
	mq.players[player.UserId()] = player
	mq.matches = append(mq.matches, player.UserId())
	if int64(len(mq.players)) >= mq.maxMatch { // 匹配成功
		mq.callback(mq.matches)
		mq.matches = make([]int64, 0)
		mq.players = make(map[int64]IPlayer, 0)
	}
}

func (mq *MatchQueue) RemoveMatch(match int64) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	delete(mq.players, match)
	for i, v := range mq.matches {
		if v == match {
			mq.matches = append(mq.matches[:i], mq.matches[i+1:]...)
			break
		}
	}
}

func (mq *MatchQueue) GetMatches() []int64 {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	return mq.matches
}
