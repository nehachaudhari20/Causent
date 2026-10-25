"""Causent HTTP API — evidence intake, assurance, and Drunix governance hooks."""

from __future__ import annotations

import uuid
from typing import Any

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

from drunix.integration.client import DrunixClient, DrunixSettings
from drunix.integration.models import (
    OpenCaseRequest,
    RecordAssessmentRequest,
    SubmitAttestationRequest,
)

app = FastAPI(
    title="Causent",
    description="Beyond the Receipt — Execution Assurance for Agentic Payments",
    version="0.1.0",
)

_drunix = DrunixClient(DrunixSettings())


class CreateCaseBody(BaseModel):
    policy_version: str = "1.0"
    case_id: str | None = None


class SubmitEvidenceBody(BaseModel):
    evidence_id: str | None = None
    event_digest: str
    issuer_organization_id: str
    event_type: str
    native_event_id: str


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok", "service": "causent-api"}


@app.post("/v1/cases")
def create_case(body: CreateCaseBody) -> dict[str, Any]:
    case_id = body.case_id or f"case_{uuid.uuid4().hex[:12]}"
    result = _drunix.open_case(OpenCaseRequest(case_id=case_id, policy_version=body.policy_version))
    if not result.success:
        raise HTTPException(status_code=502, detail=result.message or "ledger submit failed")
    return {"case_id": case_id, "ledger": result.model_dump()}


@app.post("/v1/cases/{case_id}/evidence")
def submit_evidence(case_id: str, body: SubmitEvidenceBody) -> dict[str, Any]:
    evidence_id = body.evidence_id or f"evd_{uuid.uuid4().hex[:12]}"
    result = _drunix.submit_attestation(
        SubmitAttestationRequest(
            case_id=case_id,
            evidence_id=evidence_id,
            event_digest=body.event_digest,
            issuer_organization_id=body.issuer_organization_id,
            event_type=body.event_type,
            native_event_id=body.native_event_id,
        )
    )
    if not result.success:
        raise HTTPException(status_code=502, detail=result.message or "ledger submit failed")
    return {"evidence_id": evidence_id, "ledger": result.model_dump()}


@app.get("/v1/cases/{case_id}")
def get_case(case_id: str) -> dict[str, Any]:
    record = _drunix.get_case(case_id)
    if record is None:
        return {"case_id": case_id, "ledger_state": None, "note": "stub mode or case not on ledger"}
    return record.model_dump()


class RecordAssessmentBody(BaseModel):
    assessment_id: str = Field(default_factory=lambda: f"asm_{uuid.uuid4().hex[:12]}")
    assessment_digest: str
    rule_set_version: str = "1.0"
    assurance_outcome: str


@app.post("/v1/cases/{case_id}/assessments")
def record_assessment(case_id: str, body: RecordAssessmentBody) -> dict[str, Any]:
    result = _drunix.record_assessment(
        RecordAssessmentRequest(case_id=case_id, **body.model_dump())
    )
    if not result.success:
        raise HTTPException(status_code=502, detail=result.message or "ledger submit failed")
    return {"assessment_id": body.assessment_id, "ledger": result.model_dump()}


@app.get("/v1/cases/{case_id}/graph")
def get_graph(case_id: str) -> dict[str, Any]:
    return {"case_id": case_id, "nodes": [], "edges": [], "status": "not_implemented"}


@app.get("/v1/cases/{case_id}/report")
def get_report(case_id: str) -> dict[str, Any]:
    return {"case_id": case_id, "status": "not_implemented"}


def run() -> None:
    import uvicorn

    uvicorn.run("api.main:app", host="0.0.0.0", port=8080, reload=True)


if __name__ == "__main__":
    run()
