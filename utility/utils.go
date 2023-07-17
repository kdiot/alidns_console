package utility

func StringPtr(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

func DefaultIfEmpty(value *string, defaultValue *string) *string {
	if value == nil || *value == "" {
		return defaultValue
	} else {
		return value
	}
}

func DefaultIfNull(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	} else {
		return value
	}
}
