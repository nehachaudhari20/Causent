"""Idempotency and replay protection for ingested events."""

from __future__ import annotations


class ReplayGuard:
    def __init__(self) -> None:
        self._seen: set[str] = set()

    def check_and_record(self, idempotency_key: str) -> bool:
        """Return True if new; False if duplicate."""
        if idempotency_key in self._seen:
            return False
        self._seen.add(idempotency_key)
        return True
