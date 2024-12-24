package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystrayIconResource(t *testing.T) {
	resource := GetSystrayIconResource()
	assert.NotNil(t, resource)
	assert.Equal(t, "favicon.ico", resource.Name())
	assert.NotEmpty(t, resource.Content())
}

func TestGetAppIconResource(t *testing.T) {
	resource := GetAppIconResource()
	assert.NotNil(t, resource)
	assert.Equal(t, "icon.png", resource.Name())
	assert.NotEmpty(t, resource.Content())
}
