package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Main *Logger

func init() {
	Main = New()
}

func NewMainLogger(opts ...OptionFunc) {
	Main = New(opts...)
}

type Logger struct {
	*zap.Logger
}

func newLogger(logger *zap.Logger) *Logger {
	return &Logger{
		Logger: logger,
	}
}

func New(opts ...OptionFunc) *Logger {
	o := defaultOptions()
	for _, option := range opts {
		option(o)
	}
	initialFields := map[string]any{}

	if o.subject != "" {
		initialFields[SUBJECT_KEY] = o.subject
	}
	if o.tracingId != "" {
		initialFields[TRACING_KEY] = o.tracingId
	}
	if o.service != "" {
		initialFields[SERVICE_KEY] = o.service
	}

	initialFields[LOGGER_GROUP_KEY] = uuid.New().String()

	return NewLogger(o.name, o.minLevel, initialFields)
}

func NewLogger(name string, minLevel Level, initialFields map[string]interface{}) *Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	cfg.InitialFields = initialFields
	cfg.DisableCaller = true
	cfg.EncoderConfig.StacktraceKey = STACKTRACE_KEY
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.NameKey = "name"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(minLevel))

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	if name != "" {
		logger = logger.Named(name)
	}

	return newLogger(logger)
}

func (l *Logger) WithStack(stack string) *Logger {
	return newLogger(
		l.WithOptions(
			zap.AddStacktrace(zapcore.InvalidLevel),
		).With(zap.String(STACKTRACE_KEY, stack)),
	)
}

func (l *Logger) GetSubLogger(name string) *Logger {
	return newLogger(l.WithFields(map[string]interface{}{
		LOGGER_GROUP_KEY: uuid.New().String(),
	}).Named(name))
}

func (l *Logger) WithFields(args map[string]any) *Logger {
	var newL *zap.Logger = l.Logger
	for k, v := range args {
		newL = newL.With(zap.Reflect(k, v))
	}
	return newLogger(newL)
}

func (l *Logger) WithError(err error) *Logger {
	return newLogger(l.With(zap.Error(err)))
}
