package models

type PaginationOptions struct {
	Page     int64
	PageSize int64
}

type SortOptions struct {
	Field string
	Order int // 1 for ascending, -1 for descending
}

type QueryOptions struct {
	Pagination *PaginationOptions
	Sort       *SortOptions
	Filter     map[string]interface{}
}
