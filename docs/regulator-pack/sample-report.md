# Sample verifier report

> **What this doc is.** A completed verifier report from a fictional tenant, for examiner training. Demonstrates the report structure, the pass/fail per-day detail, the anomaly section, and how to read each.

## Tenant facts (fictional)

- Tenant: `tenant_acme_prod_us_east_1`
- Examination period: 2026-04-01 through 2026-04-30
- Spec version: v1.0
- AI agent use case: Customer-service routing (low-stakes; account-level info only)
- Cadence: Daily
- Topology: BYOC (bank-operated AWS account)

## Verifier invocation

The standard examiner invocation:

```
verifier-validate.sh ./verifier
verifier verify \
  --ledger ./acme-snapshot-2026-04-30.dump \
  --root-key ./acme-public.pem \
  --master-key ./acme-ikm.bin \
  --tenant-id tenant_acme_prod_us_east_1 \
  --from 2026-04-01 \
  --to 2026-04-30 \
  --report ./acme-2026-04.pdf \
  --json-report ./acme-2026-04.json \
  --bundle ./acme-2026-04-bundle.tar.gz
```

The `--master-key` flag points at the institution's IKM file (32 raw bytes, file mode 0600); without it the verifier performs structural verification only and skips per-event HMAC equality (spec §7 fail-closed degradation per `--master-key absent`). Under `--strict` the absent IKM elevates to FAIL (`--strict requires key-bound verification`); for FFIEC examiner work the `--master-key` is provided per the institution's IKM-disclosure shape (per `customer-dispute-procedures.md` §"IKM access for customer-side verification" and `legal-disclosure.md` §"Court-ordered master-key disclosure" — both apply at examination time when key-bound verification is required).

Exit status: 0 (overall pass with anomalies)

The SOC engagement variant adds `--strict`, which elevates anomalies to FAIL and refuses dev-mode seals:

```
verifier verify --strict \
  --ledger ./acme-snapshot-2026-04-30.dump \
  --root-key ./acme-public.pem \
  --tenant-id tenant_acme_prod_us_east_1 \
  --from 2026-04-01 \
  --to 2026-04-30 \
  --report ./acme-2026-04-strict.pdf \
  --json-report ./acme-2026-04-strict.json
```

Use `--strict` for substantive SOC testing; without it for FFIEC examiner work where anomalies are evaluated in operational context.

## Cover page (Page 1)

```
FFIEC AI CHAIN-OF-CUSTODY VERIFICATION REPORT

Verifier version:     v1.0.4
Spec version:         v1.0
Generation timestamp: 2026-05-02 09:15:42 UTC

Ledger snapshot:      acme-snapshot-2026-04-30.dump
                      SHA-256: 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08

Public key:           acme-public.pem
                      Fingerprint: SHA256:HQ1cLP3OwY/...UZNiMjA=

Tenant:               tenant_acme_prod_us_east_1
Verification range:   2026-04-01 to 2026-04-30 (30 days)

SUMMARY
  Days verified:        30
  Days passed:          30
  Days failed:          0
  Events total:         4,218,391
  Runs total:           312,401
  Late-binding events:  47 (0.001%)
  Sealing delays > 1h:  3 days (see anomaly section)

OVERALL RESULT: PASS WITH ANOMALIES

Examiner signature: _________________________________
Date: _____________________
```

## Per-day detail (excerpt — page 2)

