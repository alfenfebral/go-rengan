package utils_test

import (
	"testing"

	"go-rengan/utils"

	"github.com/stretchr/testify/assert"
)

// PerPage - get per_page based on query string, the default value is 10
func TestPerPage(t *testing.T) {
	value := utils.PerPage("")
	assert.Equal(t, value, 10)

	value = utils.PerPage("10")
	assert.Equal(t, value, 10)
}

// CurrentPage - get current pages
func TestCurrentPage(t *testing.T) {
	value := utils.CurrentPage("")
	assert.Equal(t, value, 1)

	value = utils.CurrentPage("10")
	assert.Equal(t, value, 10)
}

// TotalPage - get total pages
func TestTotalPage(t *testing.T) {
	value := utils.TotalPage(20, 10)
	assert.Equal(t, value, 2)
}

// Offset - offset of pages
func TestOffset(t *testing.T) {
	value := utils.Offset(1, 10)
	assert.Equal(t, value, 0)

	value = utils.Offset(-1, 10)
	assert.Equal(t, value, 0)
}
