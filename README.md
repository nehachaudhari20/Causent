# Causent: Beyond the Receipt — Execution Assurance for Agentic Payments

**Causent** is an execution-assurance platform for multi-step, multi-agent commerce. It connects a user’s original intent and delegated authority to observable agent actions, payment events, and service outcomes across independently operated systems.

Verifying a single authorization or payment receipt does not establish whether an entire agentic execution remained consistent with the user’s intent, delegated authority, and constraints. Causent checks whether observed execution complies with explicit constraints and identifies the **earliest divergence supported by the available evidence**—and abstains when the records do not justify a stronger conclusion.

Drunix provides the **permissioned cross-organisation governance layer**: attributable evidence commitments, challenges, corrections, and case sealing. Sensitive payloads and the assurance algorithm itself stay off-chain.

---

## Table of contents

- [Problem](#problem)
- [What Causent is and is not](#what-causent-is-and-is-not)
- [Core capabilities](#core-capabilities)
- [System architecture](#system-architecture)
- [Assurance outcomes and uncertainty](#assurance-outcomes-and-uncertainty)
- [Example: amount divergence](#example-amount-divergence)
- [Drunix integration](#drunix-integration)
- [Prototype and validation](#prototype-and-validation)
- [Repository layout](#repository-layout)
- [Implementation roadmap](#implementation-roadmap)
- [Documentation](#documentation)
- [Status](#status)

---

## Problem

Agentic commerce introduces a critical **accountability gap**. As workflows become multi-step and span multiple agents, service providers, and payment participants, execution evidence fragments across independently maintained systems. Dependencies between delegated actions, payment execution, and downstream service outcomes are hard to establish across organisational boundaries. Incorrect selections, unauthorized actions, and inconsistent outcomes can propagate through dependent services, obscuring where failures began.

Existing authorization frameworks and payment protocols provide identity, mandates, payment verification, and receipts—but those records do not inherently yield a unified, **evidence-backed reconstruction** of the full execution lifecycle. Missing events, delayed confirmations, ambiguous references, and contradictory attestations further complicate verification and dispute resolution.

The gap is the absence of a reliable, **evidence-bounded** method for cross-participant execution assurance that:

- correlates independently sourced evidence,
- verifies explicit execution constraints,
- identifies the earliest **observable and evidence-supported** divergence,
- produces a reproducible investigation record,

while clearly distinguishing **confirmed violations** from **missing evidence**, **ambiguity**, and **conflicting claims**.

```text
User intent & authority
        |
        v
Agent platform ---- Delegated agents ---- Service providers
        |                    |                    |
        +--------------------+--------------------+
                             |
                    Payment participant
                             |
        Evidence fragmented; no single system
        holds the full execution story
                             |
                             v
              Causent: correlate, constrain, localise, abstain
```

Causent **complements** existing payment and authorization controls; it does not replace them.

---

## What Causent is and is not

| In scope | Outside the core claim |
|----------|-------------------------|
| Correlate authenticated evidence across participants | Recover an agent’s private LLM reasoning |
| Verify explicit authority, amount, scope, and lifecycle constraints | Replace AP2, x402, MPP, UAP, a PSP, or a payment rail |
| Model dependencies between delegation, invocation, payment, and delivery | Decide legal liability or assign final blame |
| Localise the earliest **evidence-supported** divergence | Guarantee discovery of the true root cause from incomplete records |
| Represent missing, ambiguous, and contradictory evidence explicitly | Automatically reverse payments or unwind completed services |
| Record cross-organisation assertions, challenges, corrections, and case commitments on Drunix | Claim production NPCI, UPI, bank, or PSP connectivity without authorization |

**Design distinctions (non-negotiable):**

- A valid digital signature shows that a particular key signed a representation; it does **not** prove the underlying business claim is true.
- A missing event is **not** proof that the event never occurred.
- Causent establishes what the evidence **supports**; it does not treat missing records as proof of absence or infer violations from similar amounts and timestamps alone.

**Positioning:** Causent is not “AI payments plus blockchain,” a trust score for agents, or a graph that automatically finds root cause. Its moat is the combination of **protocol-aware evidence adapters**, **typed execution relationships**, **uncertainty-aware localisation (EBDL)**, and **cross-organisation governance**—evaluated against strong baselines, not assumed.

---

## Core capabilities

### 1. Canonical evidence and cross-system correlation

Causent collects structured evidence from agent platforms, delegated agents, service providers, and payment participants. Each record preserves source identity, native event reference, timestamps (`occurred_at` vs `observed_at`), task/delegation/payment references, payload digest, and cryptographic verification metadata.

A **correlation engine** connects records using:

1. Native protocol bindings (strongest),
2. Authenticated explicit references,
3. Trace or execution context,
4. Constrained candidate matching (never silently promoted to verified).

Duplicate events, invalid signatures, conflicting records, and unresolved references are classified separately. Uncertain matches are never treated as verified dependencies.

### 2. Execution dependency model — Multi-Layer Provenance Graph (MLPG)

Causent builds an **MLPG** linking user intent, delegated authority, agent invocations, resource selection, payment authorization, settlement, service fulfilment, and downstream effects.

- **Nodes:** events or execution entities.
- **Edges:** typed relations (`delegated-by`, `authorized-by`, `invoked`, `settled-for`, `fulfilled-by`, etc.) with **evidence basis** and **verification status**.

Execution order follows **supported dependencies**, not wall-clock timestamps alone—handling delayed events, concurrency, and incomplete records without inventing a sequence the evidence cannot support.

### 3. Evidence-Bounded Divergence Localisation (EBDL)

**EBDL** is Causent’s core assurance pipeline:

1. Validate evidence (schema, signatures, issuer authority, replay/idempotency).
2. Resolve relations conservatively.
3. Build a partial-order execution graph; detect cycles and unresolved edges.
4. Evaluate versioned, deterministic assurance rules against recorded intent and delegated authority (amount ceilings, permitted providers, scope, delegation boundaries, payment–fulfilment relationships).
5. Identify **minimal confirmed violations** under the supported causal order.
6. Assess whether localization is unique, multi-candidate, or **unresolved**.
7. Emit a reproducible assessment (evidence digests, rule versions, relation justifications, unresolved questions).

Rule evaluation per constraint: `SATISFIED` · `VIOLATED` · `UNKNOWN` · `CONFLICT`. The system does **not** upgrade `UNKNOWN` to `SATISFIED` or turn candidate links into confirmed violations.

### 4. Assurance results and reproducible investigation

Each assessment produces a structured outcome and an auditable report. Investigators can trace every finding to supporting evidence and reproduce the run with the same evidence set and rule versions.

Challenges and corrections are **new attributable records**; prior history is not silently rewritten.

### 5. Drunix-based cross-organisation governance and privacy

Drunix holds **minimum on-chain governance state**: case identifiers, evidence-batch commitments, assessment digests, challenges, corrections, and sealed case versions. Encrypted payloads, credentials, personal data, and raw logs remain in off-chain storage with separate access control.

A committed digest supports **integrity checking** against registered data; it does not, by itself, prove that the underlying claim is true.

---

## System architecture

```text
Participant systems
  Agent platform · Service provider · Payment adapter · Downstream consumer
        |                    native receipts, exports, attestations
        v
+------------------------------------------------------------------+
| 1. Evidence ingestion & verification                              |
|    Schema · issuer identity · signatures · replay · authority     |
+------------------------------------------------------------------+
        v
+------------------------------------------------------------------+
| 2. Correlation & execution model                                  |
|    Canonical envelope · typed relations · protocol adapters · MLPG|
+------------------------------------------------------------------+
        v
+------------------------------------------------------------------+
| 3. Causent assurance engine (EBDL)                                |
|    Rules · evidence-qualified relations · divergence · abstention |
+------------------------------------------------------------------+
        |
        +------------------------+---------------------------+
        v                        v                           v
 Off-chain evidence plane   Drunix governance plane   Investigation API
 (encrypted objects,        (cases, commitments,         (timeline, graph,
  graph, reports)            challenges, seals)          report export)
```

| Component | Responsibility |
|-----------|----------------|
| Protocol adapters | Normalize protocol-specific evidence; preserve native references |
| Evidence verifier | Cryptographic integrity, issuer authority, freshness, replay |
| Correlation engine | Relations only at supported confidence/evidence levels |
| Graph builder | Partial order, causal edges, unresolved links |
| Assurance engine (EBDL) | Deterministic rules, localization, abstention |
| Evidence store | Encrypted raw objects and reports off-chain |
| Commitment service | Canonical batches, Merkle roots, digests |
| Drunix chaincode (`causent-assurance`) | Case governance, append-only records, endorsement policies |
| Application API | FastAPI: validation, correlation, EBDL, Drunix client integration |

**Application stack (target):** Python/FastAPI for ingestion and assurance; PostgreSQL for event index and graph metadata; encrypted object storage for payloads; **Go chaincode** on Drunix for governance.

---

## Assurance outcomes and uncertainty

Causent maintains **separate** assessments—not a single trust score.

**A. Evidence status**  
`AUTHENTICATED` · `UNVERIFIED_SOURCE` · `INVALID_EVIDENCE` · `DUPLICATE_EVENT` · `CONFLICTING_EVENT_ID`

**B. Relationship status**  
`EXPLICIT_BINDING` · `AUTHENTICATED_REFERENCE` · `CANDIDATE_LINK` · `AMBIGUOUS_LINK` · `UNLINKED`

**C. Assurance outcome**  
`CONSISTENT_WITH_EVIDENCE` · `DIVERGENCE_CONFIRMED` · `INSUFFICIENT_EVIDENCE` · `CONFLICTING_ATTESTATIONS` · `NOT_ASSESSABLE`

**D. Localization status**  
`UNIQUE_EARLIEST` · `MULTIPLE_MINIMAL_CANDIDATES` · `EARLIEST_UNRESOLVED`

When multiple violations cannot be reliably ordered, Causent reports **multiple possible earliest divergences**. When missing or conflicting evidence prevents localisation, it marks the earliest divergence as **unresolved**.

---

## Example: amount divergence

- Verified authorization establishes a ceiling of **₹9,500**.
- Authenticated settlement evidence is **reliably linked** to that execution and records **₹12,000** settled.

**Causent can report:** amount-constraint violation (`DIVERGENCE_CONFIRMED`), with localization `UNIQUE_EARLIEST` only if preceding dependencies support that uniqueness.

If settlement **cannot** be reliably linked to the authorization, Causent reports **insufficient or ambiguous evidence**—not a violation inferred from similar amounts or nearby timestamps.

---

## Drunix integration

Causent uses Drunix as the **shared governance and evidence-integrity layer** across independently operated participants. **Lite Peers** simulate and endorse transactions; **Committing Peers** validate and commit; the **stateless validation service** checks endorsement policies. Application-level assurance rules (e.g. payment-limit violations) remain in Causent’s processing layer.

| Drunix feature | Causent use |
|----------------|-------------|
| Lite / Committing peer separation | Scale endorsement vs commit as evidence submissions grow |
| Stateless validation service | Enforce authorised transaction participation—not EBDL |
| YugabyteDB state DB (target demo) | Ledger world state: case status, batch refs, assessment commitments |
| KeyDB / private data | Transient/sensitive paths where configured; **not** permanent evidence archive |
| Go chaincode | Case lifecycle and policy enforcement |

### Organisations (prototype network)

| MSP | Role |
|-----|------|
| `AgentPlatformMSP` | Authority/delegation, task events, challenges |
| `ServiceProviderMSP` | Invocation, delivery, provider attestations |
| `PaymentParticipantMSP` | Authorization, settlement, failure, reversal (sandbox/simulated) |

Submitting client identity, original evidence issuer, assurance-engine identity, and endorsing peers are **distinct** and preserved in records.

### Proposed chaincode operations (`causent-assurance`)

| Operation | Purpose |
|-----------|---------|
| `OpenCase` | Initialise case and policy version |
| `SubmitAttestation` | Participant evidence reference and digest |
| `RegisterEvidenceBatch` | Canonical batch / Merkle root commitment |
| `RecordAssessment` | Assurance report commitment |
| `RaiseChallenge` | Challenge evidence or assessment (without target approval) |
| `AppendCorrection` | Superseding/corrected evidence linked to prior record |
| `SealCase` | Final reviewable case snapshot under governance policy |

Case states (target): `OPEN` → `COLLECTING_EVIDENCE` → `ASSESSED` → `DISPUTED` → `RESOLVED` → `SEALED`.

Endorsement policies differ by action: participant evidence submission, independent challenges, multi-organisation approval for resolution/sealing. The prototype will test rejection of unauthorised submissions, impersonation, invalid transitions, and bypass of challenge or resolution procedures.

### Data placement

| Store | Contents |
|-------|----------|
| Drunix ledger / YugabyteDB | Case state, digests, governance events |
| Off-chain PostgreSQL | Event index, graph relationships, metadata |
| Encrypted object store | Original evidence and full reports |

---

## Prototype and validation

### Demonstration scenarios (minimum)

| Case | Expected behaviour |
|------|-------------------|
| **A — Consistent execution** | Authority, payment, and service evidence agree under rules → `CONSISTENT_WITH_EVIDENCE` |
| **B — Confirmed amount divergence** | ₹9,500 ceiling vs ₹12,000 settled with reliable binding → `DIVERGENCE_CONFIRMED` |
| **C — Insufficient evidence** | Similar amount/time but no trustworthy link → ambiguous / insufficient, not a definitive mispayment |

Additional injected conditions: service-scope violations, missing links, delayed events, duplicate notifications, conflicting attestations, concurrent branches, invalid signatures.

### Evaluation methodology

Begin with **30–60** synthetic multi-hop scenarios (expand if evaluation is stable). Each case is labelled for the **injected fault** and the **conclusion supportable from available evidence**—these differ when upstream evidence is missing.

**Baselines (same evidence, equivalent constraints where supported):**

- Manual investigation  
- Flat event correlation  
- Protocol-aware flat assurance (no multi-hop dependency propagation)  
- Simpler dependency graph (weaker relation qualification / abstention)  
- **Full Causent (EBDL + MLPG + governance)**

**Assurance metrics:** divergence-localisation accuracy (unique cases), localisation coverage, relationship precision/recall, **unsupported definitive rate**, abstention quality, macro outcome F1.

**Investigation effort:** median time to reproducible outcome, manual lookups, unsupported conclusions (Causent vs baseline, counterbalanced).

**Drunix metrics:** policy enforcement, challenge independence, commitment integrity, latency/throughput vs a controlled central append-only comparator.

### Provisional evaluation targets (goals—not achieved results)

- ≥ **10 percentage points** higher correct localisation vs strongest baseline on uniquely localisable cases  
- ≥ **20%** lower median investigation time vs strongest baseline  
- ≤ **5%** unsupported definitive conclusions (denominator: independently labelled evidence-supported outcomes)  
- Reproducibility metadata on every assessment (evidence refs, rule versions, commitments)

Actual results, sample sizes, failures, and limitations will be reported honestly.

---

## Repository layout

Target structure (evolving with implementation):

```text
causent/
├── README.md
├── docs/
│   ├── architecture.md
│   ├── data-model.md
│   ├── drunix-design.md
│   ├── threat-model.md
│   └── evaluation.md
├── adapters/
│   └── sandbox_payment/
├── schemas/
│   ├── canonical_event.schema.json
│   └── assessment.schema.json
├── evidence/
│   ├── signature_verifier.py
│   ├── source_registry.py
│   └── replay_guard.py
├── correlation/
│   ├── relations.py
│   └── partial_order.py
├── assurance/
│   ├── rules.py
│   ├── ebdl.py
│   └── report.py
├── storage/
│   ├── evidence_store.py
│   └── postgres_repository.py
├── drunix/
│   └── chaincode/
│       └── causent-assurance/
├── api/
│   └── main.py
├── simulator/
│   ├── scenarios/
│   └── fault_injection.py
├── evaluation/
│   ├── baselines/
│   ├── metrics.py
│   └── benchmark_runner.py
└── tests/
    ├── test_evidence.py
    ├── test_assurance.py
    ├── test_policies.py
    └── test_end_to_end.py
```

Assurance logic remains **independently testable** without entangling the Drunix client in core rule evaluation.

---

## Implementation roadmap

| Phase | Focus | Acceptance gate |
|-------|--------|-----------------|
| **1 — Evidence & rules** | Canonical schema, sandbox adapter, registry, amount/authority rules, signed fixtures | Clean, excess-amount, invalid-signature, duplicate, ambiguous-link classifications |
| **2 — EBDL & report** | MLPG, partial order, rules, localization, abstention, reproducible report | Unique / multi-candidate / unresolved localization tests pass repeatably |
| **3 — Drunix** | Test network, `causent-assurance`, three orgs, case lifecycle, batch/assessment/challenge | Read committed state, verify commitment, reject ≥1 unauthorised operation |
| **4 — Benchmark** | Demo cases A–C, baselines, ablations, Drunix performance | Published metrics with config, denominators, and known limits |

Implementation labels: `DESIGNED` → `IMPLEMENTED` → `TESTED` → `BENCHMARKED`. Planned work is not described as shipped functionality.

---

## Documentation

| Document | Contents |
|----------|----------|
| [docs/architecture.md](docs/architecture.md) | Layering and component design (being aligned with Causent / EBDL) |
| [docs/data-model.md](docs/data-model.md) | Canonical objects and events |
| [docs/drunix-design.md](docs/drunix-design.md) | On-chain model, policies, chaincode |
| [docs/threat-model.md](docs/threat-model.md) | Threats, mitigations, privacy boundaries |
| [docs/evaluation.md](docs/evaluation.md) | Benchmarks, metrics, baselines |

---

## Local development

```bash
python -m pip install -e ".[dev]"
python -m pytest tests -q
python -m api.main   # or: uvicorn api.main:app --reload --port 8080
```

Drunix chaincode (requires Go 1.22+):

```bash
cd drunix/chaincode/causent-assurance
go build .
```

With `DRUNIX_GATEWAY_ENABLED=false` (default), the API uses stub ledger responses so you can develop without a live network. See [drunix/README.md](drunix/README.md).

---

## Status

| | |
|--|--|
| **Project** | Causent — Beyond the Receipt |
| **Domain** | Evidence-bounded execution assurance for agentic payments |
| **Core method** | EBDL + Multi-Layer Provenance Graph (MLPG) |
| **Governance** | Drunix (Go chaincode, multi-org prototype) |
| **Implementation** | Documentation and target design in progress; application and chaincode to follow roadmap |

---

## Core design principle

> **Do not stop at the receipt.** Correlate what participants actually attested, verify explicit constraints against recorded intent and authority, localise only the earliest divergence the evidence supports, say clearly when the evidence cannot support more—and preserve an auditable, cross-organisation record of that investigation.
