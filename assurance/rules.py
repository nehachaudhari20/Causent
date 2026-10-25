"""Versioned deterministic assurance rules."""

from __future__ import annotations

from enum import Enum


class RuleStatus(str, Enum):
    SATISFIED = "SATISFIED"
    VIOLATED = "VIOLATED"
    UNKNOWN = "UNKNOWN"
    CONFLICT = "CONFLICT"
