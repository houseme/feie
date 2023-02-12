/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

package feie

var (
	// Gateway is the gateway of feieyun.
	// See: http://www.feieyun.com/open/index.html
	Gateway = "https://api.feieyun.cn/Api/Open/"

	// userAgent is the user agent of feieyun.
	// See: http://www.feieyun.com/open/index.html
	userAgent = `Mozilla/5.0 (lanren; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36`

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
)
