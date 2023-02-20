/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

package log

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	type args struct {
		in0  context.Context
		opts []Option
	}
	var (
		opts = []Option{WithLogPath(os.TempDir()), WithLevel(DebugLevel)}
		ctx  = context.Background()
		want = New(ctx, opts...)
	)

	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			name: "TestNew",
			args: args{
				in0:  ctx,
				opts: opts,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.in0, tt.args.opts...); reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
	type fields struct {
		op    options
		level Level
		log   *zap.Logger
	}
	type args struct {
		ctx context.Context
		v   []interface{}
	}
	var (
		opts = []Option{WithLogPath(os.TempDir()), WithLevel(DebugLevel)}
		ctx  = context.Background()
		want = New(ctx, opts...)
	)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestLogger_Debug",
			args: args{
				ctx: ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want.Debug(tt.args.ctx, "xxxx,", "ces time: ", time.Now().Format("20221015181222"))
		})
	}
}
