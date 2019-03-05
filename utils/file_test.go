package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckFolder(t *testing.T) {
	err := CheckFolder("/Users/cookie/go/gopath/src/image-server/asb/ed/sdfa/asd/g/asg/asd/ga/")
	assert.Nil(t, err)
}
