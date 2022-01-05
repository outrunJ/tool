package sqler

import (
	"git.meiqia.com/business_platform/tool"
	"strings"
	"fmt"
)

func Quote(quote string, ss ...interface{}) *[]string {
	retSS, _ := tool.SliceMap_TInterfaceRTString(&ss, func(v interface{}) (string, error) {
		return quote + fmt.Sprintf("%v", v) + quote, nil
	})
	return retSS
}

func QuoteComma(quote string, ss ...interface{}) string {
	return strings.Join(*Quote(quote, ss...), ",")
}

func Pattern(length int) string {
	return strings.Join(*tool.SliceRepeat_TString("?", length), " , ")
}
