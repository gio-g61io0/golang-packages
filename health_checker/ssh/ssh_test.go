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


	testStr = "       a "
	testStr = strings.ReplaceAll(testStr, SEPARATOR, REPLACEMENT)
	removedDuplicate = RemoveDuplicateCharWithin(testStr)
	assert.Equal(t, "|a|", removedDuplicate)

	testStr = ""
	testStr = strings.ReplaceAll(testStr, SEPARATOR, REPLACEMENT)
	removedDuplicate = RemoveDuplicateCharWithin(testStr)
	assert.Equal(t, "", removedDuplicate)

}
func TestGetParts(t *testing.T) {
	testStr := "2026-09-15T07:18:37.204448325Z  INFO:     172.18.0.1:52140"
	parts := GetParts(testStr, REPLACEMENT)
	fmt.Printf("Parts: %v\n", parts)
	assert.Equal(t, 3, len(parts))

}
func TestParseLine(t *testing.T) {

	testStr := "2026-09-28T04:36:00.253716712Z [2026-09-28 07:36:00,253: INFO/ForkPoolWorker-3] GTN report: trade_order 3716 deferred — 1 member(s) still at the broker (['E6de39a3b06c7']). An error occured"

	parsedError , err := ParseLine(testStr)

	assert.NoError(t, err)
	assert.Contains(t,parsedError.Keywords[DFNRELATED], "error")
	assert.Contains(t,parsedError.Keywords[TRADER], "error")
	assert.Contains(t,parsedError.Keywords[BROKERAGE], "error")

}
