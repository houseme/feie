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
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (f *FeiE) initLog(ctx context.Context, op options) {
	dynamicLevel := zap.NewAtomicLevel()
	dynamicLevel.SetLevel(zap.DebugLevel)
	logger := hertzzap.NewLogger(
		hertzzap.WithCores([]hertzzap.CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(humanEncoderConfig()),
				Ws:  os.Stdout,
				Lvl: dynamicLevel,
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(op.LogPath + "/all.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(op.LogPath + "/debug.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.DebugLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(op.LogPath + "/info.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.InfoLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(op.LogPath + "/warn.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.WarnLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(op.LogPath + "/error.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev >= zap.ErrorLevel
					}))),
			},
		}...),
	)
	defer logger.Sync()
	hlog.SetLogger(logger)
	hlog.SetLevel(op.Level)
	f.logger = logger
	f.logger.SetLevel(op.Level)
	f.logger.CtxInfof(ctx, "feie dada init log start level:%s", op.Level)
}

// humanEncoderConfig copy from zap
func humanEncoderConfig() zapcore.EncoderConfig {
	cfg := encoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	return cfg
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10,
		MaxBackups: 50000,
		MaxAge:     1000,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// encoderConfig encoder config for testing, copy from zap
func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
