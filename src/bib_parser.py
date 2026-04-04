import re
from dataclasses import dataclass


@dataclass
class BibEntry:
    """BibTeX 条目"""
    cite_key: str
    title: str
    entry_type: str = ""


def parse_bib_file(bib_path: str) -> list[BibEntry]:
    """解析 .bib 文件，返回 BibEntry 列表"""
    entries = []
    with open(bib_path, "r", encoding="utf-8") as f:
        content = f.read()

    # 匹配 @article{key, ... title = {...}, ... }
    entry_pattern = r'@(\w+)\s*\{\s*([^,\s]+)\s*,[^}]*?title\s*=\s*\{([^}]+)\}'
    for match in re.finditer(entry_pattern, content, re.DOTALL | re.IGNORECASE):
        entry_type = match.group(1)
        cite_key = match.group(2).strip()
        title = match.group(3).strip()
        entries.append(BibEntry(cite_key=cite_key, title=title, entry_type=entry_type))

    return entries


def filter_entries(entries: list[BibEntry], keywords: list[str]) -> list[BibEntry]:
    """根据关键词过滤条目"""
    if not keywords:
        return entries

    filtered = []
    for entry in entries:
        if any(kw.lower() in entry.title.lower() for kw in keywords):
            filtered.append(entry)
    return filtered