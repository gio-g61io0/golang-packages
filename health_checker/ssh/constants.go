
package ssh

type ERRTYPE string
const (
	DFNRELATED ERRTYPE = "dfn"
	BROKERAGE ERRTYPE = "brokerage"
	TRADER ERRTYPE = "trader"
)

const SEPARATOR = " "
const REPLACEMENT = "|"

type BadKeyword struct{
	keywords []string
	errorType ERRTYPE
}

var BADKEYWORDS = []BadKeyword{BadKeyword{
	keywords: []string{"error", "exception", "dfn", "errors", "rejected", "failed"},
	errorType: DFNRELATED,

}, BadKeyword{
	keywords: []string{"error", "exception", "errors", "failed"},
	errorType: BROKERAGE,

}, BadKeyword{
	keywords: []string{"error", "exception", "trader", "bad", "gateway", "errors", "rejected", "failed"},
	errorType:TRADER,

}}

