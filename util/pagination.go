package util

type Pagination struct {
	Page     int64
	PageSize int64
}

func (p *Pagination) Normalize() {
	p.Page = max(p.Page, 1)

	if p.PageSize <= 0 || p.PageSize >= 20 {
		p.PageSize = 20
	}
}
