# Mac 本地部署快速开始

## 一键部署

```bash
./deploy-mac.sh
```

部署脚本会自动完成所有配置，大约 5-10 分钟后可以访问 http://localhost:3000

## 手动部署

```bash
# 1. 确保 Docker Desktop 已启动
docker info

# 2. 构建并启动
docker compose -f docker-compose.local.yml up -d --build

# 3. 查看启动日志
docker logs -f infinite-canvas-local

# 4. 等待服务完全启动（看到 "Server started" 或类似消息）
# 然后访问 http://localhost:3000
```

## 默认登录

- 用户名: `admin`
- 密码: `infinite-canvas`

⚠️ **首次登录后请立即修改密码！**

## 常用命令

```bash
# 查看运行状态
docker ps

# 查看日志
docker logs -f infinite-canvas-local

# 重启
docker compose -f docker-compose.local.yml restart

# 停止
docker compose -f docker-compose.local.yml down

# 完全重建（清理缓存）
docker compose -f docker-compose.local.yml down
docker compose -f docker-compose.local.yml up -d --build --force-recreate
```

## 故障排查

### 端口占用

如果 3000 端口被占用：

```bash
# 查找占用进程
lsof -i :3000

# 杀死进程
kill -9 <PID>
```

或者修改 `docker-compose.local.yml` 中的端口映射：

```yaml
ports:
  - "3001:3000"  # 改用 3001 端口
```

### 查看详细日志

```bash
# 实时日志
docker logs -f infinite-canvas-local

# 最近 100 行
docker logs --tail 100 infinite-canvas-local

# 带时间戳
docker logs -f --timestamps infinite-canvas-local
```

### 数据库问题

如果遇到数据库错误，检查数据目录权限：

```bash
ls -la data/
sudo chown -R $(whoami) data/
```

### 重置所有数据

```bash
docker compose -f docker-compose.local.yml down
rm -rf data/*
docker compose -f docker-compose.local.yml up -d
```

## 配置修改

所有配置在 `.env` 文件中：

```bash
# 编辑配置
nano .env

# 修改后重启生效
docker compose -f docker-compose.local.yml restart
```

## 更多帮助

详细文档: [docs/mac-deployment.md](docs/mac-deployment.md)
