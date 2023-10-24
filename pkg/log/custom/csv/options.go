package gccsvlog

type CsvOption func(*csvCfg) error

func WithTruncateOnHeadersMissing() CsvOption {
	return func(ch *csvCfg) error {
		ch.truncateOnHeadersMissing = true
		return nil
	}
}
