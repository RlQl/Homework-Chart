# MedChart

本项目严格遵循软件开发标准流程，实现医疗时序数据分析与可视化。

## 技术栈选型
- **后端**：Go (处理底层二进制解析、高频并发模拟、LTTB 数据降采样)
- **前端**：React + TypeScript + Vite (UI 交互与状态管理)
- **可视化**：Apache ECharts (Canvas 高性能波形渲染、DataZoom 交互)
- **桌面框架**：Wails v2 (提供原生系统能力与高效 IPC 通信)
- **工程化**：GitHub CI/CD (自动化构建与测试)、严格的分层架构设计

## 工程结构
- `/core`：核心业务逻辑（独立纯 Go 代码，负责解析与算法，完全覆盖单元测试）
- `/docs`：构造管理文档（需求说明、架构评估、源码结构、测试报告）
- `/frontend`：前端工程目录
- `app.go` / `main.go`：Wails 进程与生命周期入口

## 快速启动
```bash
# 开发模式热重载
wails dev

# 编译独立可执行文件
wails build
