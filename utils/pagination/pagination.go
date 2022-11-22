package pagination

import (
	"math"
	"strconv"
)

// PerPage - get per_page based on query string, the default value is 10
func PerPage(value string) int {
	if value == "" {
		return 10
	}

	perPage, _ := strconv.Atoi(value)

	return perPage
}

// CurrentPage - get current pages
func CurrentPage(value string) int {
	if value == "" {
		return 1
	}
	perPage, _ := strconv.Atoi(value)

	return perPage
}

// TotalPage - get total pages
func TotalPage(total int, perPage int) int {
	totalFloat := float64(total)
	perPageFloat := float64(perPage)
	resultFloat := math.Ceil(totalFloat / perPageFloat)
	result := int(resultFloat)

	return result
}

// Offset - offset of pages
func Offset(currentPage int, perPage int) int {
	result := (currentPage - 1) * perPage
	if result < 0 {
		return 0
	}

	return result
}
