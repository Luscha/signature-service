package logger

import "go.uber.org/zap/zapcore"

type option struct {
	name      string
	service   string
	tracingId string
	subject   string
	minLevel  Level
}

type OptionFunc func(*option)

func defaultOptions() *option {
	return &option{
		minLevel: Level(zapcore.DebugLevel),
	}
}

func WithName(name string) OptionFunc {
	return func(o *option) {
		o.name = name
	}
}

func WithTracingId(tracingId string) OptionFunc {
	return func(o *option) {
		o.tracingId = tracingId
	}
}

func WithSubject(subject string) OptionFunc {
	return func(o *option) {
		o.subject = subject
	}
}

func WithService(service string) OptionFunc {
	return func(o *option) {
		o.service = service
	}
}

func WithMinLevel(minLevel string) OptionFunc {
	return func(o *option) {
		switch minLevel {
		case "DEBUG":
			o.minLevel = Level(zapcore.DebugLevel)
		case "WARN":
			o.minLevel = Level(zapcore.WarnLevel)
		case "ERROR":
			o.minLevel = Level(zapcore.ErrorLevel)
		default:
			o.minLevel = Level(zapcore.InfoLevel)
		}
	}
}
