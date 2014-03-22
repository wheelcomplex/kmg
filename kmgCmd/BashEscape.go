package kmgCmd

import (
	"strings"
)

//escape string to put it into bash
//may have some problem in extreme case
func BashEscape(inS string) (outS string) {
	return "'" + strings.Replace(inS, "'", "'''", -1) + "'"
}
