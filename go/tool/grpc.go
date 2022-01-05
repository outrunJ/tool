package tool

import "strings"

func GRPCMethodName(fullMethod string) string {
	method := fullMethod[1:]
	method = method[strings.Index(method, "/")+1:]
	return method
}
