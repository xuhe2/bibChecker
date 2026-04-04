import argparse
import logging
import sys

from src.config import Config, ensure_output_dir
from src.bib_parser import parse_bib_file, filter_entries, BibEntry
from src.fetcher import BibTeXFetcher

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
    datefmt="%H:%M:%S",
)
logger = logging.getLogger(__name__)


def save_bibtex(bibtex_list: list[str], output_path: str):
    """保存 BibTeX 到文件"""
    ensure_output_dir(output_path)
    with open(output_path, "w", encoding="utf-8") as f:
        for bibtex in bibtex_list:
            f.write(bibtex)
            f.write("\n\n")


def main():
    parser = argparse.ArgumentParser(
        description="从 .bib 文件解析论文标题并从 Google Scholar 获取 BibTeX"
    )
    parser.add_argument("bib_file", help="输入的 .bib 文件路径")
    parser.add_argument("-o", "--output", default="output.bib", help="输出的 .bib 文件路径 (默认: output.bib)")
    parser.add_argument("-t", "--title", action="append", help="只处理标题包含指定关键词的论文 (可多次使用)")

    args = parser.parse_args()

    # 加载配置
    config = Config.from_env()

    # 解析 bib 文件
    entries = parse_bib_file(args.bib_file)
    logger.info(f"解析到 {len(entries)} 条目")

    # 过滤条目
    if args.title:
        entries = filter_entries(entries, args.title)
        logger.info(f"关键词过滤后剩余 {len(entries)} 条")

    if not entries:
        logger.warning("没有需要处理的条目")
        sys.exit(0)

    # 获取 BibTeX
    fetcher = BibTeXFetcher(config)
    bibtex_list = []

    for i, entry in enumerate(entries, 1):
        logger.info(f"[{i}/{len(entries)}] 正在搜索: {entry.title}")
        bibtex = fetcher.fetch(entry.title)
        if bibtex:
            bibtex_list.append(bibtex)
            logger.info(f"  获取成功")
        else:
            logger.info(f"  获取失败")

    # 保存结果
    if bibtex_list:
        save_bibtex(bibtex_list, args.output)
        logger.info(f"已保存 {len(bibtex_list)} 条 BibTeX 到 {args.output}")
    else:
        logger.warning("没有获取到任何 BibTeX")


if __name__ == "__main__":
    main()