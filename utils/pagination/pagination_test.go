package pagination_test

import (
	"testing"

	paginationutil "go-rengan/utils/pagination"

	"github.com/stretchr/testify/assert"
)

// PerPage - get per_page based on query string, the default value is 10
func TestPerPage(t *testing.T) {
	value := paginationutil.PerPage("")
	assert.Equal(t, value, 10)

	value = paginationutil.PerPage("10")
	assert.Equal(t, value, 10)
}

// CurrentPage - get current pages
func TestCurrentPage(t *testing.T) {
	value := paginationutil.CurrentPage("")
	assert.Equal(t, value, 1)

	value = paginationutil.CurrentPage("10")
	assert.Equal(t, value, 10)
}

// TotalPage - get total pages
func TestTotalPage(t *testing.T) {
	value := paginationutil.TotalPage(20, 10)
	assert.Equal(t, value, 2)
}

// Offset - offset of pages
func TestOffset(t *testing.T) {
	value := paginationutil.Offset(1, 10)
	assert.Equal(t, value, 0)

	value = paginationutil.Offset(-1, 10)
	assert.Equal(t, value, 0)
}
