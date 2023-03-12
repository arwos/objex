package storage

type Store struct {
	ID       int64
	Lifetime int64
	Name     string
	Code     string
	Groups   map[int64]struct{}
}

func (v Store) GetGroups() []int64 {
	result := make([]int64, 0, len(v.Groups))
	for g := range v.Groups {
		result = append(result, g)
	}
	return result
}
