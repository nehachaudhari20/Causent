"""Sandbox payment protocol adapter (signed test fixtures only)."""

from __future__ import annotations

from typing import Any


def normalize_native_event(native: dict[str, Any]) -> dict[str, Any]:
    """Map a sandbox payment native event into a canonical event envelope."""
    raise NotImplementedError("Sandbox payment adapter not implemented yet")
