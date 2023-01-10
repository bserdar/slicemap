package slicemap

// A SliceMap is a map[[]K]T, that is, a map that uses a slice of K as
// its key. The map implementation is based on nested map of type
// map[K]
type SliceMap[K comparable, T any] struct {
	m        map[K]*SliceMap[K, T]
	value    T
	hasValue bool
	cnt      int
}

// Get returns the value stored for the key
func (s *SliceMap[K, T]) Get(key []K) (val T, exists bool) {
	if len(key) == 0 {
		if s.hasValue {
			return s.value, true
		}
		return
	}
	if s.m == nil {
		return
	}
	next, ok := s.m[key[0]]
	if !ok {
		return
	}
	return next.Get(key[1:])
}

// Put sets the value of key to value, and returns the old value if
// there was one, and whether or not a value was replaced.
func (s *SliceMap[K, T]) Put(key []K, value T) (oldVal T, replaced bool) {
	if len(key) == 0 {
		if s.hasValue {
			oldVal = s.value
			replaced = true
		} else {
			s.cnt++
		}
		s.value = value
		s.hasValue = true
		return
	}

	if s.m == nil {
		s.m = make(map[K]*SliceMap[K, T])
		newMap := &SliceMap[K, T]{}
		s.m[key[0]] = newMap
		oldVal, replaced = newMap.Put(key[1:], value)
		if !replaced {
			s.cnt++
		}
		return
	}

	next, exists := s.m[key[0]]
	if !exists {
		next = &SliceMap[K, T]{}
		s.m[key[0]] = next
	}
	oldVal, replaced = next.Put(key[1:], value)
	if !replaced {
		s.cnt++
	}
	return
}

// Delete the key and return if the key is deleted
func (s *SliceMap[K, T]) Delete(key []K) bool {
	if len(key) == 0 {
		if s.hasValue {
			s.hasValue = false
			s.cnt--
			return true
		}
		return false
	}
	if s.m == nil {
		return false
	}
	next, exists := s.m[key[0]]
	if !exists {
		return false
	}
	exists = next.Delete(key[1:])
	if exists {
		s.cnt--
	}
	if len(next.m) == 0 && !next.hasValue {
		delete(s.m, key[0])
	}
	return exists
}

// ForEach calls f for each element in the map until f returns false.
func (s *SliceMap[K, T]) ForEach(f func([]K, T) bool) bool {
	var zero K
	var forEach func([]K, *SliceMap[K, T]) bool
	forEach = func(key []K, mp *SliceMap[K, T]) bool {
		if mp.hasValue {
			if !f(key, mp.value) {
				return false
			}
		}
		if mp.m == nil {
			return true
		}
		nextKey := key
		ix := len(nextKey)
		nextKey = append(nextKey, zero)
		for k, v := range mp.m {
			nextKey[ix] = k
			if !forEach(nextKey, v) {
				return false
			}
		}
		return true
	}
	keyArr := make([]K, 0, 16)
	return forEach(keyArr, s)
}

// Len returns the number of elements in the map
func (s *SliceMap[K, T]) Len() int { return s.cnt }