```
2026-04-01
  Events:                142,087
  Runs:                  10,423
  Late-binding count:    2
  Key versions present:  [3]
  Key fingerprints:      [b94c...4989]
  KMS handle URI:        aws-kms:arn:aws:kms:us-east-1:111122223333:key/abc-...
  Merkle root (computed): a7f5...e391
  Recorded seal root:     a7f5...e391
  hkdf_inputs_digest:     6f8a...7d65
  Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12
  Merkle match:           PASS
  Signature verification: PASS
  HMAC chain walk:        PASS (142,087 / 142,087 events verified)
  Anomalies:              None

2026-04-02
  Events:                141,553
  Runs:                  10,381
  Late-binding count:    1
  Key versions present:  [3]
  Key fingerprints:      [b94c...4989]
  KMS handle URI:        aws-kms:arn:aws:kms:us-east-1:111122223333:key/abc-...
  Merkle root (computed): 8b2c...a4d1
  Recorded seal root:     8b2c...a4d1
  hkdf_inputs_digest:     6f8a...7d65
  Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12
  Merkle match:           PASS
  Signature verification: PASS
  HMAC chain walk:        PASS
  Anomalies:              None

[... days 03–14 omitted ...]

2026-04-15  ← MASTER-KEY ROTATION DAY
  Events:                140,108
  Runs:                  10,289
  Late-binding count:    9
  Key versions present:  [3, 4]                  ← rotation day; both generations present
  Key fingerprints:      [b94c...4989, 2eed...8537]
  KMS handle URI:        aws-kms:arn:aws:kms:us-east-1:111122223333:key/abc-...
  Merkle root (computed): 5e1a...c73f
  Recorded seal root:     5e1a...c73f
  hkdf_inputs_digest:     6f8a...7d65
  Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12
  Merkle match:           PASS
  Signature verification: PASS
  HMAC chain walk:        PASS (events under v3 verified against ikm_v3;
                                 events under v4 verified against ikm_v4)
  Anomalies:              master_key_rotation_observed (rotation completed
                            03:14 UTC per institution incident log;
                            both key generations correctly resolved at
                            verifier step 7 IKM lookup)
                          Sealing delay 3h 12m (HSM cluster failover; per
                            institution incident log)
                          Late-binding rate elevated (0.006%, threshold 0.005%)

[... days 16–30 omitted ...]
```

## Anomalies section (page N+1)

```
ANOMALIES (3 days affected)

2026-04-15: Sealing delay 3h 12m
  The seal job for 2026-04-15 ran at 03:12:18 UTC on 2026-04-16 instead of
  the expected ~01:00 UTC. The institution's incident log records an HSM
  cluster failover during the affected window, with the seal job retrying
  through the failover. The failover completed successfully; the seal is
  valid. No further action recommended.

2026-04-15: Late-binding rate 0.006%
  9 events arrived after the daily seal had been computed; included in
  2026-04-16's seal. The institution's incident log notes a brief network
  partition between an SDK pod and the ledger during the affected window.
  Operational signal, not integrity signal.

2026-04-22: Sealing delay 1h 18m
  Within tolerance; no concurrent operational event recorded. Likely
  attributable to seal job scheduling jitter. No action recommended.
```

## Methodology section (page N+2)

```
ALGORITHMS USED
  Chain key derivation:  HKDF-SHA-256 (RFC 5869)
                         info = HKDF_INFO_BASE || "|" || utf8(tenant_id)
                         HKDF_SALT      = "ffiec.chain-of-custody.v1.salt"
                         HKDF_INFO_BASE = "ffiec.chain-of-custody.v1.info"
  Chain hash:            HMAC-SHA-256 (FIPS 198-1; RFC 4868 keying floor)
  Merkle:                RFC 6962 binary Merkle, SHA-256, leaf-prefix 0x00
  Signature:             Ed25519 (FIPS 186-5)
  Constant-time compare: per spec §10.8 (fingerprint AND MAC)

ORDER OF OPERATIONS (spec §7 twelve-step procedure)
  Header pre-flight:
    1. format_version match (most-specific-first refusal)
    2. hkdf_inputs_digest match against running per-tenant constants
    3. genesis_hash match against v1 constant (32 zero bytes)
  Per event (in run_id, seq ASC order):
    4. event.tenant_id == header.tenant_id; event.run_id == header.chain_id
    5. entry.format_version == header.format_version
    6. structural seq + 1 walk; entry.prev_hash == previous payload_hash
    7. ikm_lookup(tenant_id, key_version) returns IKM (NO MAC compute on miss)
    8. fingerprint check: SHA-256(utf8(tenant_id)||ikm)[:16] == entry.key_fingerprint
       (NO MAC compute on mismatch — load-bearing rotation defence)
    9. MAC recompute using EXPECTED prev_hash (not entry.prev_hash):
       expected_mac = HMAC(session_key, expected_prev_hash || canonical_bytes)
  Per day:
    10. Streaming Merkle root over ordered payload_hashes; compare to seal
    11. Reconstruct sign_payload (algorithm||format_version||tenant_id||
        date||hex(merkle_root)||hex(hkdf_inputs_digest)); verify Ed25519
    12. Cadence and dev_mode posture check (under --strict refuses dev_mode)

PASS/FAIL RULES
  - Steps 1-11 failure              => day FAIL (with named reason + step #)
  - Step 12 cadence mismatch        => day PASS with anomaly (Observation)
  - Step 12 dev_mode under --strict => day FAIL
  - master_key_rotation_observed    => day PASS with anomaly (normal-ops)
  - Late-binding rate > threshold   => day PASS with anomaly
  - Sealing delay > 60 min          => day PASS with anomaly
  - Sealing delay > 24 hours        => day PASS with anomaly (escalate)
  - Sealing delay > 72 hours        => day PASS with anomaly (notification required)
  - audit file ends mid-line        => day FAIL (mid-write truncation refusal)
```

