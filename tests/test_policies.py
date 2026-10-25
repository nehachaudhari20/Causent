"""Chaincode policy tests run against a live Drunix network (skipped by default)."""

import pytest

pytestmark = pytest.mark.skip(reason="Requires deployed causent-assurance on Drunix test network")
