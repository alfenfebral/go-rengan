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
	return int(math.Ceil(float64(total) / float64(perPage)))
}

// Offset - offset of pages
func Offset(currentPage int, perPage int) int {
	result := (currentPage - 1) * perPage
	if result < 0 {
		return 0
	}
	return result
}
