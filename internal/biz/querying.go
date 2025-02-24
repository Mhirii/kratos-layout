package biz

type PaginationParams struct {
	Page     int
	PageSize int
	Reverse  bool
}
type SortParams struct {
	SortBy    string
	SortOrder string
}
