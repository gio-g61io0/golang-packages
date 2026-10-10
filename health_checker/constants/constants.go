package constants

type ERRTYPE string
const (
	DFNRELATED ERRTYPE = "dfn"
	BROKERAGE ERRTYPE = "brokerage"
	TRADER ERRTYPE = "trader"
)

const SEPARATOR = " "
const REPLACEMENT = "|"

type BadKeyword struct{
	Keywords []string
	ErrorType ERRTYPE
}

var BADKEYWORDS = []BadKeyword{{
	Keywords: []string{"error", "exception", "dfn", "errors", "rejected", "failed"},
	ErrorType: DFNRELATED,

}, {
	Keywords: []string{"error", "exception", "errors", "failed"},
	ErrorType: BROKERAGE,

}, {
	Keywords: []string{"error", "exception", "trader", "bad", "gateway", "errors", "rejected", "failed"},
	ErrorType:TRADER,

}}

func FormatMapToString(keywords map[ERRTYPE][]string) string{

			keywordByteFormattted := []byte{}

			ForEach(keywords, func(e constants.ERRTYPE, s []string) {

				fmt.Appendf(keywordByteFormattted, "*%s\n", e)
				for idx, keyword := range s{
					fmt.Appendf(keywordByteFormattted, "\t%d. %s\n", idx + 1, keyword)
				}
			})
}
