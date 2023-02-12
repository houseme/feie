/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

package feie

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILogger is the interface for logger
type ILogger interface {
	Print(ctx context.Context, v ...interface{})
	Printf(ctx context.Context, format string, v ...interface{})
	Debug(ctx context.Context, v ...interface{})
	Debugf(ctx context.Context, format string, v ...interface{})
	Info(ctx context.Context, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})
	Notice(ctx context.Context, v ...interface{})
	Noticef(ctx context.Context, format string, v ...interface{})
	Warning(ctx context.Context, v ...interface{})
	Warningf(ctx context.Context, format string, v ...interface{})
	Error(ctx context.Context, v ...interface{})
	Errorf(ctx context.Context, format string, v ...interface{})
	Critical(ctx context.Context, v ...interface{})
	Criticalf(ctx context.Context, format string, v ...interface{})
	Panic(ctx context.Context, v ...interface{})
	Panicf(ctx context.Context, format string, v ...interface{})
	Fatal(ctx context.Context, v ...interface{})
	Fatalf(ctx context.Context, format string, v ...interface{})
}

// Logger is the interface for logger
type Logger struct {
	logger *zap.Logger
}

// LoggerOption configures the test logger built by NewLogger.
type LoggerOption interface {
	applyLoggerOption(*loggerOptions)
}

type loggerOptions struct {
	Level      zapcore.LevelEnabler
	zapOptions []zap.Option
}

type loggerOptionFunc func(*loggerOptions)

func (f loggerOptionFunc) applyLoggerOption(opts *loggerOptions) {
	f(opts)
}

// Level controls which messages are logged by a test Logger built by
// NewLogger.
func Level(enab zapcore.LevelEnabler) LoggerOption {
	return loggerOptionFunc(func(opts *loggerOptions) {
		opts.Level = enab
	})
}

// WrapOptions adds zap.Option's to a test Logger built by NewLogger.
func WrapOptions(zapOpts ...zap.Option) LoggerOption {
	return loggerOptionFunc(func(opts *loggerOptions) {
		opts.zapOptions = zapOpts
	})
}

// NewLogger is the interface for new logger
func NewLogger(opts ...LoggerOption) *Logger {
	cfg := loggerOptions{
		Level: zapcore.DebugLevel,
	}
	for _, o := range opts {
		o.applyLoggerOption(&cfg)
	}

	// writer := newTestingWriter(t)
	zapOptions := []zap.Option{
		// Send zap errors to the same writer and mark the test as failed if
		// that happens.
		// zap.ErrorOutput(writer.WithMarkFailed(true)),
	}
	zapOptions = append(zapOptions, cfg.zapOptions...)

	return &Logger{
		logger: zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				nil,
				cfg.Level,
			),
			zapOptions...,
		),
	}
}

// Print is the interface for print
func (l *Logger) Print(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Info(v...)
}

// Printf is the interface for printf
func (l *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

// Debug is the interface for debug
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Debug(v...)
}

// Debugf is the interface for debugf
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v...)
}

// Info is the interface for info
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Info(v...)
}

// Infof is the interface for infof
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

// Notice is the interface for notice
func (l *Logger) Notice(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Info(v...)
}

// Noticef is the interface for noticef
func (l *Logger) Noticef(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

// Warning is the interface for warning
func (l *Logger) Warning(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Warn(v...)
}

// Warningf is the interface for warningf
func (l *Logger) Warningf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Warnf(format, v...)
}

// Error is the interface for error
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Error(v...)
}

// Errorf is the interface for errorf
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Errorf(format, v...)
}

// Critical is the interface for critical
func (l *Logger) Critical(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Fatal(v...)
}

// Criticalf is the interface for criticalf
func (l *Logger) Criticalf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Fatalf(format, v...)
}

// Panic is the interface for panic
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Panic(v...)
}

// Panicf is the interface for panicf
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Panicf(format, v...)
}

// Fatal is the interface for fatal
func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.logger.Sugar().Fatal(v...)
}

// Fatalf is the interface for fatalf
func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Sugar().Fatalf(format, v...)
}
