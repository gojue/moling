## MoLing MCP 服务器

[English](./README.md) | 中文

[![GitHub stars](https://img.shields.io/github/stars/gojue/moling.svg?label=Stars&logo=github)](https://github.com/gojue/moling)
[![GitHub forks](https://img.shields.io/github/forks/gojue/moling?label=Forks&logo=github)](https://github.com/gojue/moling)
[![CI](https://github.com/gojue/moling/actions/workflows/go-test.yml/badge.svg)](https://github.com/gojue/ecapture/actions/workflows/go-test.yml)
[![Github Version](https://img.shields.io/github/v/release/gojue/moling?display_name=tag&include_prereleases&sort=semver)](https://github.com/gojue/moling/releases)

---

### 简介
MoLing是一个computer-use的MCP Server，基于操作系统API实现了系统交互，可以实现文件系统的读写、合并、统计、聚合等操作，也可以执行系统命令操作。是一个无需任何依赖的本地办公自动化助手。

### 优势
> [!IMPORTANT]
> 没有任何安装依赖，直接运行，兼容Windows、Linux、macOS等操作系统。
> 再也不用苦恼NodeJS、Python等环境冲突等问题。

### 功能特性

- **文件系统操作**：读取、写入、合并、统计和聚合
- **命令行终端**：直接执行系统命令
- **浏览器控制**：基于 `github.com/chromedp/chromedp`
- **未来计划**：
    - 个人电脑资料整理
    - 文档编写辅助
    - 行程规划
    - 生活助手功能

### 支持的MCP客户端

- Claude
- Cline
- Cherry Studio

#### 配置格式
```json

```
### 运行模式

- **Stdio模式**：本地命令行交互模式，依赖于终端输入输出，适合人机交互
- **SSR模式**：远程通讯模式，适合远程部署，远程调用

### 安装指南


#### 方法一： 脚本安装
#### Linux/MacOS
```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.sh)"
```
##### Windows
```powershell
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.ps1 | iex"
```


#### 方法二：直接下载
1. 从[发布页面](https://github.com/gojue/moling/releases)下载安装包
2. 解压安装包
3. 运行服务器：
```sh
./moling
```

#### 方法三：从源码编译
1. 克隆代码库：
```sh
git clone https://github.com/gojue/moling.git
cd moling
```
2. 编译项目（需要Golang工具链）：
```sh
make build
```
3. 运行编译后的程序：
```sh
./bin/moling
```

### 使用说明
启动服务器后，使用任何支持的MCP客户端配置连接到您的MoLing服务器地址即可。

### 许可证
Apache License 2.0。详见[LICENSE](LICENSE)文件。