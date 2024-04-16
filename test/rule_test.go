package util

import (
	"fmt"
	"testing"

	"github.com/carmel/go-util/rule"
)

func TestFn(t *testing.T) {
	r := &rule.Rule{}
	r.SetExpr(`contains(a,"123555")`)
	fmt.Println(r.Eval(map[string]interface{}{
		"a": "123",
	}))
}
