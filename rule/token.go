package rule

import (
	"errors"
	"go/token"
	"reflect"

	"github.com/shopspring/decimal"
)

func operate(x, y interface{}, tk token.Token) (interface{}, error) {
	xv := reflect.ValueOf(x)
	yv := reflect.ValueOf(y)
	switch tk {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.LSS, token.GTR, token.LEQ, token.GEQ, token.EQL, token.NEQ:
		var a, b float64
		var err error
		if a, err = number(xv); err != nil {
			return nil, err
		}
		if b, err = number(yv); err != nil {
			return nil, err
		}
		switch tk {
		case token.ADD:
			return decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)), nil
		case token.SUB:
			return decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)), nil
		case token.MUL:
			return decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)), nil
		case token.QUO:
			if b == 0 {
				return 0, errors.New("x/0 error")
			}
			return a / b, nil
		case token.LSS:
			return a < b, nil
		case token.GTR:
			return a > b, nil
		case token.LEQ:
			return a <= b, nil
		case token.GEQ:
			return a >= b, nil
		case token.EQL:
			return a == b, nil
		case token.NEQ:
			return a != b, nil
		default:
			return 0, ErrUnsupportToken
		}
	case token.LAND, token.LOR:
		if xv.Kind() != reflect.Bool || yv.Kind() != reflect.Bool {
			return false, ErrNotBool
		}
		switch tk {
		case token.LAND:
			return xv.Bool() && yv.Bool(), nil
		case token.LOR:
			return xv.Bool() || yv.Bool(), nil
		default:
			return false, ErrUnsupportToken
		}
	default:
		return nil, ErrUnsupportToken
	}
}

func number(x reflect.Value) (float64, error) {
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(x.Int()), nil
	case reflect.Float32, reflect.Float64:
		return x.Float(), nil
	default:
		return 0, ErrNotNumber
	}
}

func evelUnary(x reflect.Value) {

}