## Sample failed-day output (page N+3 — paired example)

For comparison: what a single-day FAIL looks like under the same report shape, with an adjacent PASS day for visual reference. (Hypothetical; not part of the 30-day pass example above.)

```
2026-04-22  ← PASS DAY (shown for comparison; same April examination period)
  Events:                139,201
  Runs:                  10,178
  Late-binding count:    1
  Key versions present:  [4]
  Key fingerprints:      [2eed...8537]
  KMS handle URI:        aws-kms:arn:aws:kms:us-east-1:111122223333:key/abc-...
  Merkle root (computed): 4a8c...e7f2
  Recorded seal root:     4a8c...e7f2
  hkdf_inputs_digest:     6f8a...7d65
  Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12
  Merkle match:           PASS
  Signature verification: PASS
  HMAC chain walk:        PASS
  Anomalies:              None

2026-04-23  ← FAILED DAY (same April examination period)
  Events:                139,402
  Runs:                  10,201
  Late-binding count:    3
  Key versions present:  [4]
  Key fingerprints:      [2eed...8537]
  KMS handle URI:        aws-kms:arn:aws:kms:us-east-1:111122223333:key/abc-...
  Merkle root (computed): 9c2d...4b1a
  Recorded seal root:     9c2d...4b1a
  hkdf_inputs_digest:     6f8a...7d65
  Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)
  Merkle match:           N/A (chain walk failed before reaching step 10)
  Signature verification: N/A
  HMAC chain walk:        FAIL

  FAILURE RECORD
    step:                  8
    reason:                key_fingerprint mismatch at seq 17:
                           looked-up IKM does not match the entry's recorded
                           fingerprint
    run_id:                r_a3f29b71c
    seq:                   17
    tenant_id:             tenant_acme_prod_us_east_1
    key_version:           1
    expected_fingerprint:  b94c1a77b40bf5106c66ca6c1c1b4989
    recorded_fingerprint:  2eede65f0f764c97eaf3b3f306a48537
```

Cover page summary changes accordingly:

```
SUMMARY
  Days verified:        30
  Days passed:          29
  Days failed:          1                    ← non-zero
  Events total:         4,218,391
  Runs total:           312,401
  Late-binding events:  47 (0.001%)
  Sealing delays > 1h:  3 days

  FAILED DAYS:
    2026-04-23  step 8 (key_fingerprint mismatch)
                investigate institution's IKM roster row
                (tenant=tenant_acme_prod_us_east_1, key_version=1)
                per regulator-pack/examination-response-workflow.md

OVERALL RESULT: FAIL
```

The examiner's response path for this failed day follows `regulator-pack/examination-response-workflow.md` step 1 (severity from finding-language.md `step:8` row) → step 2 (IR Scenario 7) → step 3 (institution's most recent `master.reconciliation_completed`) → step 4 (lift the key_fingerprint mismatch finding paragraph) → step 5 (institution provides root cause + remediation + re-verification + control update).

## How an examiner reads this

1. **Cover page summary** — the headline. 30 days verified, 30 passed, 0 failed. Result: pass with anomalies.
2. **Per-day detail** — confirms each day's PASS status. Spot-check days where anomalies were recorded.
3. **Anomaly section** — the operational story. Were the anomalies explained by the institution's incident log? In this fictional example, two of three are explained (HSM failover, network partition). The third is small and within tolerance.
4. **Bundle** — the working-paper artifact. The examiner files the bundle's SHA-256 in IT-EX and retains the bundle as the examination's audit trail.

## Sample finding language for this report

If the examiner finds the report and the institution's explanations satisfactory:

> The institution's chain-of-custody verification for April 2026 produced an overall PASS result. Three minor sealing-delay anomalies were noted; the institution's incident log adequately explains each. The 2026-04-15 master_key_rotation_observed anomaly is normal-operations behaviour for a documented IKM rotation. No findings result.

