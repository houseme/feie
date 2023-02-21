/*
 *   Copyright `feie` Author. All Rights Reserved.
 *
 *   This Source Code Form is subject to the terms of the MIT License.
 *   If a copy of the MIT was not distributed with this file,
 *   You can obtain one at https://github.com/housme/feie.
 */

package feie

// PrinterAddReq is the request body for adding a printer.
type PrinterAddReq struct {
	User           string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime          int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig            string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName        string `json:"apiname" description:"固定值Open_printerAddlist。"`
	Debug          int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	PrinterContent string `json:"printerContent" description:"打印机编号(必填) # 打印机识别码(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100台。"`
}

// PrinterAddResp is the response body for adding a printer.
type PrinterAddResp struct {
	Ret                int              `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string           `json:"msg" description:"错误信息。"`
	Data               *PrinterRespData `json:"data,omitempty" description:"成功时返回的数据。"`
	ServerExecutedTime int64            `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// PrinterRespData is the data field of PrinterResp.
type PrinterRespData struct {
	Ok []*string `json:"ok"`
	No []*string `json:"no"`
}

// PrinterDelReq is the request body for deleting a printer.
type PrinterDelReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_printerDelList。"`
	Debug   int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	SNList  string `json:"snlist" description:"打印机编号，多台打印机请用减号“-”连接起来。"`
}

// PrinterDelResp is the response body for deleting a printer.
type PrinterDelResp struct {
	Ret                int              `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string           `json:"msg" description:"错误信息。"`
	Data               *PrinterRespData `json:"data,omitempty" description:"成功时返回的数据。"`
	ServerExecutedTime int64            `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// PrintMsgReq is the request body for printing a message.
type PrintMsgReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_printMsg。"`
	Expired int64  `json:"expired,omitempty" description:"打印任务过期时间，单位秒。"`
	BackURL string `json:"backurl,omitempty" description:"打印任务回调地址。"`
	SN      string `json:"sn" description:"打印机编号。"`
	Content string `json:"content" description:"打印内容。"`
	Times   int    `json:"times,omitempty" description:"打印联数，最大支持10联。"`
}

// PrintMsgResp is the response body for printing a message.
type PrintMsgResp struct {
	Msg                string `json:"msg" description:"错误信息。"`
	Data               string `json:"data" description:"成功时返回的数据，正确返回订单ID。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
}

// PrintLabelMsgReq is the request body for printing a label message.
type PrintLabelMsgReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_printLabelMsg。"`
	Expired int64  `json:"expired,omitempty" description:"打印任务过期时间，单位秒。"`
	BackURL string `json:"backurl,omitempty" description:"打印任务回调地址。"`
	SN      string `json:"sn" description:"打印机编号。"`
	Content string `json:"content" description:"打印内容。"`
	Times   int    `json:"times,omitempty" description:"打印联数，最大支持10联。"`
	Img     string `json:"img" description:"图片base64编码。"`
}

// PrintLabelMsgResp is the response body for printing a label message.
type PrintLabelMsgResp struct {
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string `json:"msg" description:"错误信息。"`
	Data               string `json:"data" description:"成功时返回的数据，正确返回订单ID。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// PrinterEditReq is the request body for editing a printer.
type PrinterEditReq struct {
	User     string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime    int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig      string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName  string `json:"apiname" description:"固定值Open_printerEdit。"`
	Debug    int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	SN       string `json:"sn" description:"打印机编号。"`
	Name     string `json:"name" description:"打印机备注名称。"`
	PhoneNum string `json:"phonenum,omitempty" description:"打印机流量卡号码。"`
}

// PrinterEditResp is the response body for editing a printer.
type PrinterEditResp struct {
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string `json:"msg" description:"错误信息。"`
	Data               bool   `json:"data" description:"成功时返回的数据。成功返回true，失败返回false。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// DelPrinterSQSReq is the request body for deleting a printer from SQS.
type DelPrinterSQSReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_delPrinterSQS。"`
	Debug   int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	SN      string `json:"sn" description:"打印机编号。"`
}

