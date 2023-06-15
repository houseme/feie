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

var (
	// gateway is the gateway of feieyun.
	// See: http://www.feieyun.com/open/index.html
	gateway = "https://api.feieyun.cn/Api/Open/"

	// userAgent is the user agent of feieyun.
	// See: http://www.feieyun.com/open/index.html
	userAgent = []byte(`Mozilla/5.0 (FeiE; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36`)

	// printerAddList is the API name for adding a printer.
	// 批量添加打印机，Open_printerAddlist
	// See: http://www.feieyun.com/open/index.html
	printerAddList = "Open_printerAddlist"

	// printerDelList is the API name for deleting a printer.
	// 批量删除打印机，Open_printerDelList
	// See: http://www.feieyun.com/open/index.html
	printerDelList = "Open_printerDelList"

	// printMsg is the API name for printing a message.
	// 打印订单，Open_printMsg
	// See: http://www.feieyun.com/open/index.html
	printMsg = "Open_printMsg"

	// printLabelMsg is the API name for printing a label message.
	// 打印标签，Open_printLabelMsg
	// See: http://www.feieyun.com/open/index.html
	printLabelMsg = "Open_printLabelMsg"

	// printerEdit is the API name for editing a printer.
	// 编辑打印机，Open_printerEdit
	// See: http://www.feieyun.com/open/index.html
	printerEdit = "Open_printerEdit"

	// delPrinterSqs is the API name for deleting a printer sqs.
	// 删除打印机指令，Open_delPrinterSqs 清空待打印队列
	// See: http://www.feieyun.com/open/index.html
	delPrinterSqs = "Open_delPrinterSqs"

	// queryOrderState is the API name for querying the order state.
	// 查询订单状态，OpenQueryOrderState,根据订单ID去查询订单是否打印成功,订单ID由接口Open_printMsg返回
	// See: http://www.feieyun.com/open/index.html
	queryOrderState = "Open_queryOrderState"

	// queryOrderInfoByDate is the API name for querying the order info by date.
	// 查询指定打印机某天的订单详情，Open_queryOrderInfoByDate
	// See: http://www.feieyun.com/open/index.html
	queryOrderInfoByDate = "Open_queryOrderInfoByDate"

	// queryPrinterStatus is the API name for querying the printer status.
	// 查询打印机状态，Open_queryPrinterStatus
	// 查询指定打印机状态，返回该打印机在线或离线，正常或异常的信息。
	// See: http://www.feieyun.com/open/index.html
	queryPrinterStatus = "Open_queryPrinterStatus"

	// UserField 飞鹅云后台注册用户名。
	UserField = "user"
	// SysTimeField 当前UNIX时间戳，10位，精确到秒。
	SysTimeField = "stime"
	// SigField 对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。
	SigField = "sig"
	// APINameField 请求的接口名称：Open_printMsg
	APINameField = "apiname"

	// ExpiredField 订单失效UNIX时间戳，10位，精确到秒，打印时超过该时间该订单将抛弃不打印，取值范围为：当前时间<订单失效时间≤24小时后。
	ExpiredField = "expired"

	// BackURLField 必须先在管理后台设置，回调数据格式详见《订单状态回调》
	BackURLField = "backurl"

	// DebugField debug=1返回非json格式的数据。仅测试时候使用。
	DebugField = "debug"

	// SNListField 打印机编号，多台打印机请用减号“-”连接起来。
	SNListField = "snlist"

	// PrinterContentField 打印机编号(必填) # 打印机识别码(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100台。
	PrinterContentField = "printerContent"

	// SNField 打印机编号
	SNField = "sn"

	// ContentField 打印内容,不能超过5000字节
	ContentField = "content"

	// TimesField 打印次数，默认为1。
	TimesField = "times"

	// ImgField 图片二进制数据，需配合<IMG>标签使用，最佳效果为不大于224px的正方形(宽高都为8的倍数)黑白图，支持jpg、png、bmp，不能超过10K
	ImgField = "img"

	// NameField 打印机备注名称
	NameField = "name"

	// PhoneNumField 打印机流量卡号码
	PhoneNumField = "phoneNum"

	// OrderIDField 订单ID，由接口Open_printMsg返回。
	OrderIDField = "orderid"

	// DateField 查询日期，格式为：2016-08-08
	DateField = "date"
)
