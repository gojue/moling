## MoLing MCP Server

English | [中文](./README_ZH_HANS.md)

[![GitHub stars](https://img.shields.io/github/stars/gojue/moling.svg?label=Stars&logo=github)](https://github.com/gojue/moling)
[![GitHub forks](https://img.shields.io/github/forks/gojue/moling?label=Forks&logo=github)](https://github.com/gojue/moling)
[![CI](https://github.com/gojue/moling/actions/workflows/go-test.yml/badge.svg)](https://github.com/gojue/ecapture/actions/workflows/go-test.yml)
[![Github Version](https://img.shields.io/github/v/release/gojue/moling?display_name=tag&include_prereleases&sort=semver)](https://github.com/gojue/moling/releases)

---

### Introduction
MoLing is a computer-based MCP Server that implements system interaction through operating system APIs, enabling file system operations such as reading, writing, merging, statistics, and aggregation, as well as the ability to execute system commands. It is a dependency-free local office automation assistant.

### Advantages
> [!IMPORTANT]
> Requiring no installation of any dependencies, MoLing can be run directly and is compatible with multiple operating systems, including Windows, Linux, and macOS. 
> This eliminates the hassle of dealing with environment conflicts involving Node.js, Python, and other development environments.

### Features

- **File System Operations**: Reading, writing, merging, statistics, and aggregation
- **Command-line Terminal**: Execute system commands directly
- **Browser Control**: Powered by `github.com/chromedp/chromedp`
- **Future Plans**:
    - Personal PC data organization
    - Document writing assistance
    - Schedule planning
    - Life assistant features

### Supported MCP Clients

- Claude
- Cline
- Cherry Studio

#### Configuration Format
For example, to configure the Claude client, add the following configuration:

> [!TIP]
> Only 3-6 lines of configuration are needed.

```json
{
  "mcpServers": {
    "MoLing": {
      "command": "/usr/local/bin/moling",
      "args": []
    }
  }
}
```

### Operation Modes

- **Stdio Mode**: CLI-based interactive mode for user-friendly experience
- **SSR Mode**: Server-Side Rendering mode optimized for headless/automated environments

### Installation

#### Option 1: Install via Script
##### Linux/MacOS
```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.sh)"
```
##### Windows
```powershell
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.ps1 | iex"
```

#### Option 2: Direct Download
1. Download the installation package from [releases page](https://github.com/gojue/moling/releases)
2. Extract the package
3. Run the server:
   ```sh
   ./moling
   ```

#### Option 3: Build from Source
1. Clone the repository:
```sh
git clone https://github.com/gojue/moling.git
cd moling
```
2. Build the project (requires Golang toolchain):
```sh
make build
```
3. Run the compiled binary:
```sh
./bin/moling
```

### Usage
After starting the server, connect using any supported MCP client by configuring it to point to your MoLing server address.

### License
Apache License 2.0. See [LICENSE](LICENSE) for details.