// DelPrinterSQSResp is the response body for deleting a printer from SQS.
type DelPrinterSQSResp struct {
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string `json:"msg" description:"错误信息。"`
	Data               bool   `json:"data" description:"成功时返回的数据,正确返回true。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// QueryOrderStateReq is the request body for querying the state of an order.
type QueryOrderStateReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_queryOrderState。"`
	Debug   int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	OrderID string `json:"orderid" description:"订单ID，由接口Open_printMsg返回。"`
}

// QueryOrderStateResp is the response body for querying the state of an order.
type QueryOrderStateResp struct {
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string `json:"msg" description:"错误信息。"`
	Data               bool   `json:"data" description:"成功时返回的数据,已打印返回true,未打印返回false。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// QueryOrderInfoByDateReq is the request body for querying the information of an order by date.
type QueryOrderInfoByDateReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_queryOrderInfoByDate。"`
	Debug   int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	SN      string `json:"sn" description:"打印机编号。"`
	Date    string `json:"date" description:"查询日期，格式YY-MM-DD，如：2016-09-20"`
}

// QueryOrderInfoByDateResp is the response body for querying the information of an order by date.
type QueryOrderInfoByDateResp struct {
	Ret                int                       `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string                    `json:"msg" description:"错误信息。"`
	Data               *QueryOrderInfoByDateData `json:"data" description:"成功时返回的数据,订单信息。"`
	ServerExecutedTime int64                     `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// QueryOrderInfoByDateData is the data of QueryOrderInfoByDateResp.
type QueryOrderInfoByDateData struct {
	Print   int `json:"print" description:"打印份数。"`
	Waiting int `json:"waiting" description:"等待打印份数。"`
}

// QueryPrinterStatusReq is the request body for querying the status of a printer.
type QueryPrinterStatusReq struct {
	User    string `json:"user" description:"飞鹅云后台注册用户名。"`
	STime   int64  `json:"stime,string" description:"当前UNIX时间戳，10位，精确到秒。"`
	Sig     string `json:"sig" description:"签名，详见签名算法。对参数 user+UKEY+stime拼接后（+号表示连接符）进行SHA1加密得到签名，加密后签名值为40位小写字符串。"`
	APIName string `json:"apiname" description:"固定值Open_queryPrinterStatus。"`
	Debug   int    `json:"debug,string,omitempty" description:"debug=1返回非json格式的数据。仅测试时候使用。"`
	SN      string `json:"sn" description:"打印机编号。"`
}

// QueryPrinterStatusResp is the response body for querying the status of a printer.
type QueryPrinterStatusResp struct {
	Ret                int    `json:"ret" description:"错误码，0为成功，非0为错误。"`
	Msg                string `json:"msg" description:"错误信息。"`
	Data               string `json:"data" description:"成功时返回的数据,打印机状态。返回打印机状态信息共三种：1、离线。2、在线，工作状态正常。3、在线，工作状态不正常。备注：异常一般是无纸，离线的判断是打印机与服务器失去联系超过2分钟。"`
	ServerExecutedTime int64  `json:"serverExecutedTime" description:"服务器执行时间，单位毫秒。"`
}

// AsyncPrinterResultReq is the request body for querying the result of an asynchronous request.
type AsyncPrinterResultReq struct {
	OrderID string `json:"orderId" description:"订单ID，由接口Open_printMsg返回。"`
	Sign    string `json:"sign" description:"签名，详见签名算法。SHA256WithRSA验证签名值"`
	Status  int    `json:"status" description:"订单状态"`
	Stime   int    `json:"stime" description:"订单状态变更UNIX时间戳，10位，精确到秒。"`
}

// AsyncPrinterResultResp is the response body for querying the result of an asynchronous request.
type AsyncPrinterResultResp struct {
	VerifySign bool   `json:"verifySign" description:"验签结果"`
	OrderID    string `json:"orderId" description:"订单ID，由接口Open_printMsg返回。"`
	Status     int    `json:"status" description:"订单状态"`
	Stime      int    `json:"stime" description:"订单状态变更UNIX时间戳，10位，精确到秒。"`
}
