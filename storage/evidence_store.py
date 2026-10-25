"""Encrypted off-chain evidence object storage."""

from __future__ import annotations

from pathlib import Path


class EvidenceStore:
    def __init__(self, base_path: Path) -> None:
        self.base_path = base_path
        self.base_path.mkdir(parents=True, exist_ok=True)

    def put(self, object_id: str, payload: bytes) -> str:
        path = self.base_path / f"{object_id}.bin"
        path.write_bytes(payload)
        return f"file://{path.as_posix()}"
