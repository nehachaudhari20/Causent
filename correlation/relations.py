"""Evidence-qualified event relations for the MLPG."""

from __future__ import annotations

from enum import Enum


class RelationBasis(str, Enum):
    NATIVE_PROTOCOL_REFERENCE = "NATIVE_PROTOCOL_REFERENCE"
    AUTHENTICATED_REFERENCE = "AUTHENTICATED_REFERENCE"
    EXECUTION_CONTEXT = "EXECUTION_CONTEXT"
    CANDIDATE_MATCH = "CANDIDATE_MATCH"
