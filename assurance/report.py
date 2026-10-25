"""Canonical assessment report serialization and digests."""

from __future__ import annotations

import hashlib
import json
from typing import Any


def report_digest(report: dict[str, Any]) -> str:
    canonical = json.dumps(report, sort_keys=True, separators=(",", ":"))
    digest = hashlib.sha256(canonical.encode("utf-8")).hexdigest()
    return f"sha256:{digest}"
