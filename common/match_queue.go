package common

import "sync"

type MatchQueue struct {
	maxMatch int64
	matches  []int64
	callback func([]int64)
	mutex    sync.Mutex
}

func NewMatchQueue(maxMatch int64, callback func([]int64)) *MatchQueue {
	return &MatchQueue{
		maxMatch: maxMatch,
		matches:  make([]int64, 0),
		callback: callback,
	}
}

func (mq *MatchQueue) AddMatch(match int64) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	if int64(len(mq.matches)) < mq.maxMatch {
		for _, m := range mq.matches {
			if m == match {
				return
			}
		}
		mq.matches = append(mq.matches, match)
	}
	if int64(len(mq.matches)) == mq.maxMatch {
		mq.callback(mq.matches)
		mq.matches = make([]int64, 0)
	}
}

func (mq *MatchQueue) RemoveMatch(match int64) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	for i, m := range mq.matches {
		if m == match {
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
