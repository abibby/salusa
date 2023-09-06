package rules

import "time"

func initDateRules() {
	AddTypeRule("after", &TypeRule{
		ArgCount: 1,
		Time: func(value time.Time, arguments TypeRuleArguments) bool {
			return value.After(arguments.GetTime(0))
		},
	})
	AddTypeRule("after_or_equal", &TypeRule{
		ArgCount: 1,
		Time: func(value time.Time, arguments TypeRuleArguments) bool {
			return !value.Before(arguments.GetTime(0))
		},
	})
	AddTypeRule("before", &TypeRule{
		ArgCount: 1,
		Time: func(value time.Time, arguments TypeRuleArguments) bool {
			return value.Before(arguments.GetTime(0))
		},
	})
	AddTypeRule("before_or_equal", &TypeRule{
		ArgCount: 1,
		Time: func(value time.Time, arguments TypeRuleArguments) bool {
			return !value.After(arguments.GetTime(0))
		},
	})
}
