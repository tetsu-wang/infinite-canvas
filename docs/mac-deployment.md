# Mac 本地 Docker 部署指南

本文档说明如何在 macOS 上使用 Docker 部署 infinite-canvas 项目。

## 系统要求

- macOS 10.15+
- Docker Desktop 4.0+
- 至少 4GB 可用内存
- 至少 2GB 可用磁盘空间

## 快速开始

### 方式一：使用部署脚本（推荐）

```bash
./deploy-mac.sh
```

脚本会自动：
- 检查 Docker 是否运行
- 创建配置文件和数据目录
- 构建并启动容器
- 显示访问信息

### 方式二：手动部署

```bash
# 1. 创建配置文件
cp .env.example .env

# 2. 编辑配置（可选）
nano .env

# 3. 创建数据目录
mkdir -p data

# 4. 构建并启动
docker compose -f docker-compose.local.yml up -d --build

# 5. 查看日志
docker logs -f infinite-canvas-local
```

## 访问应用

部署成功后访问：http://localhost:3000

默认登录信息：
- 用户名: `admin`
- 密码: `infinite-canvas`

**⚠️ 首次登录后请立即修改密码！**

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

# 清理数据（⚠️ 会删除所有数据）
docker compose -f docker-compose.local.yml down -v
rm -rf data/*
```

## 故障排查

### 1. 端口 3000 被占用

**症状**: `bind: address already in use`

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :3000

# 杀死进程（替换 <PID> 为实际进程 ID）
kill -9 <PID>

# 或修改端口
# 编辑 docker-compose.local.yml，将 "3000:3000" 改为 "3001:3000"
```

### 2. Docker 内存不足

**症状**: 构建过程中卡住或失败

**解决方案**:
1. 打开 Docker Desktop
2. Settings → Resources
3. 将 Memory 调整至至少 4GB
4. 点击 Apply & Restart

### 3. 数据库文件权限问题

**症状**: `unable to open database file`

**解决方案**:
```bash
# 重置数据目录权限
sudo chown -R $(whoami) data/
chmod -R 755 data/
```

### 4. Bun 安装依赖失败

**症状**: `error: failed to install dependencies`

**解决方案**:
```bash
# 清理 Docker 缓存重新构建
docker builder prune -a
docker compose -f docker-compose.local.yml up -d --build --no-cache
```

### 5. 容器启动后立即退出

**症状**: `docker ps` 看不到运行中的容器

**解决方案**:
```bash
# 查看容器退出原因
docker logs infinite-canvas-local

# 查看所有容器（包括已停止的）
docker ps -a

# 如果是配置问题，检查 .env 文件
cat .env
```

### 6. 健康检查失败

**症状**: 容器状态显示 `unhealthy`

**解决方案**:
```bash
# 查看健康检查日志
docker inspect infinite-canvas-local --format='{{json .State.Health}}'

# 手动测试健康检查端点
docker exec infinite-canvas-local wget -q -O- http://localhost:3000/api/health
```

## 架构兼容性

### Apple Silicon (M1/M2/M3)

项目已配置为使用 `linux/arm64` 平台，原生支持 Apple Silicon。

如果遇到架构问题：
```bash
# 强制使用 ARM64
export DOCKER_DEFAULT_PLATFORM=linux/arm64
docker compose -f docker-compose.local.yml up -d --build
```

### Intel Mac

如需在 Intel Mac 上运行，修改 `docker-compose.local.yml`:
```yaml
services:
  app:
    platform: linux/amd64  # 改为 amd64
```

## 性能优化

### 1. 启用文件共享缓存

Docker Desktop → Settings → Resources → File Sharing

确保项目目录在共享列表中。

### 2. 使用 Docker 卷代替绑定挂载

如果遇到性能问题，可以修改 `docker-compose.local.yml`:
```yaml
volumes:
  - canvas-data:/app/data  # 使用命名卷

volumes:
  canvas-data:
```

### 3. 调整资源限制

修改 `docker-compose.local.yml`:
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          memory: 2G
```

## 数据备份

### 备份数据

```bash
# 备份整个数据目录
tar -czf infinite-canvas-backup-$(date +%Y%m%d).tar.gz data/

# 仅备份数据库
cp data/infinite-canvas.db infinite-canvas-backup-$(date +%Y%m%d).db
```

### 恢复数据

```bash
# 停止容器
docker compose -f docker-compose.local.yml down

# 恢复数据
tar -xzf infinite-canvas-backup-20260614.tar.gz

# 重启容器
docker compose -f docker-compose.local.yml up -d
```

## 卸载

```bash
# 停止并删除容器
docker compose -f docker-compose.local.yml down

# 删除镜像
docker rmi infinite-canvas:local

# 删除数据（可选）
rm -rf data/

# 删除配置（可选）
rm .env
```

## 常见问题

**Q: 为什么首次启动很慢？**

A: 首次构建需要下载基础镜像、安装依赖、编译代码，通常需要 5-10 分钟。后续启动会快很多。

**Q: 可以同时运行多个实例吗？**

A: 可以，但需要修改端口映射，避免冲突。

**Q: 数据存储在哪里？**

A: 数据存储在项目根目录的 `data/` 文件夹中，会自动挂载到容器内的 `/app/data`。

**Q: 如何更新到最新版本？**

A: 
```bash
git pull
docker compose -f docker-compose.local.yml up -d --build
```

## 开发模式

如需进行开发调试：

```bash
# 进入容器
docker exec -it infinite-canvas-local sh

# 查看进程
ps aux

# 查看环境变量
env | grep -E 'PORT|DATABASE|ADMIN'
```

## 相关文档

- [功能介绍](features.md)
- [部署说明](deployment.md)
- [后端数据库说明](backend-database.md)
