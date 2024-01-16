package container

import "sync"

type SafeMap[K comparable, V any] struct {
	mutex sync.RWMutex
	data  map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	var sm = &SafeMap[K, V]{
		data: make(map[K]V),
	}
	return sm
}

func (sm *SafeMap[K, V]) Get(k K) (V, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	rs, ok := sm.data[k]
	return rs, ok
}

func (sm *SafeMap[K, V]) Set(k K, v V) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.data[k] = v
}

func (sm *SafeMap[K, V]) Delete(k K) V {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	var rs = sm.data[k]
	delete(sm.data, k)
	return rs
}

func (sm *SafeMap[K, V]) ReadEach(cb func(k K, v V)) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for k, v := range sm.data {
		cb(k, v)
	}
}

func (sm *SafeMap[K, V]) WriteEach(cb func(k K, v V)) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for k, v := range sm.data {
		cb(k, v)
	}
}

func (sm *SafeMap[K, V]) Read(k K, cb func(k K, v V)) V {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var rs = sm.data[k]
	cb(k, rs)
	return rs
}

func (sm *SafeMap[K, V]) Write(k K, cb func(k K, v V) V) V {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.data[k] = cb(k, sm.data[k])
	return sm.data[k]
}

func (sm *SafeMap[K, V]) Len() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return len(sm.data)
}

func (sm *SafeMap[K, V]) Clone() *SafeMap[K, V] {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var rs = NewSafeMap[K, V]()
	for k, v := range sm.data {
		rs.data[k] = v
	}
	return rs
}

func (sm *SafeMap[K, V]) AllKey() []K {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var rs = make([]K, 0, len(sm.data))
	for k := range sm.data {
		rs = append(rs, k)
	}
	return rs
}

func (sm *SafeMap[K, V]) AllValue() []V {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var rs = make([]V, 0, len(sm.data))
	for _, v := range sm.data {
		rs = append(rs, v)
	}
	return rs
}

func (sm *SafeMap[K, V]) Clean() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.data = make(map[K]V, 10)
}
