import argparse
import logging
import sys

from src.config import Config, ensure_output_dir
from src.bib_parser import parse_bib_file, filter_entries
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


def cmd_search(args):
    """search 子命令：按关键词搜索 Google Scholar"""
    config = Config.from_env()
    fetcher = BibTeXFetcher(config)

    logger.info(f"搜索关键词: {args.query}")
    bibtex_list = fetcher.search_by_query(args.query, args.num)
    logger.info(f"获取到 {len(bibtex_list)} 条 BibTeX")

    if bibtex_list:
        if args.output:
            save_bibtex(bibtex_list, args.output)
            logger.info(f"已保存到 {args.output}")
        else:
            for bibtex in bibtex_list:
                print(bibtex)
                print()
    else:
        logger.warning("没有获取到任何 BibTeX")


def cmd_update(args):
    """update 子命令：解析 .bib 文件并更新 BibTeX"""
    config = Config.from_env()
    fetcher = BibTeXFetcher(config)

    entries = parse_bib_file(args.bib_file)
    logger.info(f"解析到 {len(entries)} 条目")

    if args.title:
        entries = filter_entries(entries, args.title)
        logger.info(f"标题过滤后剩余 {len(entries)} 条")

    if not entries:
        logger.warning("没有需要处理的条目")
        sys.exit(0)

    bibtex_list = []
    for i, entry in enumerate(entries, 1):
        logger.info(f"[{i}/{len(entries)}] 正在搜索: {entry.title}")
        bibtex = fetcher.fetch(entry.title)
        if bibtex:
            bibtex_list.append(bibtex)
            logger.info(f"  获取成功")
        else:
            logger.info(f"  获取失败")

    if bibtex_list:
        save_bibtex(bibtex_list, args.output)
        logger.info(f"已保存 {len(bibtex_list)} 条 BibTeX 到 {args.output}")
    else:
        logger.warning("没有获取到任何 BibTeX")


def main():
    parser = argparse.ArgumentParser(
        description="bibChecker - 从 Google Scholar 获取和更新 BibTeX"
    )
    subparsers = parser.add_subparsers(dest="command", help="可用命令")

    # search 子命令
    search_parser = subparsers.add_parser("search", help="按关键词搜索 Google Scholar")
    search_parser.add_argument("query", help="搜索关键词")
    search_parser.add_argument("-o", "--output",
                               help="输出文件路径 (不指定则输出到终端)")
    search_parser.add_argument("-n", "--num", type=int, default=5,
                               help="返回结果数量 (默认: 5)")

    # update 子命令
    update_parser = subparsers.add_parser("update", help="更新 .bib 文件中的条目")
    update_parser.add_argument("bib_file", help="输入的 .bib 文件路径")
    update_parser.add_argument("-o", "--output", default="output.bib",
                               help="输出文件路径 (默认: output.bib)")
    update_parser.add_argument("-t", "--title", action="append",
                               help="只更新标题包含指定关键词的条目 (可多次使用)")

    args = parser.parse_args()

    if args.command == "search":
        cmd_search(args)
    elif args.command == "update":
        cmd_update(args)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
