package rules

func initBoolRules() {
	AddTypeRule("accepted", &TypeRule{
		ArgCount: 0,
		Bool: func(value bool, arguments TypeRuleArguments) bool {
			return value
		},
	})
	AddTypeRule("declined", &TypeRule{
		ArgCount: 0,
		Bool: func(value bool, arguments TypeRuleArguments) bool {
			return !value
		},
	})
}
