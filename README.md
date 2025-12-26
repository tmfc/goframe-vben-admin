# Repository Guidelines

基于 GoFrame 和 Vben Admin 的 SaaS 管理平台。下面给出后端与前端的开发环境配置与运行方式。

## 环境准备

- Go: 1.24.11（或兼容 1.24.x）
- Node.js: >= 20.12.0
- pnpm: >= 10.0.0
- PostgreSQL: 本地默认连接 `127.0.0.1:5433`，数据库名 `gva`

后端配置在 `backend/manifest/config/config.yaml`，按需修改数据库连接。

## 后端开发与运行

在 `backend/` 目录执行：

```bash
go mod tidy
go run ./main.go
```

## 数据库迁移

在项目根目录执行：

```bash
migrate -path backend/db/migrations -database "postgres://gva:gva@127.0.0.1:5433/gva?sslmode=disable" up
```

连接字符串请按 `backend/config.toml` 调整。

## 运行后端测试

在 `backend/` 目录执行：

```bash
go test ./...
```

常用命令：

- `go test ./...` 运行后端测试
- `make build` 使用 GoFrame 构建二进制（需安装 `gf` CLI）

## 前端开发与运行

在 `frontend/` 目录执行：

```bash
pnpm install
pnpm dev
```

常用应用入口（任选一个）：

- `pnpm dev:antd` 启动 Ant Design 版本
- `pnpm dev:naive` 启动 Naive UI 版本
- `pnpm dev:ele` 启动 Element Plus 版本
- `pnpm dev:tdesign` 启动 TDesign 版本

## 联调建议

- 先启动后端，再启动前端，确保 API 可访问。
- 如需调整 API 地址，按前端应用的环境配置文件修改（对应应用目录内的 `.env*`）。
