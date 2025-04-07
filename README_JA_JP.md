## MoLing MCP サーバー

[English](./README.md) | [中文](./README_ZH_HANS.md) | 日本語

[![GitHub stars](https://img.shields.io/github/stars/gojue/moling.svg?label=Stars&logo=github)](https://github.com/gojue/moling/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/gojue/moling?label=Forks&logo=github)](https://github.com/gojue/moling/forks)
[![CI](https://github.com/gojue/moling/actions/workflows/go-test.yml/badge.svg)](https://github.com/gojue/moling/actions/workflows/go-test.yml)
[![Github Version](https://img.shields.io/github/v/release/gojue/moling?display_name=tag&include_prereleases&sort=semver)](https://github.com/gojue/moling/releases)

---

![](./images/logo.svg)

### 紹介
MoLingは、オペレーティングシステムAPIを介してシステム操作を実装するコンピュータ使用およびブラウザ使用のMCPサーバーであり、ファイルシステム操作（読み取り、書き込み、マージ、統計、集計）やシステムコマンドの実行を可能にします。依存関係のないローカルオフィス自動化アシスタントです。

### 利点
> [!IMPORTANT]
> 依存関係のインストールを必要とせず、MoLingは直接実行でき、Windows、Linux、macOSなどの複数のオペレーティングシステムと互換性があります。
> これにより、Node.js、Python、Dockerなどの開発環境に関連する環境の競合を処理する手間が省けます。

### 機能

> [!CAUTION]
> コマンドライン操作は危険であり、慎重に使用する必要があります。

- **ファイルシステム操作**：読み取り、書き込み、マージ、統計、集計
- **コマンドラインターミナル**：システムコマンドを直接実行
- **ブラウザ制御**：`github.com/chromedp/chromedp`によって提供される
- **将来の計画**：
    - 個人PCデータの整理
    - ドキュメント作成支援
    - スケジュール計画
    - 生活アシスタント機能

> [!WARNING]
> 現在、MoLingはmacOSでのみテストされており、他のオペレーティングシステムでは問題が発生する可能性があります。

### サポートされているMCPクライアント

- [Claude](https://claude.ai/)
- [Cline](https://cline.bot/)
- [Cherry Studio](https://cherry-ai.com/)
- その他（MCPプロトコルをサポートするクライアント）

#### スクリーンショット

https://github.com/user-attachments/assets/229c4dd5-23b4-4b53-9e25-3eba8734b5b7

[Claude](https://claude.ai/)に統合されたMoLing
![](./images/screenshot_claude.png)

#### 設定形式

##### MCPサーバー（MoLing）設定

設定ファイルは`/Users/username/.moling/config/config.json`に生成され、必要に応じて内容を変更できます。

ファイルが存在しない場合は、`moling config --init`を使用して作成できます。

##### MCPクライアント設定
例として、Claudeクライアントを設定するには、次の設定を追加します：

> [!TIP]
> 
> 3〜6行の設定のみが必要です。
> 
> Claude設定パス：`~/Library/Application\ Support/Claude/claude_desktop_config`

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

また、`/usr/local/bin/moling`はダウンロードしたMoLingサーバーバイナリのパスです。

**自動設定**

`moling client --install`を実行して、MCPクライアントの設定を自動的にインストールします。

MoLingはMCPクライアントを自動的に検出し、設定をインストールします。Cline、Claude、Roo Codeなどを含みます。

### 動作モード

- **Stdioモード**：CLIベースのインタラクティブモードで、ユーザーフレンドリーな体験を提供
- **SSRモード**：ヘッドレス/自動化環境に最適化されたサーバーサイドレンダリングモード

### インストール

#### オプション1：スクリプトを使用してインストール
##### Linux/MacOS
```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.sh)"
```
##### Windows

> [!WARNING]
> テストされていないため、動作するかどうかは不明です。

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/gojue/moling/HEAD/install/install.ps1 | iex"
```

#### オプション2：直接ダウンロード
1. [リリースページ](https://github.com/gojue/moling/releases)からインストールパッケージをダウンロード
2. パッケージを解凍
3. サーバーを実行：
   ```sh
   ./moling
   ```

#### オプション3：ソースからビルド
1. リポジトリをクローン：
```sh
git clone https://github.com/gojue/moling.git
cd moling
```
2. プロジェクトをビルド（Golangツールチェーンが必要）：
```sh
make build
```
3. コンパイルされたバイナリを実行：
```sh
./bin/moling
```

### 使用方法
サーバーを起動した後、サポートされているMCPクライアントを使用して、MoLingサーバーアドレスに接続します。

### ライセンス
Apache License 2.0。詳細は[LICENSE](LICENSE)ファイルを参照してください。
