/*
 *  Copyright `FeiE` Author(https://houseme.github.io/feie/). All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  You can obtain one at https://github.com/houseme/feie.
 */

package feie

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/houseme/gocrypto"
	"github.com/houseme/gocrypto/rsa"

	"github.com/houseme/feie/internal"
)

// New returns a new feie client.
func New(ctx context.Context, opts ...Option) *Client {
	op := options{
		TimeOut:   30 * time.Second,
		Gateway:   gateway,
		UserAgent: userAgent,
		DataType:  gocrypto.Base64,
		HashType:  gocrypto.SHA256,
		LogPath:   os.TempDir(),
		Level:     Level(hlog.LevelDebug),
	}
	for _, option := range opts {
		option(&op)
	}
	c := &Client{
		op: op,
		secretInfo: rsa.SecretInfo{
			PublicKey:          op.PublicKey,
			PrivateKey:         "",
			PublicKeyDataType:  op.DataType,
			PrivateKeyDataType: op.DataType,
			PrivateKeyType:     gocrypto.PKCS8,
			HashType:           op.HashType,
		},
		logger:   internal.InitLog(ctx, op.LogPath, hlog.Level(op.Level)),
		request:  &protocol.Request{},
		response: &protocol.Response{},
		user:     op.User,
		ukey:     op.UKey,
	}
	c.logger.SetLevel(hlog.Level(op.Level))
	c.logger.CtxDebugf(ctx, "feie client init success %s", c.op.Level)
	return c
}

// SetRequest sets the request.
func (c *Client) SetRequest(request *protocol.Request) {
	c.request = request
}

// Response return request response content
func (c *Client) Response() *protocol.Response {
	return c.response
}

// SetLogger set feie logger
func (c *Client) SetLogger(logger hlog.FullLogger) {
	c.logger = logger
}

// SetUserKey sets the user key.
func (c *Client) SetUserKey(ukey string) {
	c.ukey = ukey
}

// Reset reset the feie client.
func (c *Client) Reset() {
	if strings.TrimSpace(c.op.User) != "" {
		c.user = c.op.User
	}
	if strings.TrimSpace(c.op.UKey) != "" {
		c.ukey = c.op.UKey
	}
}

// sha1Sign returns the sha1 sign.
func (c *Client) sha1Sign() string {
	s := sha1.Sum([]byte(c.user + c.ukey + c.sysTime))
	return hex.EncodeToString(s[:])
}

// generateTime Generate current time
func (c *Client) generateTime() {
	c.sysTime = strconv.FormatInt(time.Now().Unix(), 10)
}

// doRequest does the request.
func (c *Client) doRequest(ctx context.Context, formData map[string]string) error {
	c.generateTime()

	formData[UserField] = c.user
	formData[SysTimeField] = c.sysTime
	formData[SigField] = c.sha1Sign()
	c.logger.Debug(ctx, "formData:", formData)
	c.request.SetMultipartFormData(formData)
	c.request.SetRequestURI(c.op.Gateway)
	c.request.Header.SetMethod(consts.MethodPost)
	c.request.Header.SetUserAgentBytes(c.op.UserAgent)
	c.logger.Debug(ctx, "request content: ", c.request)

	hc, err := client.NewClient(client.WithTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}), client.WithDialTimeout(c.op.TimeOut))
	if err != nil {
		return err
	}

	c.logger.Debug(ctx, "do request start")
	err = hc.Do(ctx, c.request, c.response)
	if err != nil {
		return err
	}
	c.logger.Debug(ctx, "do request end")
	return nil
}

