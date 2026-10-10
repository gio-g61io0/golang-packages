package utils

import (
	"fmt"
	"personal-http-server/health_checker/constants"
)


//ForEach map's (key, value) it runs the passed callback function
func ForEach(iter map[constants.ERRTYPE][]string, cb func (constants.ERRTYPE, []string)){
	for key, value := range iter{
		cb(key, value)
	}
}

//Formats Detected (Error Type, Keywords) from map to string
func FormatMapToString(keywords map[constants.ERRTYPE][]string) string{

	keywordByteFormattted := []byte{}

	ForEach(keywords, func(e constants.ERRTYPE, s []string) {

		_ = fmt.Appendf(keywordByteFormattted, "*%s\n", e)
		for idx, keyword := range s{
			_ = fmt.Appendf(keywordByteFormattted, "\t%d. %s\n", idx + 1, keyword)
		}
	})
	return string(keywordByteFormattted)
}

