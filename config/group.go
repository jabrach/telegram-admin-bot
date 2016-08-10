package config

type jGroup struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Muted    []int64 `json:"muted"`
	NoImages []int64 `json:"img_muted`
}

type Group struct {
	ID       int64
	Name     string
	Muted    map[int64]bool
	NoImages map[int64]bool
}

func (j *jGroup) Init() *Group {
	g := &Group{
		ID:       j.ID,
		Name:     j.Name,
		Muted:    map[int64]bool{},
		NoImages: map[int64]bool{},
	}
	for _, id := range j.Muted {
		g.Muted[id] = true
	}
	for _, id := range j.NoImages {
		g.NoImages[id] = true
	}

	return g
}
