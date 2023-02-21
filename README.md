# FeiE Printer SDK for Golang
Flying Goose Cloud print go language sdk

## Installation

```bash
go get -u -v github.com/houseme/feie@main 
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/houseme/feie"
)

func main() {
    ctx := context.Background()
    f, err := feie.New(ctx)
    if err != nil {
        panic(err)
    }
    // 添加打印机
    printerAddResp, err := f.OpenPrinterAddList(ctx, &feie.PrinterAddReq{
        PrinterContent: "xxxxxx",
    })
    
    if err != nil {
    
    }
    fmt.Println("PrinterAddResp:", printerAddResp)
    
    // 执行打印
    printMsgReq, err := f.OpenPrintMsg(ctx, &feie.PrintMsgReq{
        SN: "xxxxx",
    })
    if err != nil {
    
    }
    fmt.Println("PrintMsgResp:", printMsgReq)
}

```


## License
FeiE is primarily distributed under the terms of both the [Apache License (Version 2.0)](LICENSE)