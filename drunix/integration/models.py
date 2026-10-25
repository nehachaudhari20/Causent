"""Request/response models for Drunix chaincode integration."""

from __future__ import annotations

from pydantic import BaseModel, Field


class OpenCaseRequest(BaseModel):
    case_id: str
    policy_version: str = "1.0"


class SubmitAttestationRequest(BaseModel):
    case_id: str
    evidence_id: str
    event_digest: str
    issuer_organization_id: str
    event_type: str
    native_event_id: str = Field(..., description="Idempotency key within issuer/source namespace")


class RegisterEvidenceBatchRequest(BaseModel):
    case_id: str
    batch_id: str
    batch_digest: str
    merkle_root: str | None = None


class RecordAssessmentRequest(BaseModel):
    case_id: str
    assessment_id: str
    assessment_digest: str
    rule_set_version: str
    assurance_outcome: str


class RaiseChallengeRequest(BaseModel):
    case_id: str
    challenge_id: str
    target_type: str
    target_id: str
    reason_code: str


class AppendCorrectionRequest(BaseModel):
    case_id: str
    correction_id: str
    prior_record_id: str
    correction_digest: str


class SealCaseRequest(BaseModel):
    case_id: str
    seal_digest: str


class CaseRecord(BaseModel):
    case_id: str
    state: str
    policy_version: str
    owner_msp: str
    created_at: str
    updated_at: str


class LedgerSubmitResult(BaseModel):
    success: bool
    mode: str
    transaction_id: str | None = None
    payload: dict | None = None
    message: str | None = None
