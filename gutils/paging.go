package gutils

type Paging struct {
	Skip  int64
	Limit int64
	Total int64
}

func (p *Paging) FullFil() {
	if p.Skip < 0 {
		p.Skip = 0
	}
	if p.Limit <= 0 || p.Limit > 30 {
		p.Limit = 10
	}
}
