package functions

import "text/template"

// https://golang.org/pkg/text/template/

var Templatefuncs = template.FuncMap{
	"translated": translated,
	"account":    account,
	"tax":        tax,
}
