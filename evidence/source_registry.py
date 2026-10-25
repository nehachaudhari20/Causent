"""Source key registry: which keys may attest which claim types."""

from __future__ import annotations

from pathlib import Path
from typing import Any


def load_registry(path: Path) -> dict[str, Any]:
    raise NotImplementedError
