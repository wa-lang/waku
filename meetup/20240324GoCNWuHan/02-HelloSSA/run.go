package main

import (
	"bytes"
	"fmt"
	"go/constant"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
)

func runFunc(fn *ssa.Function) {
	fmt.Println("--- runFunc begin ---")
	defer fmt.Println("--- runFunc end   ---")

	// 从第0个Block开始执行
	if len(fn.Blocks) > 0 {
		for blk := fn.Blocks[0]; blk != nil; {
			blk = runFuncBlock(fn, fn.Blocks[0])
		}
	}
}

// 运行Block, 返回下一个Block, 如果返回nil表示结束
func runFuncBlock(fn *ssa.Function, block *ssa.BasicBlock) (nextBlock *ssa.BasicBlock) {
	for _, ins := range block.Instrs {
		switch ins := ins.(type) {
		case *ssa.Call:
			doCall(ins)
		case *ssa.Return:
			doReturn(ins)
		default:
			panic("Not Implemented.")
		}
	}
	return nil
}

func doCall(ins *ssa.Call) {
	switch {
	case ins.Call.Method == nil: // 普通函数调用
		switch callFn := ins.Call.Value.(type) {
		case *ssa.Builtin:
			callBuiltin(callFn, ins.Call.Args...)
		default:
			// 普通函数
			panic("Not Implemented.")
		}

	default:
		// 方法或接口调用
		panic("Not Implemented.")
	}
}

func callBuiltin(fn *ssa.Builtin, args ...ssa.Value) {
	switch fn.Name() {
	case "println":
		var buf bytes.Buffer
		for i := 0; i < len(args); i++ {
			if i > 0 {
				buf.WriteRune(' ')
			}
			switch arg := args[i].(type) {
			case *ssa.Const: // 处理常量参数
				if t, ok := arg.Type().Underlying().(*types.Basic); ok {
					switch t.Kind() {
					case types.Int, types.UntypedInt:
						fmt.Fprintf(&buf, "%d", int(arg.Int64()))
					case types.String:
						fmt.Fprintf(&buf, "%s", constant.StringVal(arg.Value))
					default:
						// 其它常量类型，暂不支持
						panic("Not Implemented.")
					}
				}
			default:
				// 暂不支持非常量参数
				panic("Not Implemented.")
			}
		}
		buf.WriteRune('\n')
		os.Stdout.Write(buf.Bytes())

	default:
		// 其它内置函数，暂不支持
		panic("Not Implemented.")
	}
}

func doReturn(ins *ssa.Return) {
	return // ins.Results[...]
}
