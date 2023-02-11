## 柴树杉

凹语言本身目前缺少完备的异常和错误等规范，因此海不具备构建单元测试工具的条件。但随着项目标准库的启动和编译器本身的复杂化，同时为了向可用性的目标逼近对于代码的稳健性和测试都更加迫切。因此迫切需要一个适合当前需求的测试方案。

为了方便被集成，凹语言提供了api包，可以作为Go语言环境内的脚本被执行：

```go
package main

import (
	"fmt"
	"wa-lang.org/wa/api"
)

func main() {
	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", `func main { println(42) }`)
	fmt.Print(string(output), err)
}
```

以上可以通过 `api.RunCode` 来执行一个凹程序，同时获得返回结果或错误信息。因此可以尝试基于该API来设计临时的测试工具。

比如有 `hash/crc32` 包，对应的目录文件如下：

```
/hash/crc32/
  |-- crc32.wa // 实现代码
  |-- test/test_01.wa         // 测试代码1
  |-- test/test_01.stdout.txt // 测试代码1 正常输出文件
  |-- test/test_01.stderr.txt // 测试代码1 错误输出文件
```

因此可以构建一个自动扫描目录大工具，当发现 `test_xxx.wa` 模式的文件时，自动通过 `api.RunCode` 加载该文件并比对输出结果和错误信息。每个 `test_xxx.wa` 文件类似集成测试是独立执行的，不像单元测试那种效率更好。不过应该可以满足当前阶段测试需求。

基于以上的思路，每个标准库的测试代码可以在 test 子目录单独维护管理，也可以用于普通模块的测试管理。

抛砖引玉下，欢迎提供更多建议！
