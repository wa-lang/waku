
## 丁尔男

在某些底层 `.wa` 代码，比如堆管理模块中，会出现两个特别的编译标记：`#wa:runtime_getter` 和 `wa:runtime_setter`，它们一般用于需要直接通过地址读写内存的情况。用法如下：

- `#wa:runtime_getter` 使用这个标记的函数，签名形式必须为: `func XXX(ptr u32) => T`  
其中，`T` 是一个基本数据类型或由基本数据类型组合而成的结构体类型。使用该标记修饰的函数，只声明不实现，编译器会自动生成其函数体，用于从 `ptr` 地址的内存中读取一个类型为 `T` 的值。

- `#wa:runtime_setter` 使用这个标记的函数，签名形式必须为: `func YYY(ptr u32, data T)`  
其中，`T` 是一个基本数据类型或由基本数据类型组合而成的结构体类型。使用该标记修饰的函数，只声明不实现，编译器会自动生成其函数体，用于将 `data` 值写入内存的 `ptr` 地址处。

一个实际的用例如下：

```wa
//版权 @2022 凹语言 作者。保留所有权利。

type Block struct{
	size: u32
	used: u32
    next: u32
}

#wa:runtime_getter
func GetBlock(ptr: u32) => Block

#wa:runtime_setter
func SetBlock(ptr: u32, data: Block)

// 堆基址
var heap_base u32

func main(){
	//从heap_base地址中读取Block：
	head := GetBlock(heap_base)

	//从head.next地址中读取next：
	next := GetBlock(head.next)
	next.used = 1

	//next写入head.next地址：
	SetBlock(head.next, next)
}
```

上述代码中，函数 `GetBlock` 和 `SetBlock` 的声明分别使用了 `#wa:runtime_getter` 和 `wa:runtime_setter` 标记，编译器会自动生成它们的函数体，比如 wasm 后端会生成如下：

```wat
(func $__main__.GetBlock (export "__main__.GetBlock") (param $addr i32) (result i32 i32 i32)
  local.get $addr
  i32.load offset=0 align=1
  local.get $addr
  i32.load offset=4 align=1
  local.get $addr
  i32.load offset=8 align=1
) ;;__main__.GetBlock

(func $__main__.SetBlock (export "__main__.SetBlock") (param $addr i32) (param $data.size i32) (param $data.used i32) (param $data.next i32)
  local.get $addr
  local.get $data.size
  i32.store offset=0 align=1
  local.get $addr
  local.get $data.used
  i32.store offset=4 align=1
  local.get $addr
  local.get $data.next
  i32.store offset=8 align=1
) ;;__main__.SetBlock
```

若没有这两个标记，通过内存地址直接读写结构体需要手动计算各个成员的偏移值，这两个编译标记可以简化该操作，有助于提高底层代码的开发效率和可维护性。