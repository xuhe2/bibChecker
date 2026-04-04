import os
from dataclasses import dataclass
from pathlib import Path
from dotenv import load_dotenv

# Load .env file
load_dotenv()


@dataclass
class Config:
    """Application configuration"""
    http_proxy: str | None
    https_proxy: str | None
    timeout: int = 10
    delay: float = 1.0  # 请求间隔(秒)

    @classmethod
    def from_env(cls) -> "Config":
        """从环境变量创建配置"""
        return cls(
            http_proxy=os.getenv("HTTP_PROXY") or os.getenv("http_proxy"),
            https_proxy=os.getenv("HTTPS_PROXY") or os.getenv("https_proxy"),
            timeout=int(os.getenv("TIMEOUT", "10")),
            delay=float(os.getenv("DELAY", "1.0")),
        )

    def get_proxy_dict(self) -> dict | None:
        """获取代理字典"""
        if self.https_proxy:
            return {"https": self.https_proxy}
        if self.http_proxy:
            return {"http": self.http_proxy}
        return None


def ensure_output_dir(output_path: str) -> Path:
    """确保输出目录存在"""
    path = Path(output_path)
    path.parent.mkdir(parents=True, exist_ok=True)
    return path