"""Drunix / Fabric gateway client for causent-assurance chaincode."""

from __future__ import annotations

import json
import logging
import subprocess
from typing import Any

from pydantic_settings import BaseSettings, SettingsConfigDict

from drunix.integration.models import (
    AppendCorrectionRequest,
    CaseRecord,
    LedgerSubmitResult,
    OpenCaseRequest,
    RaiseChallengeRequest,
    RecordAssessmentRequest,
    RegisterEvidenceBatchRequest,
    SealCaseRequest,
    SubmitAttestationRequest,
)

logger = logging.getLogger(__name__)


class DrunixSettings(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="DRUNIX_", env_file=".env", extra="ignore")

    channel_name: str = "causent-channel"
    chaincode_name: str = "causent-assurance"
    msp_id: str = "AgentPlatformMSP"
    peer_endpoint: str = "localhost:7051"
    tls_enabled: bool = True
    gateway_enabled: bool = False
    peer_cli_path: str = "peer"


class DrunixClient:
    """
    Thin wrapper around chaincode invoke/query.

    When ``gateway_enabled`` is false (default), invokes are recorded locally and
    returned as stub results so the API and tests can run without a live network.
    Set ``DRUNIX_GATEWAY_ENABLED=true`` and configure the peer CLI / connection
    profile after bringing up the Drunix test network.
    """

    def __init__(self, settings: DrunixSettings | None = None) -> None:
        self.settings = settings or DrunixSettings()

    def open_case(self, req: OpenCaseRequest) -> LedgerSubmitResult:
        return self._invoke("OpenCase", req.case_id, req.policy_version)

    def submit_attestation(self, req: SubmitAttestationRequest) -> LedgerSubmitResult:
        return self._invoke(
            "SubmitAttestation",
            req.case_id,
            req.evidence_id,
            req.event_digest,
            req.issuer_organization_id,
            req.event_type,
            req.native_event_id,
        )

    def register_evidence_batch(self, req: RegisterEvidenceBatchRequest) -> LedgerSubmitResult:
        root = req.merkle_root or ""
        return self._invoke(
            "RegisterEvidenceBatch",
            req.case_id,
            req.batch_id,
            req.batch_digest,
            root,
        )

    def record_assessment(self, req: RecordAssessmentRequest) -> LedgerSubmitResult:
        return self._invoke(
            "RecordAssessment",
            req.case_id,
            req.assessment_id,
            req.assessment_digest,
            req.rule_set_version,
            req.assurance_outcome,
        )

    def raise_challenge(self, req: RaiseChallengeRequest) -> LedgerSubmitResult:
        return self._invoke(
            "RaiseChallenge",
            req.case_id,
            req.challenge_id,
            req.target_type,
            req.target_id,
            req.reason_code,
        )

    def append_correction(self, req: AppendCorrectionRequest) -> LedgerSubmitResult:
        return self._invoke(
            "AppendCorrection",
            req.case_id,
            req.correction_id,
            req.prior_record_id,
            req.correction_digest,
        )

    def seal_case(self, req: SealCaseRequest) -> LedgerSubmitResult:
        return self._invoke("SealCase", req.case_id, req.seal_digest)

    def get_case(self, case_id: str) -> CaseRecord | None:
        raw = self._query("GetCase", case_id)
        if raw is None:
            return None
        return CaseRecord.model_validate(raw)

    def verify_commitment(self, case_id: str, batch_id: str, digest: str) -> dict[str, Any]:
        raw = self._query("VerifyCommitment", case_id, batch_id, digest)
        return raw or {"verified": False}

    def _invoke(self, fn: str, *args: str) -> LedgerSubmitResult:
        payload = {"function": fn, "args": list(args)}
        if not self.settings.gateway_enabled:
            logger.info("Drunix stub invoke: %s", payload)
            return LedgerSubmitResult(
                success=True,
                mode="stub",
                transaction_id=f"stub-{fn}-{args[0] if args else 'none'}",
                payload=payload,
                message="Enable DRUNIX_GATEWAY_ENABLED to submit to a live network.",
            )
        return self._peer_chaincode_invoke(fn, list(args))

    def _query(self, fn: str, *args: str) -> dict[str, Any] | None:
        if not self.settings.gateway_enabled:
            return None
        return self._peer_chaincode_query(fn, list(args))

    def _peer_chaincode_invoke(self, fn: str, args: list[str]) -> LedgerSubmitResult:
        arg_json = json.dumps({"Args": [fn, *args]})
        cmd = [
            self.settings.peer_cli_path,
            "chaincode",
            "invoke",
            "-C",
            self.settings.channel_name,
            "-n",
            self.settings.chaincode_name,
            "-c",
            arg_json,
        ]
        try:
            completed = subprocess.run(cmd, capture_output=True, text=True, check=False)
        except FileNotFoundError:
            return LedgerSubmitResult(
                success=False,
                mode="peer_cli",
                message="peer CLI not found; install Fabric peer binary or disable gateway.",
            )
        if completed.returncode != 0:
            return LedgerSubmitResult(
                success=False,
                mode="peer_cli",
                message=completed.stderr or completed.stdout,
            )
        return LedgerSubmitResult(success=True, mode="peer_cli", message=completed.stdout)

    def _peer_chaincode_query(self, fn: str, args: list[str]) -> dict[str, Any] | None:
        arg_json = json.dumps({"Args": [fn, *args]})
        cmd = [
            self.settings.peer_cli_path,
            "chaincode",
            "query",
            "-C",
            self.settings.channel_name,
            "-n",
            self.settings.chaincode_name,
            "-c",
            arg_json,
        ]
        try:
            completed = subprocess.run(cmd, capture_output=True, text=True, check=False)
        except FileNotFoundError:
            return None
        if completed.returncode != 0:
            return None
        try:
            return json.loads(completed.stdout)
        except json.JSONDecodeError:
            return {"raw": completed.stdout}
