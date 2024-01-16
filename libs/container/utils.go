package container

func MergeMap[K comparable, V any](ms ...map[K]V) map[K]V {
	var rs = make(map[K]V)
	for _, itMap := range ms {
		for k, v := range itMap {
			rs[k] = v
		}
	}
	return rs
}
