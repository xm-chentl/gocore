package noticeex

import (
	"errors"
	"sync"
)

const DEFAULT = "default"

var (
	mt            sync.Mutex
	keyOfInstance = make(map[string]INotice)
)

func Default() (INotice, error) {
	inst, ok := keyOfInstance[DEFAULT]
	if !ok {
		return nil, errors.New("default INotice is not initialize")
	}

	return inst, nil
}

func Get(key string) (INotice, error) {
	mt.Lock()
	defer mt.Unlock()

	inst, ok := keyOfInstance[key]
	if !ok {
		return nil, errors.New("key install is not exist")
	}

	return inst, nil
}

func Set(key string, inst INotice) {
	mt.Lock()
	defer mt.Unlock()

	if inst == nil {
		return
	}
	keyOfInstance[key] = inst
}

func SetDefault(inst INotice) {
	mt.Lock()
	defer mt.Unlock()

	if inst == nil {
		return
	}
	keyOfInstance[DEFAULT] = inst
}
