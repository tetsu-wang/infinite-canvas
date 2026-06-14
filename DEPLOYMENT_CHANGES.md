# Mac Docker 部署配置修改摘要

## 修改时间
2026-06-14

## 修改目的
将项目配置为可在 macOS 上使用 Docker 进行本地部署，特别优化了 Apple Silicon (M1/M2/M3) 支持。

## 详细修改清单

### 1. 修复 Dockerfile (关键)
**文件**: `Dockerfile`
**行号**: 13
**修改内容**: 
- 修改前: `FROM golang:1.25-alpine AS api-build`
- 修改后: `FROM golang:1.23-alpine AS api-build`
- **原因**: Go 1.25 版本不存在，修正为当前稳定版本 1.23

### 2. 优化 docker-compose.local.yml
**文件**: `docker-compose.local.yml`
**新增配置**:
```yaml
platform: linux/arm64              # ARM64 架构支持（M1/M2/M3）
container_name: infinite-canvas-local  # 固定容器名称
build:
  platforms:
    - linux/arm64                  # 明确指定构建平台
healthcheck:                        # 健康检查
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
environment:
  - TZ=Asia/Shanghai               # 设置时区
```

### 3. 配置环境变量文件
**文件**: `.env` (从 `.env.example` 创建)
**关键修改**:
- `DATABASE_DSN` 改为 `/app/data/infinite-canvas.db` (Docker 容器内绝对路径)
- 保留默认管理员账号密码（部署后需修改）

### 4. 创建数据目录
**目录**: `data/`
**用途**: 持久化存储数据库和用户上传文件

### 5. 新增部署脚本
**文件**: `deploy-mac.sh` (可执行)
**功能**:
- 自动检查 Docker 运行状态
- 检测系统架构 (ARM64/AMD64)
- 创建必要的配置文件和目录
- 停止并删除旧容器
- 构建并启动新容器
- 显示部署结果和访问信息

### 6. 新增文档
**文件**: 
- `docs/mac-deployment.md` - 完整部署文档（故障排查、优化建议）
- `MAC_QUICKSTART.md` - 快速开始指南

### 7. 优化 .dockerignore
**新增忽略项**:
```
*.md (除了 VERSION 和 CHANGELOG.md)
.DS_Store
.vscode
*.swp, *.swo, *~
```

## 部署流程

### 快速部署
```bash
./deploy-mac.sh
```

### 手动部署
```bash
docker compose -f docker-compose.local.yml up -d --build
```

## 架构兼容性

### Apple Silicon (M1/M2/M3) ✅
- 原生支持 ARM64 架构
- 已在配置中指定 `platform: linux/arm64`

### Intel Mac ✅
- 如需切换到 AMD64，修改 `docker-compose.local.yml`:
  ```yaml
  platform: linux/amd64
  ```

## 访问信息

- **地址**: http://localhost:3000
- **默认账号**: admin
- **默认密码**: infinite-canvas

⚠️ **重要**: 首次登录后立即修改密码！

## 常用命令

```bash
# 查看容器状态
docker ps

# 查看日志
docker logs -f infinite-canvas-local

# 停止服务
docker compose -f docker-compose.local.yml down

# 重启服务
docker compose -f docker-compose.local.yml restart

# 重新构建
docker compose -f docker-compose.local.yml up -d --build --force-recreate
```

## 故障排查要点

1. **端口占用**: 使用 `lsof -i :3000` 检查
2. **内存不足**: Docker Desktop 设置至少 4GB 内存
3. **权限问题**: `sudo chown -R $(whoami) data/`
4. **架构不匹配**: 设置 `export DOCKER_DEFAULT_PLATFORM=linux/arm64`

## 数据备份

```bash
# 备份
tar -czf infinite-canvas-backup-$(date +%Y%m%d).tar.gz data/

# 恢复
docker compose -f docker-compose.local.yml down
tar -xzf infinite-canvas-backup-YYYYMMDD.tar.gz
docker compose -f docker-compose.local.yml up -d
```

## 技术栈

- **前端**: Next.js 16.2 + React + TypeScript + Tailwind CSS
- **后端**: Go 1.23 + Gin + GORM
- **构建工具**: Bun 1.3.13 (前端) + Go (后端)
- **运行时**: Node.js 22 (容器最终运行环境)
- **数据库**: SQLite (默认)

## 相关文档

- 详细部署指南: [docs/mac-deployment.md](docs/mac-deployment.md)
- 快速开始: [MAC_QUICKSTART.md](MAC_QUICKSTART.md)
- 功能介绍: [docs/features.md](docs/features.md)
- 项目主文档: [README.md](README.md)

## 未修改的文件

以下文件保持原样：
- `docker-compose.yml` - 生产环境配置（使用远程镜像）
- `main.go` - Go 后端入口
- `web/` - 前端源码
- 其他业务逻辑代码

## 验证清单

部署前检查：
- [ ] Docker Desktop 已安装并运行
- [ ] 至少 4GB 可用内存
- [ ] 至少 2GB 可用磁盘空间
- [ ] 端口 3000 未被占用

部署后验证：
- [ ] 容器正常运行 (`docker ps`)
- [ ] 健康检查通过 (`docker inspect infinite-canvas-local`)
- [ ] 可访问 http://localhost:3000
- [ ] 可正常登录
- [ ] 数据库文件已创建 (`ls data/`)

## 已知限制

1. 首次构建需要 5-10 分钟（下载依赖、编译）
2. SQLite 不支持高并发写入（多用户场景建议切换 MySQL/PostgreSQL）
3. 容器内服务以非特权用户运行，部分系统操作受限

## 后续优化建议

1. **生产部署**: 修改 JWT_SECRET 和管理员密码
2. **性能优化**: 考虑使用 MySQL 或 PostgreSQL
3. **反向代理**: 使用 Nginx/Caddy 处理 HTTPS
4. **监控**: 接入 Prometheus/Grafana
5. **备份**: 设置定时备份任务

---

*最后更新: 2026-06-14*
*修改者: Claude Code*
