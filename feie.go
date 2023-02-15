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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/houseme/gocrypto"
	"github.com/houseme/gocrypto/rsa"

	"github.com/houseme/feie/log"
)

type options struct {
	User      string
	UKey      string
	Gateway   string
	PublicKey string
	TimeOut   time.Duration
	UserAgent string
	Logger    log.ILogger     // 日志
	DataType  gocrypto.Encode // 数据类型
	HashType  gocrypto.Hash   // Hash类型
	APIName   string          // API名称
	LogPath   string          // 日志路径
	Level     log.Level
}

// FeiE is the feie client.
type FeiE struct {
	op         options
	request    *Request
	secretInfo rsa.SecretInfo
	response   *Response
}

// Option The option is a payment option.
type Option func(o *options)

// WithUser sets the user.
func WithUser(user string) Option {
	return func(o *options) {
		o.User = user
	}
}

// WithUKey sets the ukey.
func WithUKey(ukey string) Option {
	return func(o *options) {
		o.UKey = ukey
	}
}

// WithGateway sets the gateway.
func WithGateway(gateway string) Option {
	return func(o *options) {
		o.Gateway = gateway
	}
}

// WithPublicKey sets the public key.
func WithPublicKey(publicKey string) Option {
	return func(o *options) {
		o.PublicKey = publicKey
	}
}

// WithTimeOut sets the timeout.
func WithTimeOut(timeout time.Duration) Option {
	return func(o *options) {
		o.TimeOut = timeout
	}
}

// WithUserAgent sets the user agent.
func WithUserAgent(userAgent string) Option {
	return func(o *options) {
		o.UserAgent = userAgent
	}
}

// WithLogger sets the logger.
func WithLogger(logger log.ILogger) Option {
	return func(o *options) {
		o.Logger = logger
	}
}

// WithDataType sets the data type.
func WithDataType(dataType gocrypto.Encode) Option {
	return func(o *options) {
		o.DataType = dataType
	}
}

// WithHashType sets the hash type.
func WithHashType(hashType gocrypto.Hash) Option {
	return func(o *options) {
		o.HashType = hashType
	}
}

// WithAPIName sets the api name.
func WithAPIName(apiName string) Option {
	return func(o *options) {
		o.APIName = apiName
	}
}

// WithLogPath sets the log path.
func WithLogPath(logPath string) Option {
	return func(o *options) {
		o.LogPath = logPath
	}
}

// WithLevel sets the level.
func WithLevel(level log.Level) Option {
	return func(o *options) {
		o.Level = level
	}
}

// NewFeiE returns a new feie client.
func NewFeiE(ctx context.Context, opts ...Option) (*FeiE, error) {
	op := options{
		TimeOut:   30 * time.Second,
		Gateway:   Gateway,
		UserAgent: userAgent,
		DataType:  gocrypto.Base64,
		HashType:  gocrypto.SHA256,
		LogPath:   os.TempDir(),
		Level:     log.DebugLevel,
		Logger:    log.New(ctx, log.WithLevel(log.DebugLevel), log.WithLogPath(os.TempDir())),
	}
	for _, option := range opts {
		option(&op)
	}
	return &FeiE{
		op: op,
		secretInfo: rsa.SecretInfo{
			PublicKey:          op.PublicKey,
			PrivateKey:         "",
			PublicKeyDataType:  op.DataType,
			PrivateKeyDataType: op.DataType,
			PrivateKeyType:     gocrypto.PKCS8,
			HashType:           op.HashType,
		},
	}, nil
}

// SetAPIName sets the api name.
func (f *FeiE) SetAPIName(apiName string) {
	f.op.APIName = apiName
}

// SetRequest sets the request.
func (f *FeiE) SetRequest(request *Request) {
	f.request = request
}

// Sha1Sign returns the sha1 sign.
func (f *FeiE) Sha1Sign() string {
	s := sha1.Sum([]byte(f.op.User + f.op.UKey + time.Now().Format("20060102"))) // 20060102150405
	return hex.EncodeToString(s[:])
}

// DoRequest does the request.
func (f *FeiE) DoRequest(ctx context.Context, req *Request) (resp *Response, err error) {
	c, err := client.NewClient()
	if err != nil {
		return
	}

	var postArgs protocol.Args
	postArgs.Set("arg", "a") // Set post args
	status, body, _ := c.Post(ctx, nil, "https://www.example.com", &postArgs)
	fmt.Printf("status=%v body=%v\n", status, string(body))
	// Marshal
	output, err := sonic.Marshal(resp)
	// Unmarshal
	err = sonic.Unmarshal(output, &resp)

	return
}
