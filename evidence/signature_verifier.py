"""Cryptographic verification of participant evidence."""

from __future__ import annotations

from typing import Any


def verify_signature(payload: bytes, signature: bytes, key_material: Any) -> bool:
    raise NotImplementedError
