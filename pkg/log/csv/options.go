package gclogcsv

type LogOption func(any) error

func WithDelimiter(delim rune) LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.cw.comma = string(delim)
		return nil
	}
}

func WithTruncateOnHeadersMissing() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.truncateOnHeadersMissing = true
		return nil
	}
}

func WithDisableTimestamp() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.disableTimestamp = true
		return nil
	}
}

func WithEnableLevel() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableLevel = true
		return nil
	}
}

func WithEnableCaller() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableCaller = true
		return nil
	}
}

func WithEnableError() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableError = true
		return nil
	}
}
