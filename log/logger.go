/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

// Package log is the logger for feie.
package log

import (
	"context"
	"errors"
	"os"

	"github.com/natefinch/lumberjack"
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

var (
	// ErrInvalidKey is the error for invalid key.
	ErrInvalidKey = errors.New("invalid key")
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = Level(zap.DebugLevel)
	// InfoLevel is the default logging priority.
	InfoLevel = Level(zap.InfoLevel)
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = Level(zap.WarnLevel)
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = Level(zap.ErrorLevel)
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = Level(zap.DPanicLevel)
	// PanicLevel logs a message, then panics.
	PanicLevel = Level(zap.PanicLevel)
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = Level(zapcore.FatalLevel)
)

type (
	// Level is the level of logger.
	Level zapcore.Level
	// Logger is the global logger instance.
	Logger struct {
		op    options
		level Level
		log   *zap.Logger
	}
	options struct {
		LogPath string
		Level   Level
	}
	// Option is the option for logger.
	Option func(o *options)
)

// WithLogPath is the option for log path.
func WithLogPath(path string) Option {
	return func(o *options) {
		o.LogPath = path
	}
}

// WithLevel is the option for log level.
func WithLevel(level Level) Option {
	return func(o *options) {
		o.Level = level
	}
}

// New is the global logger instance.
func New(_ context.Context, opts ...Option) *Logger {
	var (
		coreArr []zapcore.Core
		op      = options{
			LogPath: os.TempDir(),
			Level:   InfoLevel,
		}
	)
	for _, option := range opts {
		option(&op)
	}

	// ???????????????
	encoderConfig := zap.NewProductionEncoderConfig()            // NewJSONEncoder()??????json?????????NewConsoleEncoder()????????????????????????
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // ??????????????????
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // ???????????????????????????????????????????????????zapcore.CapitalLevelEncoder????????????
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder       // ????????????????????????
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder // ??????????????????
	encoderConfig.EncodeName = zapcore.FullNameEncoder           // ????????????????????????
	encoder := zapcore.NewConsoleEncoder(encoderConfig)          // ??????????????????????????????????????????????????????????????????????????????

	// ????????????
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error??????
		return lev >= zap.ErrorLevel
	})
	if op.Level <= InfoLevel {
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info???debug??????,debug??????????????????
			return lev < zap.ErrorLevel && lev >= zap.DebugLevel
		})

		// info??????writeSyncer
		infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   op.LogPath + "/log/info.log", // ??????????????????????????????????????????????????????????????????
			MaxSize:    2,                            // ??????????????????,??????MB
			MaxBackups: 100,                          // ??????????????????????????????
			MaxAge:     30,                           // ????????????????????????
			Compress:   false,                        // ??????????????????
		})
		infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority) // ?????????????????????????????????????????????????????????,ErrorLevel???????????????error???????????????
		coreArr = append(coreArr, infoFileCore)
	}
	// error??????writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   op.LogPath + "/log/error.log", // ????????????????????????
		MaxSize:    1,                             // ??????????????????,??????MB
		MaxBackups: 5,                             // ??????????????????????????????
		MaxAge:     30,                            // ????????????????????????
		Compress:   false,                         // ??????????????????
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) // ?????????????????????????????????????????????????????????,ErrorLevel???????????????error???????????????

	coreArr = append(coreArr, errorFileCore)
	return &Logger{
		level: op.Level,
		log:   zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()), // zap.AddCaller()???????????????????????????????????????
	}
}

// Print is the interface for print
func (l *Logger) Print(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Info(v...)
}

// Printf is the interface for printf
func (l *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Infof(format, v...)
}

// Debug is the interface for debug
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Debug(v...)
}

// Debugf is the interface for debugf
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Debugf(format, v...)
}

// Info is the interface for info
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Info(v...)
}

// Infof is the interface for infof
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Infof(format, v...)
}

// Notice is the interface for notice
func (l *Logger) Notice(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Info(v...)
}

// Noticef is the interface for noticef
func (l *Logger) Noticef(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Infof(format, v...)
}

// Warning is the interface for warning
func (l *Logger) Warning(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Warn(v...)
}

// Warningf is the interface for warningf
func (l *Logger) Warningf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Warnf(format, v...)
}

// Error is the interface for error
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Error(v...)
}

// Errorf is the interface for errorf
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Errorf(format, v...)
}

// Critical is the interface for critical
func (l *Logger) Critical(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Fatal(v...)
}

// Criticalf is the interface for criticalf
func (l *Logger) Criticalf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Fatalf(format, v...)
}

// Panic is the interface for panic
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Panic(v...)
}

// Panicf is the interface for panicf
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Panicf(format, v...)
}

// Fatal is the interface for fatal
func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.log.Sugar().Fatal(v...)
}

// Fatalf is the interface for fatalf
func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.log.Sugar().Fatalf(format, v...)
}
