# gitgram

GitHub Webhook to Telegram notification service written in Go.

## 功能

- 接收 GitHub repository Webhook。
- 驗證 `X-Hub-Signature-256`。
- 支援重要通知：`pull_request`、`issues`、`pull_request_review`、`release`。
- 支援 `workflow_run`，但只通知 `conclusion=failure` 的 GitHub Actions 失敗。
- 使用 Telegram Bot API `sendMessage` 發送 HTML 格式訊息。
- 使用 `X-GitHub-Delivery` 做 in-memory 去重。
- 透過 `ALLOWED_REPOS` 限制允許的 repository。

## 環境變數

可以先從範例建立本機設定：

```sh
cp .env.example .env
```

目前服務會讀取 process environment，不會自動載入 `.env` 檔。若用 `.env`，請透過 shell、部署平台或 dotenv 工具載入。

```text
PORT=8080
GITHUB_WEBHOOK_SECRET=your-github-webhook-secret
TELEGRAM_BOT_TOKEN=123456:telegram-bot-token
TELEGRAM_CHAT_ID=-1001234567890
ALLOWED_REPOS=owner/repo-a,owner/repo-b
```

可選：

```text
MAX_BODY_BYTES=10485760
DELIVERY_CACHE_SIZE=10000
PUBLIC_BASE_URL=https://your-domain.example
```

`ALLOWED_REPOS` 支援：

- `owner/repo`
- `owner/*`
- `*`

## 執行

```sh
go run ./cmd/gitgram
```

## Docker

Build image：

```sh
docker build -t gitgram .
```

Run container：

```sh
docker run --rm -p 8080:8080 --env-file .env gitgram
```

健康檢查：

```text
GET /healthz
```

GitHub Webhook：

```text
POST /webhooks/github
```

## GitHub Webhook 設定

- Payload URL: `https://your-domain.example/webhooks/github`
- Content type: `application/json`
- Secret: 和 `GITHUB_WEBHOOK_SECRET` 相同
- Events: 建議選擇 `Pull requests`、`Issues`、`Pull request reviews`、`Releases`、`Workflow runs`。即使先用 `Send me everything`，服務也只會通知重要事件。

## Telegram 設定

1. 用 BotFather 建立 bot，取得 `TELEGRAM_BOT_TOKEN`。
2. 把 bot 加入 Telegram group 或 channel。
3. 取得目標 chat id，設定到 `TELEGRAM_CHAT_ID`。

## 訊息格式

服務會使用 Telegram `HTML` parse mode 發送訊息。Repo 名稱會是粗體，branch 和 actor 會是 monospace，GitHub URL 會顯示成 `Open on GitHub` link。

Pull Request：

```text
owner/repo pull request opened
#12 Add login flow
by octocat
feature/login -> main
Open on GitHub
```

`pull_request` 只會通知 `opened`、`reopened`、`ready_for_review`、`review_requested`，以及已 merge 的 `closed`。

Issue：

```text
owner/repo issue opened
#42 Cannot login with OAuth
by octocat
Open on GitHub
```

`issues` 只會通知 `opened` 和 `reopened`。Issue/PR comments 會被忽略。

Pull Request review requested：

```text
owner/repo pull request review requested
#12 Add login flow
by octocat
feature/login -> main

review requested from mona

Open on GitHub
```

Pull Request review：

```text
owner/repo pull request review changes requested
#12 Add login flow
by octocat

Please handle the empty token case before merge.

Open on GitHub
```

`pull_request_review` 只會通知 `approved` 和 `changes_requested`。一般 review comment 會被忽略。

Release：

```text
owner/repo release published
v1.0.0
by octocat
Open on GitHub
```

GitHub Actions workflow failure：

```text
owner/repo workflow failed
CI on main
by octocat
Open on GitHub
```

`workflow_run` 只會在 `action=completed` 且 `conclusion=failure` 時發送。成功、取消、略過的 workflow run 會被忽略。

## 測試

```sh
go test ./...
```
