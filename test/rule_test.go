package util

import (
	"fmt"
	"goUtil/rule"
	"testing"
)

func TestFn(t *testing.T) {
	r := &rule.Rule{}
	r.SetExpr(`contains(a,"123555")`)
	fmt.Println(r.Eval(map[string]interface{}{
		"a": "123",
	}))
}
