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
	sysTime    string
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

// New returns a new feie client.
func New(ctx context.Context, opts ...Option) (*FeiE, error) {
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

// SetRequest sets the request.
func (f *FeiE) SetRequest(request *protocol.Request) {
	f.request = request
}

// Response return request response content
func (f *FeiE) Response() *protocol.Response {
	return f.response
}

// sha1Sign returns the sha1 sign.
func (f *FeiE) sha1Sign() string {
	s := sha1.Sum([]byte(f.op.User + f.op.UKey + f.sysTime)) // 20060102150405
	return hex.EncodeToString(s[:])
}

// generateTime Generate current time
func (f *FeiE) generateTime() {
	f.sysTime = strconv.FormatInt(time.Now().Unix(), 10)
}

// doRequest does the request.
func (f *FeiE) doRequest(ctx context.Context, formData map[string]string) error {
	f.generateTime()
	formData[UserField] = f.op.User
	formData[SysTimeField] = f.sysTime
	formData[SigField] = f.sha1Sign()
	f.logger.Debug(ctx, "formData:", formData)
	f.request.SetMultipartFormData(formData)
	f.request.SetRequestURI(gateway)
	f.request.Header.SetMethod(consts.MethodPost)
	f.request.Header.SetUserAgentBytes(userAgent)
	f.logger.Debug(ctx, "request content: ", f.request)

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	f.logger.Debug(ctx, "do request start")
	err = c.Do(ctx, f.request, f.response)
	if err != nil {
		return err
	}
	f.logger.Debug(ctx, "do request end")
	return nil
}

// OpenPrintMsg 打印订单
// 发送用户需要打印的订单内容给飞鹅云小票打印机 （该接口只能是小票机使用，如购买标签机请使用标签机专用接口）
// see: http://help.feieyun.com/document.php
func (f *FeiE) OpenPrintMsg(ctx context.Context, sn, content string) (resp *PrintMsgResp, err error) {
	var formData = make(map[string]string, 6)
	formData[APINameFiled] = PrintMsg
	formData[SNFiled] = sn
	formData[ContentField] = content

	if err = f.doRequest(ctx, formData); err != nil {
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

// OpenPrinterAddList 批量添加打印机
// 批量添加打印机，请严格参照格式说明：
// 批量添加规则：
//
// 打印机编号SN(必填) # 打印机识别码KEY(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100行(台)。
// 每次最多添加100台。
func (f *FeiE) OpenPrinterAddList(ctx context.Context, content string) (resp *PrinterAddResp, err error) {
	var formData = make(map[string]string, 6)
	formData[APINameFiled] = PrinterAddList
	formData[PrinterContentField] = content

	if err = f.doRequest(ctx, formData); err != nil {
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

// OpenPrinterDelList 删除批量打印机
// content 打印机编号，多台打印机请用减号“-”连接起来。
// see: http://help.feieyun.com/document.php
func (f *FeiE) OpenPrinterDelList(ctx context.Context, content string) (resp *PrinterDelResp, err error) {
	var formData = make(map[string]string, 6)
	formData[SNListField] = content

	if err = f.doRequest(ctx, formData); err != nil {
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
