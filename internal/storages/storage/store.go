package storage

type Store struct {
	ID       uint64
	Lifetime uint64
	Name     string
	Code     string
	Groups   map[uint64]struct{}
}

func (v Store) GetGroups() []uint64 {
	result := make([]uint64, 0, len(v.Groups))
	for g := range v.Groups {
		result = append(result, g)
	}
	return result
}
