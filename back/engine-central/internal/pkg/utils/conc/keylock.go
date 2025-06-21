package conc

import "sync"

type keyMutex struct {
	localLockMap map[string]*lock
	globalLock   sync.Mutex
}

type lock struct {
	mux      *sync.Mutex
	refCount int
}

// NewKeyMutex creates a new key mutex
// like sync.Mutex but for a given key
func NewKeyMutex() *keyMutex {
	return &keyMutex{localLockMap: map[string]*lock{}}
}

// Lock locks the mutex for the given key
func (km *keyMutex) Lock(key string) {
	km.globalLock.Lock()

	wl, locked := km.localLockMap[key]

	if !locked {
		wl = &lock{
			mux:      new(sync.Mutex),
			refCount: 0,
		}
		km.localLockMap[key] = wl
	}

	wl.refCount++

	km.globalLock.Unlock()

	wl.mux.Lock()
}

// Unlock unlocks the mutex for the given key
func (km *keyMutex) Unlock(key string) {
	km.globalLock.Lock()

	wl, locked := km.localLockMap[key]

	if !locked {
		km.globalLock.Unlock()
		return
	}

	wl.refCount--

	if wl.refCount <= 0 {
		delete(km.localLockMap, key)
	}

	km.globalLock.Unlock()

	wl.mux.Unlock()
}
