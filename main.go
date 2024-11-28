package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const pgtypeContent = `package db

import (
	"time"
)

type InfinityModifier int8

// Bits represents the PostgreSQL bit and varbit types.
type Bits struct {
	Bytes []byte
	Len   int32 // Number of bits
	Valid bool
}

type Bool struct {
	Bool  bool
	Valid bool
}

type Date struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

type Vec2 struct {
	X float64
	Y float64
}

type Box struct {
	P     [2]Vec2
	Valid bool
}

type Circle struct {
	P     Vec2
	R     float64
	Valid bool
}

type Float4 struct {
	Float32 float32
	Valid   bool
}

type Float8 struct {
	Float64 float64
	Valid   bool
}

type Int2 struct {
	Int16 int16
	Valid bool
}

type Interval struct {
	Microseconds int64
	Days         int32
	Months       int32
	Valid        bool
}

type JSONCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

type JSONBCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

type Line struct {
	A, B, C float64
	Valid   bool
}

type Lseg struct {
	P     [2]Vec2
	Valid bool
}

type Path struct {
	P      []Vec2
	Closed bool
	Valid  bool
}

type Polygon struct {
	P     []Vec2
	Valid bool
}

type Text struct {
	String string
	Valid  bool
}

type TID struct {
	BlockNumber  uint32
	OffsetNumber uint16
	Valid        bool
}

type Time struct {
	Microseconds int64 // Number of microseconds since midnight
	Valid        bool
}

// Timestamp represents the PostgreSQL timestamp type.
type Timestamp struct {
	Time             time.Time // Time zone will be ignored when encoding to PostgreSQL.
	InfinityModifier InfinityModifier
	Valid            bool
}

// Timestamptz represents the PostgreSQL timestamptz type.
type Timestamptz struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

// Uint32 is the core type that is used to represent PostgreSQL types such as OID, CID, and XID.
type Uint32 struct {
	Uint32 uint32
	Valid  bool
}

type UUID struct {
	Bytes [16]byte
	Valid bool
}

type XMLCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}`

const tsTypeContent = `    export type InfinityModifier = number

    export type Text = string | null

    export type Timestamp = string | null

    export type Time = number | null

    export type Timestamptz = string | null

    export type Bool = boolean | null

    export type Date = string | null

    export type Float4 = number | null

    export type Float8 = number | null

    export type Int2 = number | null	

    export type JSONCodec = any | null

    export type JSONBCodec = any | null

    export type Interval = number | null

    export type Uint32 = number | null

    export type UUID = string | null

    export type XMLCodec = any | null`

func main() {
	// 添加命令行参数解析
	command := flag.String("cmd", "pgtype", "要执行的命令: pgtype 或 ts")
	flag.Parse()

	switch *command {
	case "pgtype":
		executePgtypeTask()
	case "ts":
		executeTypeScriptTask()
	default:
		fmt.Printf("未知的命令: %s\n", *command)
		fmt.Println("可用命令: pgtype, ts")
	}
}

// 将原来的 main 函数内容移到这个新函数中
func executePgtypeTask() {
	// 1. 创建 pgtype.go 文件
	err := ioutil.WriteFile("db/pgtype.go", []byte(pgtypeContent), 0644)
	if err != nil {
		fmt.Printf("创建pgtype.go失败: %v\n", err)
		return
	}

	// 2. 创建 params 目录和 params.go 文件
	err = os.MkdirAll("db/params", 0755)
	if err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return
	}

	// 创建并写入 params.go 文件头
	paramsFile, err := os.Create("db/params/params.go")
	if err != nil {
		fmt.Printf("创建params.go失败: %v\n", err)
		return
	}
	defer paramsFile.Close()

	// 写入包声明和导入语句
	_, err = paramsFile.WriteString("package p\n\nimport \"encore.app/db\"\n\n")
	if err != nil {
		fmt.Printf("写入文件头失败: %v\n", err)
		return
	}

	// 3. 读取所有 *.sql.go 文件
	files, err := filepath.Glob("db/*.sql.go")
	if err != nil {
		fmt.Printf("查找sql.go文件失败: %v\n", err)
		return
	}

	// 只保留 Params 结构体的正则表达式
	structRegex := regexp.MustCompile(`type\s+\w+Params\s+struct\s*{[^}]+}`)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("读取文件 %s 失败: %v\n", file, err)
			continue
		}

		// 只查找 Params 结构体定义
		matches := structRegex.FindAll(content, -1)
		for _, match := range matches {
			// 将找到的结构体写入 params.go
			_, err = paramsFile.Write(match)
			if err != nil {
				fmt.Printf("写入结构体失败: %v\n", err)
				continue
			}
			_, err = paramsFile.WriteString("\n\n")
			if err != nil {
				fmt.Printf("写入换行符失败: %v\n", err)
				continue
			}
		}
	}

	// 4. 读取整个文件内容并处理 params.go
	paramsFile.Seek(0, 0)
	content, err := ioutil.ReadAll(paramsFile)
	if err != nil {
		fmt.Printf("读取params.go失败: %v\n", err)
		return
	}

	// 替换 pgtype. 为空字符串
	newContent := strings.ReplaceAll(string(content), "pgtype.", "db.")

	// 重写文件
	err = ioutil.WriteFile("db/params/params.go", []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("重写params.go失败: %v\n", err)
		return
	}

	fmt.Println("处理完成！")
}

// 新增的 TypeScript 相关任务函数
func executeTypeScriptTask() {
	filePath := "src/lib/encore/generated.ts"
	
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 使用更复杂的逻辑来处理嵌套的花括号
	startPattern := regexp.MustCompile(`export\s+namespace\s+pgtype\s*\{`)
	contentStr := string(content)
	startIdx := startPattern.FindStringIndex(contentStr)
	
	if startIdx == nil {
		fmt.Println("未找到 pgtype namespace")
		return
	}

	// 从找到的位置开始计算花括号的配对
	braceCount := 0
	endIdx := startIdx[1]
	
	for i := startIdx[1]; i < len(contentStr); i++ {
		if contentStr[i] == '{' {
			braceCount++
		} else if contentStr[i] == '}' {
			braceCount--
			if braceCount == 0 {
				endIdx = i + 1
				break
			}
		}
	}

	// 构建新的内容
	newContent := contentStr[:startIdx[0]] +
		fmt.Sprintf("export namespace pgtype {\n%s\n}", tsTypeContent) +
		contentStr[endIdx:]

	// 写回文件
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Println("TypeScript 类型定义更新完成！")
}