If the late-binding rate had exceeded a defined threshold for multiple consecutive days:

> The institution's chain-of-custody verification for April 2026 produced PASS results with persistent late-binding event rate above the threshold of [X]%. The institution should investigate the operational cause and document remediation in its next quarterly control review.

## SOC team appendix — verifier line → TSC criterion mapping

For SOC 2 engagements consuming this report as substantive evidence behind PI1.1 / PI1.2 and ancillary CC subcriteria, map verifier output lines to TSC criteria as follows:

| Verifier output line | TSC criterion / SOC procedure |
|---|---|
| `Spec §7 steps executed: 1..12` | PI1.1 (processing integrity — input completeness); evidence the verifier exercised every defense-in-depth step |
| `Merkle match: PASS` | PI1.2 (processing integrity — output completeness); the day's events transitively produce the sealed root |
| `Signature verification: PASS` | CC6.7 (cryptographic controls), CC6.8 (prevent unauthorized modification); the HSM-signed root is integrity-bearing |
| `HMAC chain walk: PASS` | PI1.1; chain links from each event to the previous |
| `Key versions present: [v]` + `Key fingerprints: [hex]` | CC6.1 (logical access — key identity); P-6 reconciliation evidence at the audit-procedures level |
| `hkdf_inputs_digest: [hex]` | CC7.2 (system monitoring — format-drift detection) |
| `KMS handle URI` (non-`plaintext-` prefix) | CC6.7 (HSM custody); verification that production seals are HSM-backed |
| `Anomalies: master_key_rotation_observed` | CC6.7 (key lifecycle); CUEC-CRY-04 evidence the institution operates documented rotation |
| `Anomalies: sealing delay (within tolerance)` | A1.2 (availability — operational resilience); not an integrity finding |
| `FAILURE RECORD step: 8 key_fingerprint mismatch` | CC6.1 + CC6.7; identity-mismatch finding at the IKM-roster layer |
| `FAILURE RECORD step: 9 payload_hash MAC mismatch` | PI1.1 + CC6.8; content-tampering finding at the chain layer |
| `FAILURE RECORD step: 10 merkle root mismatch` | PI1.2 + CC6.8; ledger-content-tampering finding |
| `FAILURE RECORD step: 11 signature verification failed` | CC6.7; HSM-key-compromise or registry-mismatch finding |
| `FAILURE RECORD step: 1 format_version not supported by this verifier` | (Not a TSC finding — verifier-version skew; examiner obtains updated verifier; not institutional) |
| `FAILURE RECORD step: 2 header HKDF inputs do not match running v1 inputs` | CC6.7 (cryptographic controls — HKDF parameter integrity); CC8.1 if traced to SDK build drift |
| `FAILURE RECORD step: 3 header genesis_hash does not match v1 constant` | CC6.7 + CC8.1 (format-construction defect at SDK build) |
| `FAILURE RECORD step: 4 cross-chain lift detected at seq N` | CC6.7 + CC9.2 (vendor management — evidence-handling) |
| `FAILURE RECORD step: 5 format_version mismatch at seq N` | CC8.1 (change management — SDK format_version handling) |
| `FAILURE RECORD step: 6 chain link broken at seq N` | PI1.1 + CC6.8; structural chain integrity broken (insertion/deletion) |
| `FAILURE RECORD step: 7 unknown key_version: no IKM for ...` | CC6.1 + ID.AM-08 — IKM-roster retention/provisioning failure |
| `FAILURE RECORD step: 12 cadence mismatch` | CC8.1 (control-description-accuracy) |
| `FAILURE RECORD step: 12 dev-mode seal in production verification — refused` | CC6.8 + CC8.1 (compile-time exclusion bypassed; production-config defect) |
| `FAILURE RECORD (file pre-flight) audit file ends mid-line` | A1.2 (availability — operational resilience); recovery via SDK-local SQLite or upstream OTLP |
| `Late-binding count: N` (per-day) | PI1.2 (processing integrity — output completeness); the event is sealed in the next day's seal, not lost |
| `KMS handle URI` with `plaintext-` prefix under `--strict` | CC6.8 (prevent unauthorized software); verifier-side refusal of dev-mode seals; pairs with spec §10.7 compile-time exclusion |

The SOC engagement working paper carries the verifier output verbatim; the appendix above is the cross-walk the engagement team uses to allocate evidence to TSC criteria.
