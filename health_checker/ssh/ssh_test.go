package ssh

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicate(t *testing.T) {
	testStr := "2026-09-15T07:18:37.204448325Z  INFO:     172.18.0.1:52140"
	testStr = strings.ReplaceAll(testStr, " ", "-")
	removedDuplicate := RemoveDuplicateCharWithin(testStr)

	assert.Equal(t, "2026-09-15T07:18:37.204448325Z-INFO:-172.18.0.1:52140", removedDuplicate)
	testStr = "            "
	testStr = strings.ReplaceAll(testStr, " ", "-")
	removedDuplicate = RemoveDuplicateCharWithin(testStr)
	assert.Equal(t, "-", removedDuplicate)

}
