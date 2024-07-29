package rules

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type ValidationOptions struct {
	Value     any
	Arguments []string
	Request   *http.Request
	Name      string
}

type ValidationRule func(options *ValidationOptions) bool

var typeTime = reflect.TypeFor[time.Time]()

var rules = map[string]ValidationRule{}

func AddRule(key string, rule ValidationRule) {
	rules[key] = rule
}

var initRules = sync.OnceFunc(func() {
	initNumericRules()
	initStringRules()
	initGenericRules()
	initBoolRules()
	initDateRules()
})

func GetRule(key string) (ValidationRule, bool) {
	initRules()

	r, ok := rules[key]
	return r, ok
}

type TypeRuleArguments []string

func (a TypeRuleArguments) GetString(i int) string {
	if len(a) < i {
		return ""
	}
	return a[i]
}

func (a TypeRuleArguments) GetInt(i int) int64 {
	s := a.GetString(i)
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (a TypeRuleArguments) GetUint(i int) uint64 {
	s := a.GetString(i)
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}

func (a TypeRuleArguments) GetFloat(i int) float64 {
	s := a.GetString(i)
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func (a TypeRuleArguments) GetTime(i int) time.Time {
	s := a.GetString(i)
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func (a TypeRuleArguments) GetBoolean(i int) bool {
	s := a.GetString(i)
	return s == "1" || s == "true" || s == "yes"
}

type TypeRule struct {
	ArgCount int
	Int      func(value int64, arguments TypeRuleArguments) bool
	Uint     func(value uint64, arguments TypeRuleArguments) bool
	Float    func(value float64, arguments TypeRuleArguments) bool
	String   func(value string, arguments TypeRuleArguments) bool
	Bool     func(value bool, arguments TypeRuleArguments) bool
	Time     func(value time.Time, arguments TypeRuleArguments) bool
	Array    func(value reflect.Value, arguments TypeRuleArguments) bool
}

func AddTypeRule(key string, rule *TypeRule) {
	AddRule(key, func(options *ValidationOptions) bool {
		if len(options.Arguments) < rule.ArgCount {
			log.Printf("%s must have %d argument(s)", key, rule.ArgCount)
			return true
		}
		val := reflect.ValueOf(options.Value)

		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		if val.Type() == typeTime {
			if rule.Time == nil {
				log.Printf("no rule for int fields")
				return true
			}
			return rule.Time(val.Interface().(time.Time), options.Arguments)
		}

		if val.CanInt() {
			if rule.Int == nil {
				log.Printf("no rule for int fields")
				return true
			}
			return rule.Int(val.Int(), options.Arguments)
		}
		if val.CanUint() {
			if rule.Uint == nil {
				log.Printf("no rule for uint fields")
				return true
			}
			return rule.Uint(val.Uint(), options.Arguments)
		}
		if val.CanFloat() {
			if rule.Float == nil {
				log.Printf("no rule for float fields")
				return true
			}
			return rule.Float(val.Float(), options.Arguments)
		}
		if val.Kind() == reflect.String {
			if rule.String == nil {
				log.Printf("no rule for string fields")
				return true
			}
			return rule.String(val.String(), options.Arguments)
		}

		if val.Kind() == reflect.Slice {
			if rule.Array == nil {
				log.Printf("no rule for string fields")
				return true
			}
			return rule.Array(val, options.Arguments)
		}

		log.Printf("using a numeric rule on a non numeric field")
		return true
	})
}
