# Drunix integration

Causent uses Drunix as the **governance plane**: case lifecycle, evidence digests, assessment commitments, challenges, and seals. Assurance logic (EBDL, MLPG) stays in the Python application layer.

## Chaincode

| Path | Description |
|------|-------------|
| `chaincode/causent-assurance/` | Go chaincode (`causent-assurance`) |

### Operations

- `OpenCase(caseID, policyVersion)`
- `SubmitAttestation(caseID, evidenceID, eventDigest, issuerOrg, eventType, nativeEventID)`
- `RegisterEvidenceBatch(caseID, batchID, batchDigest, merkleRoot)`
- `RecordAssessment(caseID, assessmentID, assessmentDigest, ruleSetVersion, assuranceOutcome)`
- `RaiseChallenge(caseID, challengeID, targetType, targetID, reasonCode)`
- `AppendCorrection(caseID, correctionID, priorRecordID, correctionDigest)`
- `SealCase(caseID, sealDigest)`
- `GetCase(caseID)` — query
- `VerifyCommitment(caseID, batchID, digest)` — query

## Application client

`integration/client.py` exposes `DrunixClient`:

- Default: **stub mode** (`DRUNIX_GATEWAY_ENABLED=false`) for local API development.
- Live network: set env vars from `configs/example.env`, deploy chaincode on your Drunix test network, then enable the gateway and use the Fabric `peer` CLI (or replace with Fabric Gateway SDK once connection profiles are wired).

## Deploy (outline)

1. Start the [Drunix](https://github.com/drunix/drunix) test network (YugabyteDB state DB per your network config).
2. Package and install `causent-assurance` on participating peers (`AgentPlatformMSP`, `ServiceProviderMSP`, `PaymentParticipantMSP`).
3. Approve and commit chaincode on `causent-channel` (name may match your network script).
4. Copy connection profile / TLS CAs into `configs/` (gitignored) and set `DRUNIX_GATEWAY_ENABLED=true`.

See `scripts/deploy-chaincode.sh` for a template invoke sequence—adjust paths to your Drunix checkout.
