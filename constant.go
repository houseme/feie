/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

package feie

var (
	// gateway is the gateway of feieyun.
	// See: http://www.feieyun.com/open/index.html
	gateway = "https://api.feieyun.cn/Api/Open/"

	// userAgent is the user agent of feieyun.
	// See: http://www.feieyun.com/open/index.html
	userAgent = []byte(`Mozilla/5.0 (lanren; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36`)

	// PrinterAddList is the API name for adding a printer.
	// 批量添加打印机，Open_printerAddlist
	// See: http://www.feieyun.com/open/index.html
	PrinterAddList = "Open_printerAddlist"

	// PrinterDelList is the API name for deleting a printer.
	// 批量删除打印机，Open_printerDelList
	// See: http://www.feieyun.com/open/index.html
	PrinterDelList = "Open_printerDelList"

	// PrintMsg is the API name for printing a message.
	// 打印订单，Open_printMsg
	// See: http://www.feieyun.com/open/index.html
	PrintMsg = "Open_printMsg"

	// PrintLabelMsg is the API name for printing a label message.
	// 打印标签，Open_printLabelMsg
	// See: http://www.feieyun.com/open/index.html
	PrintLabelMsg = "Open_printLabelMsg"

	// PrinterEdit is the API name for editing a printer.
	// 编辑打印机，Open_printerEdit
	// See: http://www.feieyun.com/open/index.html
	PrinterEdit = "Open_printerEdit"

	// DelPrinterSqs is the API name for deleting a printer sqs.
	// 删除打印机指令，Open_delPrinterSqs 清空待打印队列
	// See: http://www.feieyun.com/open/index.html
	DelPrinterSqs = "Open_delPrinterSqs"

	// QueryOrderState is the API name for querying the order state.
	// 查询订单状态，Open_queryOrderState 根据订单ID,去查询订单是否打印成功,订单ID由接口Open_printMsg返回
	// See: http://www.feieyun.com/open/index.html
	QueryOrderState = "Open_queryOrderState"

	// QueryOrderInfoByDate is the API name for querying the order info by date.
	// 查询指定打印机某天的订单详情，Open_queryOrderInfoByDate
	// See: http://www.feieyun.com/open/index.html
	QueryOrderInfoByDate = "Open_queryOrderInfoByDate"

	// QueryPrinterStatus is the API name for querying the printer status.
	// 查询打印机状态，Open_queryPrinterStatus
	// 查询指定打印机状态，返回该打印机在线或离线，正常或异常的信息。
	// See: http://www.feieyun.com/open/index.html
	QueryPrinterStatus = "Open_queryPrinterStatus"

	// UserField 飞鹅云后台注册用户名。
	UserField = "user"
	// SysTimeField 当前UNIX时间戳，10位，精确到秒。
	SysTimeField = "stime"
	// SigField 对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。
	SigField = "sig"
	// APINameFiled 请求的接口名称：Open_printMsg
	APINameFiled = "apiname"

	// ExpiredField 订单失效UNIX时间戳，10位，精确到秒，打印时超过该时间该订单将抛弃不打印，取值范围为：当前时间<订单失效时间≤24小时后。
	ExpiredField = "expired"

	// BackURLField 必须先在管理后台设置，回调数据格式详见《订单状态回调》
	BackURLField = "backurl"

	// DebugField debug=1返回非json格式的数据。仅测试时候使用。
	DebugField = "debug"

	// PrinterContentField 打印机编号(必填) # 打印机识别码(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100台。
	PrinterContentField = "printerContent"

	// SNFiled 打印机编号
	SNFiled = "sn"

	// ContentField 打印内容,不能超过5000字节
	ContentField = "content"
)
