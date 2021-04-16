package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configuration := GetConfig()
	assert.NotNil(t, configuration)
}
