package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

// 错误定义
var (
	ErrRuleEmpty      = errors.New("rule is empty")
	ErrUnsupportToken = errors.New("unsupport token")
	ErrUnsupportExpr  = errors.New("unsupport expr")
	ErrUnsupportParam = errors.New("unsupport param")
	ErrNotNumber      = errors.New("not a number")
	ErrIndexNotNumber = errors.New("index not a number")
	ErrNotBool        = errors.New("not boolean")
	ErrKeyNotFound    = errors.New("map key not found")
	fns               = map[string]fn{
		"contains": func(args []interface{}) (interface{}, error) {
			fmt.Println("contains print: ", args)
			return nil, nil
		},
	}
)

type fn = func(args []interface{}) (interface{}, error)

type Rule struct {
	expr ast.Expr
}

func (r *Rule) SetExpr(expr string) error {
	if len(expr) == 0 {
		return ErrRuleEmpty
	}
	if exp, err := parser.ParseExpr(expr); err != nil {
		return err
	} else {
		r.expr = exp
	}

	return nil
}

func (r *Rule) Bool(database map[string]interface{}) (bool, error) {
	if r.expr != nil {
		b, err := r.Eval(database)
		if err != nil {
			return false, err
		}
		if r, ok := b.(bool); ok {
			return r, nil
		}
	}

	return false, errors.New("expr is nil")
}

func (r *Rule) Int(database map[string]interface{}) (int64, error) {
	if r.expr != nil {
		b, err := r.Eval(database)
		if err != nil {
			return 0, err
		}
		switch b := b.(type) {
		case int64:
			return b, nil
		case float64:
			return int64(b), nil

		}
	}
	return 0, errors.New("expr is nil")

}

func (r *Rule) Float(database map[string]interface{}) (float64, error) {
	if r.expr != nil {
		b, err := r.Eval(database)
		if err != nil {
			return 0, err
		}
		switch b := b.(type) {
		case int64:
			return float64(b), nil
		case float64:
			return b, nil

		}
	}
	return 0, errors.New("expr is nil")
}

func (r *Rule) Eval(datasource map[string]interface{}) (interface{}, error) {
	switch t := r.expr.(type) {
	case *ast.UnaryExpr: // 一元表达式
		r.expr = t.X
		operand, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}
		oprd := reflect.ValueOf(operand)

		switch t.Op {
		case token.NOT: // !
			if oprd.Kind() != reflect.Bool {
				return false, ErrNotBool
			}
			return !oprd.Bool(), nil
		case token.SUB: // -
			if x, err := number(oprd); err == nil {
				return (-1.0) * x, nil
			}
			return 0.0, ErrNotNumber
		}
	case *ast.BinaryExpr: // 二元表达式
		r.expr = t.X
		x, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}
		r.expr = t.Y
		y, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}
		return operate(x, y, t.Op)
	case *ast.Ident: // 标志符（已定义变量或常量（bool））
		return evalIdent(t.Name, datasource)
	case *ast.BasicLit: // 基本类型文字（当作字符串存储）
		switch t.Kind {
		case token.STRING:
			return strings.Trim(t.Value, "\""), nil
		case token.INT:
			return strconv.ParseInt(t.Value, 10, 64)
		case token.FLOAT:
			return strconv.ParseFloat(t.Value, 64)
		default:
			return nil, ErrUnsupportParam
		}
	case *ast.ParenExpr: // 圆括号内表达式
		r.expr = t.X
		return r.Eval(datasource)
	case *ast.SelectorExpr: // 属性或方法选择表达式
		r.expr = t.X
		v, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}
		return evalIdent(t.Sel.Name, v.(map[string]interface{}))
	case *ast.IndexExpr: // 中括号内表达式——map或slice索引
		r.expr = t.X
		data, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}

		r.expr = t.Index
		idx, err := r.Eval(datasource)
		if err != nil {
			return nil, err
		}

		switch data := data.(type) {
		case map[string]interface{}:
			if idx, isString := idx.(string); isString {
				return data[idx], nil
			} else {
				return nil, fmt.Errorf("map here index must be string")
			}
		case []interface{}:
			switch idx := idx.(type) {
			case int:
				return data[int64(idx)], nil
			case int64:
				return data[idx], nil
			default:
				return nil, fmt.Errorf("slice index index must be number")
			}
		default:
			return nil, fmt.Errorf("IndexExpr: unsupport data type")
		}
	case *ast.CallExpr: // 方法调用表达式
		// r.expr = t.Fun
		// f, err := r.Eval(fns)

		// if err != nil {
		// 	return nil, err
		// }
		// switch f := f.(type) {
		// case func(args []interface{}) (interface{}, error):
		// 	if params, err := evalArg(t.Args, datasource); err != nil {
		// 		return nil, err
		// 	} else {
		// 		return f(params)
		// 	}
		// }
		if params, err := evalArg(t.Args, datasource); err != nil {
			return nil, err
		} else {
			return fns[t.Fun.(*ast.Ident).Name](params)
		}
	}
	return nil, ErrUnsupportExpr
}

func evalArg(args []ast.Expr, datasource map[string]interface{}) ([]interface{}, error) {
	var result []interface{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case *ast.BasicLit:
			result = append(result, arg.Value)
		case *ast.Ident:
			if val, err := evalIdent(arg.Name, datasource); err != nil {
				return nil, err
			} else {
				result = append(result, val)
			}
		}
	}
	return result, nil
}

func evalIdent(key string, datasource map[string]interface{}) (interface{}, error) {
	// while bool type is Ident
	if key == "true" {
		return true, nil
	} else if key == "false" {
		return false, nil
	}

	if value, ok := datasource[key]; ok {
		return value, nil
	} else {
		return nil, ErrKeyNotFound
	}
}
