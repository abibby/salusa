package rules

// Nullable
// Size

// Required
// Required If
// Required Unless
// Required With
// Required With All
// Required Without
// Required Without All
// Required Array Keys
// Sometimes

func initGenericRules() {
	AddTypeRule("gt", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value > arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value > arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value > arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return value > arguments.GetString(0)
		},
	})
	AddTypeRule("gte", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return value >= arguments.GetString(0)
		},
	})
	AddTypeRule("lt", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value < arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value < arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value < arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return value < arguments.GetString(0)
		},
	})
	AddTypeRule("lte", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return value <= arguments.GetString(0)
		},
	})
	AddTypeRule("max", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value <= arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return len(value) <= int(arguments.GetInt(0))
		},
	})
	AddTypeRule("min", &TypeRule{
		ArgCount: 1,
		Int: func(value int64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetInt(0)
		},
		Uint: func(value uint64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetUint(0)
		},
		Float: func(value float64, arguments TypeRuleArguments) bool {
			return value >= arguments.GetFloat(0)
		},
		String: func(value string, arguments TypeRuleArguments) bool {
			return len(value) >= int(arguments.GetInt(0))
		},
	})
}