// OpenPrintMsg 打印订单
// 发送用户需要打印的订单内容给飞鹅云小票打印机 （该接口只能是小票机使用，如购买标签机请使用标签机专用接口）
// see: http://help.feieyun.com/document.php
// ----------接口返回值说明----------
// 正确例子：{"msg":"ok","ret":0,"data":"xxxx_xxxx_xxxxxxxxx","serverExecutedTime":6}
// 错误：{"msg":"错误信息.","ret":非零错误码,"data":null,"serverExecutedTime":5}
func (c *Client) OpenPrintMsg(ctx context.Context, req *PrintMsgReq) (resp *PrintMsgResp, err error) {
	var formData = make(map[string]string)
	formData[APINameField] = printMsg
	formData[SNField] = req.SN
	formData[ContentField] = req.Content
	if req.User != "" {
		c.user = req.User
	}
	if req.Expired > time.Now().Unix() {
		formData[ExpiredField] = strconv.FormatInt(req.Expired, 10)
	}
	if req.Times > 1 {
		formData[TimesField] = strconv.Itoa(req.Times)
	}
	if strings.TrimSpace(req.BackURL) != "" {
		formData[BackURLField] = req.BackURL
	}

	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenPrinterAddList 批量添加打印机
// 批量添加打印机，请严格参照格式说明：
// 批量添加规则：
//
// 打印机编号SN(必填) # 打印机识别码KEY(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100行(台)。
// 每次最多添加100台。
// 提示：打印机编号(必填) # 打印机识别码(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100行(台)。
// snlist := "sn1#key1#remark1#carnum1\nsn2#key2#remark2#carnum2"
//
// ----------接口返回值说明----------
// 正确例子：{"msg":"ok","ret":0,"data":{"ok":["sn#key#remark#carnum","316500011#abcdefgh#快餐前台"],"no":["316500012#abcdefgh#快餐前台#13688889999  （错误：识别码不正确）"]},"serverExecutedTime":3}
// 错误：{"msg":"参数错误 : 该帐号未注册.","ret":-2,"data":null,"serverExecutedTime":37}
func (c *Client) OpenPrinterAddList(ctx context.Context, req *PrinterAddReq) (resp *PrinterAddResp, err error) {
	var formData = make(map[string]string, 5)
	formData[APINameField] = printerAddList
	formData[PrinterContentField] = req.PrinterContent
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenPrinterDelList 删除批量打印机
// content 打印机编号，多台打印机请用减号“-”连接起来。
// see: http://help.feieyun.com/document.php
func (c *Client) OpenPrinterDelList(ctx context.Context, req *PrinterDelReq) (resp *PrinterDelResp, err error) {
	var formData = make(map[string]string, 5)
	formData[SNListField] = req.SNList
	formData[APINameField] = printerDelList
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenPrintLabelMsg 标签机打印订单
// 发送用户需要打印的订单内容给飞鹅云标签打印机（该接口只能是标签机使用，其它型号打印机请勿使用该接口）
// see: http://help.feieyun.com/document.php
func (c *Client) OpenPrintLabelMsg(ctx context.Context, req *PrintLabelMsgReq) (resp *PrintLabelMsgResp, err error) {
	var formData = make(map[string]string)
	formData[APINameField] = printLabelMsg
	formData[SNField] = req.SN
	formData[ContentField] = req.Content
	if req.User != "" {
		c.user = req.User
	}
	if req.Expired > time.Now().Unix() {
		formData[ExpiredField] = strconv.FormatInt(req.Expired, 10)
	}
	if req.Times > 1 {
		formData[TimesField] = strconv.Itoa(req.Times)
	}
	if strings.TrimSpace(req.BackURL) != "" {
		formData[BackURLField] = req.BackURL
	}

	if len(req.Img) > 0 {
		formData[ImgField] = req.Img
	}

	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenPrinterEdit 修改打印机信息
// 修改打印机信息
// see: http://help.feieyun.com/document.php
func (c *Client) OpenPrinterEdit(ctx context.Context, req *PrinterEditReq) (resp *PrinterEditResp, err error) {
	var formData = make(map[string]string)
	formData[SNField] = req.SN
	formData[APINameField] = printerEdit
	formData[NameField] = req.Name
	if req.User != "" {
		c.user = req.User
	}
	if len(strings.TrimSpace(req.PhoneNum)) > 0 {
		formData[PhoneNumField] = strings.TrimSpace(req.PhoneNum)
	}

	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenDelPrinterSQS 清空待打印队列
// see: http://help.feieyun.com/document.php
func (c *Client) OpenDelPrinterSQS(ctx context.Context, req *DelPrinterSQSReq) (resp *DelPrinterSQSResp, err error) {
	var formData = make(map[string]string, 5)
	formData[SNField] = req.SN
	formData[APINameField] = delPrinterSqs
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenQueryOrderState 查询订单是否打印成功
// 根据订单ID,去查询订单是否打印成功,订单ID由接口Open_printMsg返回
// see: http://help.feieyun.com/document.php
func (c *Client) OpenQueryOrderState(ctx context.Context, req *QueryOrderStateReq) (resp *QueryOrderStateResp, err error) {
	var formData = make(map[string]string, 5)
	formData[OrderIDField] = req.OrderID
	formData[APINameField] = queryOrderState
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenQueryOrderInfoByDate 查询指定打印机某天的订单统计数
// 根据打印机编号和日期，查询该打印机某天的订单统计数,查询指定打印机某天的订单详情，返回已打印订单数和等待打印数。
// see: http://help.feieyun.com/document.php
func (c *Client) OpenQueryOrderInfoByDate(ctx context.Context, req *QueryOrderInfoByDateReq) (resp *QueryOrderInfoByDateResp, err error) {
	var formData = make(map[string]string, 6)
	formData[SNField] = req.SN
	formData[DateField] = req.Date
	formData[APINameField] = queryOrderInfoByDate
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// OpenQueryPrinterStatus 查询打印机状态
// 根据打印机编号，查询打印机状态，返回打印机状态。
// 查询指定打印机状态，返回该打印机在线或离线，正常或异常的信息。
// see: http://help.feieyun.com/document.php
func (c *Client) OpenQueryPrinterStatus(ctx context.Context, req *QueryPrinterStatusReq) (resp *QueryPrinterStatusResp, err error) {
	var formData = make(map[string]string, 5)
	formData[SNField] = req.SN
	formData[APINameField] = queryPrinterStatus
	if req.User != "" {
		c.user = req.User
	}
	if err = c.doRequest(ctx, formData); err != nil {
		return
	}
	c.logger.Debug(ctx, "do request response body:", string(c.response.Body()))
	if !c.response.HasBodyBytes() {
		err = errors.New("response is empty")
		return
	}
	if err = sonic.Unmarshal(c.response.Body(), &resp); err != nil {
		return
	}
	c.logger.Debug(ctx, "json Unmarshal resp result:", resp)
	return
}

// AsyncPrinterResult 异步打印结果
// 接口提供方式
// 回调接口统一使用https post方式，contentType 为“application/x-www-form-urlencoded”。
// 需要服务商在开发者后台配置回调地址白名单，并在提交打印订单时提交回调地址。
// 回调地址为服务商提供的接口地址，服务商需要在接口中处理打印结果。
// 回调参数:
//
//	| - orderId 订单ID，由接口Open_printMsg返回。
//	| - status 打印状态，1：打印成功，1：打印失败。
//	| - stime 订单状态变更UNIX时间戳，10位，精确到秒。
//	| - sign 数字签名
//
// 数字签名验证说明
// 2.1.1 获取待验证签名字符串
// 获取所有飞鹅云开放平台的 post 内容，不包括字节类型参数，如文件、字节流，剔除 sign 字段，剔除值为空的参数；按照第一个字符的键值 ASCII 码递增排序（字母升序排序），如果遇到相同字符则按照第二个字符的键值 ASCII 码递增排序，以此类推；将排序后的参数与其对应值，组合成 参数=参数值 的格式，并且把这些参数用 & 字符连接起来，此时生成的字符串为待验签字段为，如：
// orderId：816501678_20160919184316_1419533539
// status：1 stime：1625194910 则待验签字段为：orderId=816501678_20160919184316_1419533539&status=1&stime=1625194910
// 2.1.2 取出签名值sign
// NW1BNm4oTxyyPBdXHPwuI5gjh2onvyHavrSLnrPAGCp4TnoX1IJTwwX+tXFybdi+bo+OM/1FoIeU4H70fPw0m/z/Fz6EYdDpsBbUZFbbUdj9OJrzY/sdnArkynnYoVkLGOwV0DM1WvCn3iqlskD5O1K6POFDc0006xMK+d3/SSNegSUPMuIvuXG6VKGiDN0rO9hOdXFjrp0b1Td14ofPXKibmGXV7XikC2suU45nWmCBC8lKzhazCiInS/tkRAF8WsS2AiACeMvmonyrT/LZWbsfrd9k6M+kATCOz7EjPEd9z+W8N8Rtbur1m3MZdjAshMfduqQEpRU+w7U6R4sxQA==
// 2.1.3 将签名参数（sign）使用 base64 解码为签名值字节码串
// 2.1.4 使用飞鹅云公钥、待验证签名字符串、签名值字节码串进行SHA256WithRSA验证签名值是否正确。
// 3 返回示例：
// 注：开发者接收后需立即返回SUCCESS，如5秒内不返回或返回数据格式错误，平台会重新推送。
// `SUCCESS`
// see: http://help.feieyun.com/document.php
func (c *Client) AsyncPrinterResult(ctx context.Context, req *AsyncPrinterResultReq) (resp *AsyncPrinterResultResp, err error) {
	var (
		handle     = rsa.NewRSACrypt(c.secretInfo)
		content    = "orderId=" + req.OrderID + "&status=" + strconv.Itoa(req.Status) + "&stime=" + strconv.Itoa(req.Stime)
		verifySign bool
	)

	if verifySign, err = handle.VerifySign(content, req.Sign, c.op.DataType); err != nil {
		c.logger.Debug(ctx, "AsyncPrinterResult rsa VerifySign failed:", err)
		return
	}
	resp = &AsyncPrinterResultResp{
		VerifySign: verifySign,
		OrderID:    req.OrderID,
		Stime:      req.Stime,
		Status:     req.Status,
	}
	return
}
