# pgtype_patch 工具使用说明

## 功能介绍
这是一个用于处理 PostgreSQL 类型定义的工具，主要提供两个功能：
1. 生成 Go 语言的 pgtype 类型定义和参数文件
2. 更新 TypeScript 的类型定义文件

## 安装方法

go install github.com/dylansong/pgtype_patch@v1.0.6

## 使用方法
工具提供两个命令模式，通过 `-cmd` 参数指定：

### 1. pgtype 模式


```bash
pgtype_patch -cmd pgtype
```


执行效果：
- 在 `db` 目录下创建 `pgtype.go` 文件
- 在 `db/params` 目录下创建并处理 `params.go` 文件
- 自动处理所有 `*.sql.go` 文件中的 Params 结构体

### 2. TypeScript 模式

pgtype_patch -cmd ts


执行效果：
- 在 `src/lib/encore/generated.ts` 文件中：
  - 将原有的 `export namespace pgtype` 重命名为 `export namespace pgtypeBak`
  - 在文件末尾添加新的 pgtype 类型定义

## 命令行参数

```bash
-cmd string
    可选值: "pgtype" 或 "ts"
    默认值: "pgtype"
```


## 发布新版本
如需发布新版本，执行以下命令：

```bash
git tag v1.0.6
git push origin v1.0.6
```


## 目录结构要求
使用工具时，请确保项目符合以下目录结构：
```
项目根目录
├── db/
│   ├── *.sql.go
│   └── params/
└── src/
    └── lib/
        └── encore/
            └── generated.ts
```


## 注意事项
1. 运行工具前请确保有相应目录的写入权限
2. 建议在执行操作前备份重要文件
3. TypeScript 模式会修改 `generated.ts` 文件，请确保该文件存在
```

这个 README.md 文件包含了完整的工具使用说明，包括安装、使用方法、参数说明以及注意事项。你可以将它放在项目的根目录下。需要我修改或补充什么内容吗？


