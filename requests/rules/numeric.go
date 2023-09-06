package rules

import (
	"math"
)

func initNumericRules() {
	AddTypeRule("multiple_of", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value%arguments.GetInt(0) == 0
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value%arguments.GetUint(0) == 0
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			n := value / arguments.GetFloat(0)
			return n == math.Floor(n)
		},
	})
}
