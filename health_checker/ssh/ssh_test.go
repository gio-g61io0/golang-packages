package ssh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicate(t *testing.T) {
	testStr := "2026-09-15T07:18:37.204448325Z  INFO:     172.18.0.1:52140"
	testStr = strings.ReplaceAll(testStr, SEPARATOR, REPLACEMENT)
	removedDuplicate := RemoveDuplicateCharWithin(testStr)

	assert.Equal(t, "2026-09-15T07:18:37.204448325Z|INFO:|172.18.0.1:52140", removedDuplicate)
	testStr = "            "
	testStr = strings.ReplaceAll(testStr, SEPARATOR, REPLACEMENT)
	removedDuplicate = RemoveDuplicateCharWithin(testStr)
	assert.Equal(t, "|", removedDuplicate)

}
func TestGetParts(t *testing.T) {
	testStr := "2026-09-15T07:18:37.204448325Z  INFO:     172.18.0.1:52140"
	parts := GetParts(testStr, REPLACEMENT)
	fmt.Printf("Parts: %v\n", parts)
	assert.Equal(t, 3, len(parts))

}
