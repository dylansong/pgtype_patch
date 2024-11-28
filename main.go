package main

import (
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

func main() {
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

	// 1. 创建 rows 目录和 rows.go 文件
	err = os.MkdirAll("db/rows", 0755)
	if err != nil {
		fmt.Printf("创建 rows 目录失败: %v\n", err)
		return
	}

	// 创建并写入 rows.go 文件头
	rowsFile, err := os.Create("db/rows/rows.go")
	if err != nil {
		fmt.Printf("创建 rows.go 失败: %v\n", err)
		return
	}
	defer rowsFile.Close()

	// 写入包声明和导入语句
	_, err = rowsFile.WriteString("package r\n\n\n")
	if err != nil {
		fmt.Printf("写入 rows.go 文件头失败: %v\n", err)
		return
	}

	// 3. 读取所有 *.sql.go 文件
	files, err := filepath.Glob("db/*.sql.go")
	if err != nil {
		fmt.Printf("查找sql.go文件失败: %v\n", err)
		return
	}

	// 用于匹配结构体定义的正则表达式
	structRegex := regexp.MustCompile(`type\s+\w+Params\s+struct\s*{[^}]+}`)
	rowStructRegex := regexp.MustCompile(`type\s+\w+Row\s+struct\s*{[^}]+}`)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("读取文件 %s 失败: %v\n", file, err)
			continue
		}

		// 查找所有匹配的 Params 结构体定义
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

		// 查找所有匹配的 Row 结构体定义
		rowMatches := rowStructRegex.FindAll(content, -1)
		for _, match := range rowMatches {
			// 将找到的结构体写入 rows.go
			_, err = rowsFile.Write(match)
			if err != nil {
				fmt.Printf("写入 Row 结构体失败: %v\n", err)
				continue
			}
			_, err = rowsFile.WriteString("\n\n")
			if err != nil {
				fmt.Printf("写入换行符失败: %v\n", err)
				continue
			}
		}
	}

	// 4. 读取整个文件内容
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

	// 处理 rows.go 文件
	rowsFile.Seek(0, 0)
	rowsContent, err := ioutil.ReadAll(rowsFile)
	if err != nil {
		fmt.Printf("读取 rows.go 失败: %v\n", err)
		return
	}

	// 定义类型替换映射
	typeReplacements := map[string]string{
		"pgtype.Text":        "string",
		"pgtype.UUID":        "[16]byte",
		"pgtype.Uint32":      "uint32",
		"pgtype.Timestamptz": "time.Time",
		"pgtype.Timestamp ":  "time.Time",
		"pgtype.Time ":       "int64",
		"pgtype.Int2":        "int16",
		"pgtype.Float8":      "float64",
		"pgtype.Float4":      "float32",
		"pgtype.Date":        "time.Time",
		"pgtype.Bool":        "bool",
		"pgtype.Bits":        "[]byte",
	}

	// 修改替换逻辑
	newRowsContent := string(rowsContent)
	// 使用正则表达式进行更精确的替换
	newRowsContent = regexp.MustCompile(`pgtype\.Timestamptz\b`).ReplaceAllString(newRowsContent, "time.Time")
	newRowsContent = regexp.MustCompile(`pgtype\.Timestamp\b`).ReplaceAllString(newRowsContent, "time.Time")
	newRowsContent = regexp.MustCompile(`pgtype\.Time\b`).ReplaceAllString(newRowsContent, "int64")
	
	// 执行其他类型替换
	for old, new := range typeReplacements {
		if !strings.Contains(old, "Timestamp") && !strings.Contains(old, "Time") {
			newRowsContent = strings.ReplaceAll(newRowsContent, old, new)
		}
	}

	// 重写文件
	err = ioutil.WriteFile("db/rows/rows.go", []byte(newRowsContent), 0644)
	if err != nil {
		fmt.Printf("重写 rows.go 失败: %v\n", err)
		return
	}

	// 复制 models.go 到 db/rows 目录
	modelsSrc := "db/models.go"
	modelsDst := "db/rows/models.go"
	
	// 读取源文件内容
	modelsContent, err := ioutil.ReadFile(modelsSrc)
	if err != nil {
		fmt.Printf("读取 models.go 失败: %v\n", err)
		return
	}

	// 替换 package db 为 package r
	newModelsContent := strings.Replace(string(modelsContent), "package db", "package r", 1)
	// 删除 pgtype 导入
	newModelsContent = regexp.MustCompile(`\s*"github\.com/jackc/pgx/v5/pgtype"\n`).ReplaceAllString(newModelsContent,  "\n    \"time\"\n")



	// 执行所有类型替换
	for old, new := range typeReplacements {
		newModelsContent = strings.ReplaceAll(newModelsContent, old, new)
	}

	// 写入新文件
	err = ioutil.WriteFile(modelsDst, []byte(newModelsContent), 0644)
	if err != nil {
		fmt.Printf("写入 db/rows/models.go 失败: %v\n", err)
		return
	}

	fmt.Println("处理完成！")
}