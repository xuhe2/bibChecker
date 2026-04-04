import logging
import time
import os
from scholarly import scholarly
from src.config import Config

logger = logging.getLogger(__name__)


class BibTeXFetcher:
    """BibTeX 获取器"""

    def __init__(self, config: Config):
        self.config = config
        self._setup_proxy()

    def _setup_proxy(self):
        """设置代理"""
        proxy = self.config.get_proxy_dict()
        if proxy:
            logger.info(f"使用代理: {proxy}")
            os.environ["HTTP_PROXY"] = self.config.http_proxy or ""
            os.environ["HTTPS_PROXY"] = self.config.https_proxy or ""

    def fetch(self, title: str) -> str | None:
        """通过标题获取 BibTeX"""
        try:
            search_result = next(scholarly.search_pubs(title))
            time.sleep(self.config.delay)
            return scholarly.bibtex(search_result)
        except StopIteration:
            logger.warning(f"未找到论文: {title}")
            return None
        except Exception as e:
            logger.error(f"搜索出错 [{title}]: {e}")
            return None

    def fetch_batch(self, titles: list[str]) -> list[tuple[str, str | None]]:
        """批量获取 BibTeX"""
        results = []
        for i, title in enumerate(titles, 1):
            logger.info(f"[{i}/{len(titles)}] 正在搜索: {title}")
            bibtex = self.fetch(title)
            results.append((title, bibtex))
        return results