package pagination

type Pages struct {
	Limit  int64 `validate:"neglect"`
	Page   int64 `validate:"neglect"`
	Offset int64
}

func (p *Pages) Calculate(limit int64) {
	if p.Limit < 1 {
		p.Limit = limit
	}
	if p.Page > 0 {
		p.Offset = p.Limit * (p.Page - 1)
	}
}
