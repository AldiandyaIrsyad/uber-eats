package models

type Pagination struct {
	Data     []interface{} `json:"data,omitempty"`
	Page     int64         `json:"page" validate:"required,min=1"`
	PageSize int64         `json:"pageSize" validate:"required,min=1,max=100"`
	Total    int64         `json:"total,omitempty"`
}

// NewPagination creates a new Pagination with default values
func NewPagination(page, pageSize int64) *Pagination {
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Validate checks if pagination values are valid
func (p *Pagination) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return nil
}

// GetSkip returns the number of documents to skip
func (p *Pagination) GetSkip() int64 {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the number of documents to return
func (p *Pagination) GetLimit() int64 {
	return p.PageSize
}

// SetTotal sets the total number of documents
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}

// GetTotalPages returns the total number of pages
func (p *Pagination) GetTotalPages() int64 {
	if p.Total == 0 {
		return 0
	}
	totalPages := p.Total / p.PageSize
	if p.Total%p.PageSize > 0 {
		totalPages++
	}
	return totalPages
}
