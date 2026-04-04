# BibChecker

从 `.bib` 文件解析论文标题，并从 Google Scholar 获取更新后的 BibTeX。

## 安装

```bash
uv pip install -e .
```

## 配置

复制 `.env.example` 为 `.env`，根据需要修改代理配置：

```bash
cp .env.example .env
```

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `HTTP_PROXY` | HTTP 代理地址 | - |
| `HTTPS_PROXY` | HTTPS 代理地址 | - |
| `TIMEOUT` | 请求超时时间(秒) | 10 |
| `DELAY` | 请求间隔(秒) | 1.0 |

## 使用

```bash
# 处理整个 bib 文件
uv run main.py input.bib -o output.bib

# 只处理标题包含指定关键词的论文
uv run main.py input.bib -t "attention" -t "transformer"

# 指定输出文件
uv run main.py input.bib -o my_papers.bib
```

### 参数说明

- `bib_file`: 输入的 `.bib` 文件路径
- `-o, --output`: 输出文件路径 (默认: `output.bib`)
- `-t, --title`: 标题过滤关键词，可多次使用

## 项目结构

```
bibChecker/
├── main.py              # 入口文件
├── src/
│   ├── __init__.py
│   ├── bib_parser.py    # Bib 文件解析
│   ├── config.py        # 配置管理
│   └── fetcher.py       # Google Scholar 获取
├── .env.example         # 环境变量示例
└── pyproject.toml
```

## 依赖

- [scholarly](https://github.com/scholarly-python-package/scholarly) - Google Scholar API
- [python-dotenv](https://github.com/theskumar/python-dotenv) - 环境变量管理