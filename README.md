# Agent OS SDK

多个服务共用的 Go 协议与传输库，包含版本化消息契约、有界 HTTP、MCP 传输/网络安全、资源目录 DTO 和可配置的进程托管。最初从集成提交 `400afaea7ed6c9c4aa47b8002839efae2d6bf55d` 拆出。

## 边界

SDK 不包含业务 SQL、迁移、记录仓库、角色授权表或角色到 schema 的探测规则。各服务入口通过 `servicehost.Options` 提供自己的就绪与鉴权策略；连接池属于基础设施，业务查询留在所属服务。`contracts` 只定义协议结构与输入约束，不执行业务流程。

此仓库是依赖库，**不是独立运行的微服务**，没有业务端口或服务专属默认配置。`RUNTIME_URL` 可开启实例心跳，省略则关闭；密钥与地址由使用方部署配置提供。

`feat/model-native-limits` 增加可选的 `ModelConfig.Limits` 协议候选：保存具体模型、供应商端点、协议、物理窗口/输出范围及观测来源。空指针保留旧默认及验证行为；提供新能力时必须与任务配置一致，不再套用旧32K输出验证和小上下文额度。Registry仍须服务端解析能力，不能信任客户端自填；本SDK不访问供应商。32条本地测试及vet通过，但服务集成、新Harness及正式发布尚未完成，不能将此协议候选说成运行中的限制已取消。

## 构建与验证

需要 Go 1.27.2。

```powershell
$env:GOWORK='off'
go test -count=1 ./...
go vet ./...
```

SDK 测试不需要数据库或模型密钥；具体服务集成测试由各服务负责。Linux race 与 Windows 检查由 CI 定义，当前 Actions 日常关闭，等待统一验收。

## 版本分发

`v0.1.0-rc.4` 增加可选 `Agent.final_output`：在不可变快照中指定候选工具、顶层输入字符串字段和响应的对象键路径。只有路径值为 JSON 布尔 `true` 才可作为候选；这不表示事实已核实。旧 Agent 不增加 JSON 字段。引擎所有者还须校验工具绑定、原始字符串、当前约束和版本兼容；本 SDK 不执行成稿。协议测试与 vet 已通过，新 Harness / Runtime 的集成和正式安装仍待完成。

使用方固定模块版本并提交需要的 `vendor` 源码和许可证，以 `GOWORK=off`、`-mod=vendor` 独立构建，不依赖个人令牌或邻近目录 `replace`。集成发布记录精确 SDK 提交与模块校验和。行为变化需发布新版本，不能覆盖既有版本。

早期手工生成的 rc.1/rc.2 本地代理包曾因 Windows 换行和 ZIP 目录条目产生非标准校验和，已弃用。应针对固定 Git tag 执行 `go mod download`，再 `go mod vendor`，不要手工拼模块 ZIP。

rc.2 标准模块校验和为 `h1:SunD5Oyh1le4sv3CMJTpKCa0SsKBQDqXoBoAC0aji98=`，对应提交 `f45ca168dee81f864533a8f77f8c203960caad1a`；原始标签和历史保留。[历史说明](docs/history/service-notes-20261014.en.md) 记录来源及分发约束。
