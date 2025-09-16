# Redmine CLI

RedmineのIssueやProjectを管理するためのコマンドラインツールです。

## 機能

- **Issue管理**: Redmine Issueの一覧表示・詳細表示
- **認証管理**: APIトークンの管理
- **プロファイル管理**: 複数のRedmineサーバー接続情報の管理

## インストール

### 前提条件

- Go 1.24.3 以上

### ビルド

```bash
go mod tidy
go build -o redmine
```

## 設定

### プロファイルの追加

最初にRedmineサーバーの接続情報を設定します:

```bash
./redmine profile add <profile_name> <redmine_url> <api_token>
```

例:

```bash
./redmine profile add production https://redmine.example.com abcd1234567890
```

### プロファイル管理

```bash
# プロファイル一覧
./redmine profile list

# デフォルトプロファイルの設定
./redmine profile use <profile_name>

# プロファイル詳細表示
./redmine profile show [profile_name]

# プロファイル削除
./redmine profile remove <profile_name>
```

## 使い方

### Issue管理

#### Issue一覧の表示

```bash
./redmine issues list
```

オプション:

- `--limit`: 取得件数 (デフォルト: 25)
- `--offset`: オフセット (デフォルト: 0)
- `--project`: プロジェクトIDでフィルタ
- `--status`: ステータスIDでフィルタ

例:

```bash
./redmine issues list --limit 50 --project 1 --status 1
```

#### Issue詳細の表示

```bash
./redmine issues show <issue_id>
```

オプション:

- `--comments`, `-c`: コメント（journal）を含めて表示

例:

```bash
./redmine issues show 123 --comments
```

### 認証管理（非推奨）

```bash
# APIトークンの設定（profile addの使用を推奨）
./redmine auth token add <token>

# Redmine URLの設定（profile addの使用を推奨）
./redmine config set-url <url>

# 設定の表示（profile showの使用を推奨）
./redmine config show
```

### プロファイルの使用

特定のプロファイルを一時的に使用する場合:

```bash
./redmine --profile <profile_name> issues list
```

## 依存関係

- [spf13/cobra](https://github.com/spf13/cobra): CLIフレームワーク
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3): YAML設定ファイル処理

## ライセンス

MIT - 詳細は [LICENSE.md](LICENSE.md) を参照してください。
