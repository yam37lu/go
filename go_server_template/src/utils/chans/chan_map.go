package chans

import (
	"sync"
)

type ChanMap interface {
	Get(k string) chan bool
	Set(k string, v chan bool)
}

type chanMap struct {
	Data map[string]chan bool
	Lock sync.RWMutex
}

func (cm *chanMap) Get(k string) chan bool {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	if _, ok := cm.Data[k]; ok {
		return cm.Data[k]
	} else {
		return nil
	}
}
func (cm *chanMap) Set(k string, v chan bool) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	cm.Data[k] = v
}

func NewChanMap() ChanMap {
	return &chanMap{Data: make(map[string]chan bool)}
}
