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
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
	UserAgent []byte
	DataType  gocrypto.Encode // 数据类型
	HashType  gocrypto.Hash   // Hash类型
	APIName   string          // API名称
	LogPath   string          // 日志路径
	Level     log.Level
}

// FeiE is the feie client.
type FeiE struct {
	request    *protocol.Request
	response   *protocol.Response
	logger     log.ILogger
	op         options
	secretInfo rsa.SecretInfo
	sysTime    time.Time
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
func WithUserAgent(userAgent []byte) Option {
	return func(o *options) {
		o.UserAgent = userAgent
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
		Gateway:   gateway,
		UserAgent: userAgent,
		DataType:  gocrypto.Base64,
		HashType:  gocrypto.SHA256,
		LogPath:   os.TempDir(),
		Level:     log.DebugLevel,
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
		logger:   log.New(ctx, log.WithLevel(op.Level), log.WithLogPath(op.LogPath)),
		request:  &protocol.Request{},
		response: &protocol.Response{},
	}, nil
}

// SetAPIName sets the api name.
func (f *FeiE) SetAPIName(apiName string) {
	f.op.APIName = apiName
}

// SetRequest sets the request.
func (f *FeiE) SetRequest(request *protocol.Request) {
	f.request = request
}

// Sha1Sign returns the sha1 sign.
func (f *FeiE) Sha1Sign() string {
	s := sha1.Sum([]byte(f.op.User + f.op.UKey + f.sysTime.Format("20060102"))) // 20060102150405
	return hex.EncodeToString(s[:])
}

// GenerateTime Generate current time
func (f *FeiE) GenerateTime() {
	f.sysTime = time.Now()
}

// DoRequest does the request.
func (f *FeiE) DoRequest(ctx context.Context) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	f.request.SetRequestURI(gateway)
	f.logger.Debug(ctx, "do request start")
	err = c.Do(ctx, f.request, f.response)
	if err != nil {
		return err
	}
	f.logger.Debug(ctx, "do request end")
	return nil
}

// OpenPrintMsg 打印订单
// see: http://help.feieyun.com/document.php
func (f *FeiE) OpenPrintMsg(ctx context.Context, sn, content string) (resp *PrintMsgResp, err error) {
	var formData = make(map[string]string, 6)
	formData[UserField] = f.op.User
	formData[SysTimeField] = strconv.FormatInt(f.sysTime.Unix(), 10)
	formData[APINameFiled] = PrintMsg
	formData[SNFiled] = sn
	formData[ContentField] = content
	formData[SigField] = f.Sha1Sign()
	f.logger.Debug(ctx, "formData:", formData)
	f.request.SetMultipartFormData(formData)
	f.request.Header.SetMethod(consts.MethodPost)
	f.request.Header.SetUserAgentBytes(userAgent)
	f.logger.Debug(ctx, f.request)
	if err = f.DoRequest(ctx); err != nil {
		return
	}
	f.logger.Debug(ctx, "do request response body:", string(f.response.Body()))
	if !f.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(f.response.Body(), &resp); err != nil {
		return
	}
	f.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}
