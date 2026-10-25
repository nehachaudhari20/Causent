"""Benchmark metrics for localization, linking, and abstention quality."""

from __future__ import annotations


def unsupported_definitive_rate(_predictions: list[dict], _labels: list[dict]) -> float:
    raise NotImplementedError
