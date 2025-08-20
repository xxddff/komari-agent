# komari-agent

监控代理，用于收集系统信息并向监控服务器报告。

## Cloudflare Access 支持

如果您的监控服务器受到 Cloudflare Access 保护，您可以配置 agent 添加必要的访问头部。

### 配置方式

#### 1. 环境变量（推荐）

```bash
export KOMARI_CF_ACCESS_CLIENT_ID="your-client-id"
export KOMARI_CF_ACCESS_CLIENT_SECRET="your-client-secret"
```

#### 2. 命令行参数

```bash
./komari-agent --cf-access-client-id="your-client-id" --cf-access-client-secret="your-client-secret" --endpoint="https://your-server.com" --token="your-token"
```

### 获取 Cloudflare Access 凭据

1. 登录到 Cloudflare Dashboard
2. 进入 Zero Trust > Access > Service Tokens
3. 创建新的 Service Token
4. 复制 Client ID 和 Client Secret
5. 将其配置到 komari-agent 中

### 功能说明

配置后，agent 会在以下请求中自动添加 `CF-Access-Client-Id` 和 `CF-Access-Client-Secret` 头部：

- 基本信息上报 (`/api/clients/uploadBasicInfo`)
- 任务结果上报 (`/api/clients/task/result`)
- WebSocket 连接 (`/api/clients/report`)
- 终端连接 (`/api/clients/terminal`)
- 自动发现注册 (`/api/clients/register`)

注意：HTTP ping 测试和外部 IP 查询不会添加这些头部，因为它们访问的是外部服务。