package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	caseStateOpen              = "OPEN"
	caseStateCollecting        = "COLLECTING_EVIDENCE"
	caseStateAssessed          = "ASSESSED"
	caseStateDisputed          = "DISPUTED"
	caseStateResolved          = "RESOLVED"
	caseStateSealed            = "SEALED"
)

// AssuranceContract implements Causent case governance on Drunix.
type AssuranceContract struct {
	contractapi.Contract
}

type caseRecord struct {
	CaseID        string `json:"case_id"`
	State         string `json:"state"`
	PolicyVersion string `json:"policy_version"`
	OwnerMSP      string `json:"owner_msp"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type evidenceRecord struct {
	CaseID               string `json:"case_id"`
	EvidenceID           string `json:"evidence_id"`
	EventDigest          string `json:"event_digest"`
	IssuerOrganizationID string `json:"issuer_organization_id"`
	EventType            string `json:"event_type"`
	NativeEventID        string `json:"native_event_id"`
	SubmitterMSP         string `json:"submitter_msp"`
	SubmittedAt          string `json:"submitted_at"`
}

type batchRecord struct {
	CaseID      string `json:"case_id"`
	BatchID     string `json:"batch_id"`
	BatchDigest string `json:"batch_digest"`
	MerkleRoot  string `json:"merkle_root"`
	SubmitterMSP string `json:"submitter_msp"`
	SubmittedAt string `json:"submitted_at"`
}

type assessmentRecord struct {
	CaseID            string `json:"case_id"`
	AssessmentID      string `json:"assessment_id"`
	AssessmentDigest  string `json:"assessment_digest"`
	RuleSetVersion    string `json:"rule_set_version"`
	AssuranceOutcome  string `json:"assurance_outcome"`
	SubmitterMSP      string `json:"submitter_msp"`
	SubmittedAt       string `json:"submitted_at"`
}

type challengeRecord struct {
	CaseID       string `json:"case_id"`
	ChallengeID  string `json:"challenge_id"`
	TargetType   string `json:"target_type"`
	TargetID     string `json:"target_id"`
	ReasonCode   string `json:"reason_code"`
	ChallengerMSP string `json:"challenger_msp"`
	SubmittedAt  string `json:"submitted_at"`
}

type correctionRecord struct {
	CaseID           string `json:"case_id"`
	CorrectionID     string `json:"correction_id"`
	PriorRecordID    string `json:"prior_record_id"`
	CorrectionDigest string `json:"correction_digest"`
	SubmitterMSP     string `json:"submitter_msp"`
	SubmittedAt      string `json:"submitted_at"`
}

type sealRecord struct {
	CaseID     string `json:"case_id"`
	SealDigest string `json:"seal_digest"`
	SealedBy   string `json:"sealed_by"`
	SealedAt   string `json:"sealed_at"`
}

func caseKey(caseID string) string {
	return fmt.Sprintf("CASE~%s", caseID)
}

func evidenceKey(caseID, evidenceID string) string {
	return fmt.Sprintf("EVIDENCE~%s~%s", caseID, evidenceID)
}

func evidenceIdempotencyKey(caseID, nativeEventID, issuerOrg string) string {
	return fmt.Sprintf("EVIDEMP~%s~%s~%s", caseID, issuerOrg, nativeEventID)
}

func batchKey(caseID, batchID string) string {
	return fmt.Sprintf("BATCH~%s~%s", caseID, batchID)
}

func assessmentKey(caseID, assessmentID string) string {
	return fmt.Sprintf("ASSESSMENT~%s~%s", caseID, assessmentID)
}

func challengeKey(caseID, challengeID string) string {
	return fmt.Sprintf("CHALLENGE~%s~%s", caseID, challengeID)
}

func correctionKey(caseID, correctionID string) string {
	return fmt.Sprintf("CORRECTION~%s~%s", caseID, correctionID)
}

func sealKey(caseID string) string {
	return fmt.Sprintf("SEAL~%s", caseID)
}

func (c *AssuranceContract) OpenCase(ctx contractapi.TransactionContextInterface, caseID string, policyVersion string) error {
	if err := requireNonEmpty("caseID", caseID); err != nil {
		return err
	}
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}

	key := caseKey(caseID)
	existing, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("case %s already exists", caseID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	rec := caseRecord{
		CaseID:        caseID,
		State:         caseStateOpen,
		PolicyVersion: policyVersion,
		OwnerMSP:      msp,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return putJSON(ctx, key, rec)
}

func (c *AssuranceContract) SubmitAttestation(ctx contractapi.TransactionContextInterface, caseID, evidenceID, eventDigest, issuerOrg, eventType, nativeEventID string) error {
	if err := requireNonEmpty("caseID", caseID); err != nil {
		return err
	}
	if err := requireNonEmpty("evidenceID", evidenceID); err != nil {
		return err
	}
	if err := requireNonEmpty("eventDigest", eventDigest); err != nil {
		return err
	}
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}

	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State == caseStateSealed {
		return fmt.Errorf("case %s is sealed; cannot submit attestation", caseID)
	}

	idemKey := evidenceIdempotencyKey(caseID, nativeEventID, issuerOrg)
	idemBytes, err := ctx.GetStub().GetState(idemKey)
	if err != nil {
		return err
	}
	if idemBytes != nil {
		return fmt.Errorf("duplicate native event for issuer in case: %s", nativeEventID)
	}

	evKey := evidenceKey(caseID, evidenceID)
	evBytes, err := ctx.GetStub().GetState(evKey)
	if err != nil {
		return err
	}
	if evBytes != nil {
		return fmt.Errorf("evidence ID %s already exists for case %s", evidenceID, caseID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	ev := evidenceRecord{
		CaseID:               caseID,
		EvidenceID:           evidenceID,
		EventDigest:          eventDigest,
		IssuerOrganizationID: issuerOrg,
		EventType:            eventType,
		NativeEventID:        nativeEventID,
		SubmitterMSP:         msp,
		SubmittedAt:          now,
	}
	if err := putJSON(ctx, evKey, ev); err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(idemKey, []byte(evidenceID)); err != nil {
		return err
	}

	if cr.State == caseStateOpen {
		cr.State = caseStateCollecting
		cr.UpdatedAt = now
		return putJSON(ctx, caseKey(caseID), cr)
	}
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) RegisterEvidenceBatch(ctx contractapi.TransactionContextInterface, caseID, batchID, batchDigest, merkleRoot string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State == caseStateSealed {
		return fmt.Errorf("case %s is sealed", caseID)
	}

	bKey := batchKey(caseID, batchID)
	if b, err := ctx.GetStub().GetState(bKey); err != nil {
		return err
	} else if b != nil {
		return fmt.Errorf("batch %s already registered", batchID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	batch := batchRecord{
		CaseID:       caseID,
		BatchID:      batchID,
		BatchDigest:  batchDigest,
		MerkleRoot:   merkleRoot,
		SubmitterMSP: msp,
		SubmittedAt:  now,
	}
	if err := putJSON(ctx, bKey, batch); err != nil {
		return err
	}
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) RecordAssessment(ctx contractapi.TransactionContextInterface, caseID, assessmentID, assessmentDigest, ruleSetVersion, assuranceOutcome string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State == caseStateSealed {
		return fmt.Errorf("case %s is sealed", caseID)
	}

	aKey := assessmentKey(caseID, assessmentID)
	if b, err := ctx.GetStub().GetState(aKey); err != nil {
		return err
	} else if b != nil {
		return fmt.Errorf("assessment %s already exists", assessmentID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	rec := assessmentRecord{
		CaseID:           caseID,
		AssessmentID:     assessmentID,
		AssessmentDigest: assessmentDigest,
		RuleSetVersion:   ruleSetVersion,
		AssuranceOutcome: assuranceOutcome,
		SubmitterMSP:     msp,
		SubmittedAt:      now,
	}
	if err := putJSON(ctx, aKey, rec); err != nil {
		return err
	}
	cr.State = caseStateAssessed
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) RaiseChallenge(ctx contractapi.TransactionContextInterface, caseID, challengeID, targetType, targetID, reasonCode string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State != caseStateAssessed && cr.State != caseStateDisputed {
		return fmt.Errorf("challenges require case in ASSESSED or DISPUTED state, got %s", cr.State)
	}

	chKey := challengeKey(caseID, challengeID)
	if b, err := ctx.GetStub().GetState(chKey); err != nil {
		return err
	} else if b != nil {
		return fmt.Errorf("challenge %s already exists", challengeID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	ch := challengeRecord{
		CaseID:        caseID,
		ChallengeID:   challengeID,
		TargetType:    targetType,
		TargetID:      targetID,
		ReasonCode:    reasonCode,
		ChallengerMSP: msp,
		SubmittedAt:   now,
	}
	if err := putJSON(ctx, chKey, ch); err != nil {
		return err
	}
	cr.State = caseStateDisputed
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) AppendCorrection(ctx contractapi.TransactionContextInterface, caseID, correctionID, priorRecordID, correctionDigest string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State == caseStateSealed {
		return fmt.Errorf("case %s is sealed", caseID)
	}

	cKey := correctionKey(caseID, correctionID)
	if b, err := ctx.GetStub().GetState(cKey); err != nil {
		return err
	} else if b != nil {
		return fmt.Errorf("correction %s already exists", correctionID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	rec := correctionRecord{
		CaseID:           caseID,
		CorrectionID:     correctionID,
		PriorRecordID:    priorRecordID,
		CorrectionDigest: correctionDigest,
		SubmitterMSP:     msp,
		SubmittedAt:      now,
	}
	if err := putJSON(ctx, cKey, rec); err != nil {
		return err
	}
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) RecordResolution(ctx contractapi.TransactionContextInterface, caseID, resolutionID, resolutionDigest string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	if !isAllowedParticipantMSP(msp) {
		return fmt.Errorf("submitter MSP %s is not an allowed participant", msp)
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State != caseStateDisputed {
		return fmt.Errorf("resolution requires DISPUTED state, got %s", cr.State)
	}
	_ = resolutionID
	_ = resolutionDigest
	_ = msp
	now := time.Now().UTC().Format(time.RFC3339)
	cr.State = caseStateResolved
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) SealCase(ctx contractapi.TransactionContextInterface, caseID, sealDigest string) error {
	msp, err := submitterMSP(ctx)
	if err != nil {
		return err
	}
	cr, err := c.loadCase(ctx, caseID)
	if err != nil {
		return err
	}
	if cr.State != caseStateAssessed && cr.State != caseStateResolved {
		return fmt.Errorf("sealing requires ASSESSED or RESOLVED state, got %s", cr.State)
	}
	if s, err := ctx.GetStub().GetState(sealKey(caseID)); err != nil {
		return err
	} else if s != nil {
		return fmt.Errorf("case %s already sealed", caseID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	seal := sealRecord{
		CaseID:     caseID,
		SealDigest: sealDigest,
		SealedBy:   msp,
		SealedAt:   now,
	}
	if err := putJSON(ctx, sealKey(caseID), seal); err != nil {
		return err
	}
	cr.State = caseStateSealed
	cr.UpdatedAt = now
	return putJSON(ctx, caseKey(caseID), cr)
}

func (c *AssuranceContract) GetCase(ctx contractapi.TransactionContextInterface, caseID string) (*caseRecord, error) {
	return c.loadCase(ctx, caseID)
}

func (c *AssuranceContract) VerifyCommitment(ctx contractapi.TransactionContextInterface, caseID, batchID, digest string) (map[string]interface{}, error) {
	bBytes, err := ctx.GetStub().GetState(batchKey(caseID, batchID))
	if err != nil {
		return nil, err
	}
	if bBytes == nil {
		return map[string]interface{}{"verified": false, "reason": "batch not found"}, nil
	}
	var batch batchRecord
	if err := json.Unmarshal(bBytes, &batch); err != nil {
		return nil, err
	}
	verified := batch.BatchDigest == digest
	return map[string]interface{}{
		"verified":      verified,
		"case_id":       caseID,
		"batch_id":      batchID,
		"committed":     batch.BatchDigest,
		"merkle_root":   batch.MerkleRoot,
	}, nil
}

func (c *AssuranceContract) loadCase(ctx contractapi.TransactionContextInterface, caseID string) (*caseRecord, error) {
	bytes, err := ctx.GetStub().GetState(caseKey(caseID))
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, fmt.Errorf("case %s not found", caseID)
	}
	var cr caseRecord
	if err := json.Unmarshal(bytes, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

func submitterMSP(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to read submitter MSP: %w", err)
	}
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("empty submitter MSP")
	}
	return id, nil
}

func isAllowedParticipantMSP(msp string) bool {
	switch msp {
	case "AgentPlatformMSP", "ServiceProviderMSP", "PaymentParticipantMSP", "Org1MSP", "Org2MSP":
		return true
	default:
		return false
	}
}

func requireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func putJSON(ctx contractapi.TransactionContextInterface, key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, bytes)
}
