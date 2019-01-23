package main

import (
    "github.com/zywaited/php-go/phpgo"
)

func module_startup(ptype int, module_number int) int {
    return 1
}

func module_shutdown(ptype int, module_number int) int {
    return 1
}

func request_startup(ptype int, module_number int) int {
    return 1
}

func request_shutdown(ptype int, module_number int) int {
    return 1
}

func init() {
    phpgo.InitExtension("rsp", "1.0")
    phpgo.RegisterInitFunctions(module_startup, module_shutdown, request_startup, request_shutdown)
    phpgo.AddClass("Waited\\RSP", NewRsp)
}

func main() {
    // 引入cellnet
}
