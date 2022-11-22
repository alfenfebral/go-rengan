package timeutil_test

import (
	timeutil "go-rengan/utils/time"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	value := timeutil.GetTimeNow()
	assert.Equal(t, value, value)
}
