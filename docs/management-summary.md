# Management summary

> **What this doc is.** Plain-English overview an examiner can hand to a bank CEO. Short. No jargon. Tells the CEO what was tested, what was found, what to do next.

## What was tested

The chain-of-custody system the bank operates for AI agent decisions. The system captures every AI agent decision and signs the day's collection of decisions with a tamper-evident seal that an outside party can verify.

Three things were checked:

1. Were any AI decisions altered after they were captured? (Did anyone change the record?)
2. Were the day-end seals signed by the bank's authorized hardware? (Are the seals genuine?)
3. Did the chain operate continuously through the examination period? (Was the system always on?)

## What was found

The findings fall into one of three buckets per day in the examination period:

| Bucket | Meaning | Examiner action |
|---|---|---|
| Pass | Everything verified. No issues for the day. | Note the result. |
| Pass with anomalies | Verified, but with operational notes (network blip, brief HSM unavailability). | Confirm the bank's incident log explains the anomaly. |
| Fail | The verifier could not confirm integrity for the day. | Investigate; require management response. |

The full report shows the per-day result. The summary on the cover page tells the CEO whether the result was overall pass, pass-with-anomalies, or fail.

## What the CEO needs to know

### If the result is pass

The system worked as designed. AI agent decisions captured during the examination period are integrity-bearing — no one altered them, and the bank can prove it. This is the headline.

The anomalies (if any) are operational items the bank's IT team handles in the normal course. They do not affect the integrity of the captured decisions. The bank should review them in its next quarterly control review and address any patterns.

### If the result is pass with anomalies that warrant attention

Same as above, plus: the examiner has identified specific anomaly patterns the bank should track. Typical patterns:

- Sealing delays beyond 24 hours: investigate the operational cause
- Late-binding event rate above the threshold: investigate network or SDK configuration
- Clock-skew anomalies: investigate time-synchronization

The examiner will follow up at the next examination cycle.

### If the result is fail

A specific day or set of days could not be verified. The bank must:

- Determine the root cause (the verifier output names the specific failure mode)
- Take immediate containment action if tampering is suspected
- Notify the regulator within 36 hours per the cyber-incident notification rule
- Produce a remediation plan with a timeline

The examiner provides the specific failure modes and the language for the management response.

## What the chain does NOT do

To be clear about scope:

- The chain does not check whether AI decisions were correct. It checks whether they were captured and unaltered. Whether the model made a good call is a separate matter — model risk management.
- The chain does not prevent bad AI behavior in real time. It catches bad behavior after the fact, with evidence.
- The chain does not replace the bank's broader cybersecurity posture. It is one control alongside many others.

## What this means for the bank

A passing chain-of-custody report is the bank's evidence that its AI decision-making is auditable. It satisfies a specific regulatory expectation (FFIEC IT Handbook II.C.10 on logging integrity) and supports the bank's broader model-risk and audit programs.

A bank that maintains a passing chain over time has a stronger position when:

- A customer disputes a decision (the bank can produce the integrity-bearing record)
- A regulator inquires about a specific decision or pattern (same)
- An external auditor reviews the bank's controls (the chain output is direct evidence)
- A vendor relationship changes (the chain's output is independent of any specific vendor)

A failing chain is a control failure that requires immediate response, the same as any other control failure of similar severity.

## Frequently asked questions

**"Why are we doing this if we already have logs from our AI vendor?"**

Vendor logs are the vendor's word that an event happened. The chain proves it. The distinction matters when the vendor's interests diverge from the bank's — in a regulatory inquiry, in a customer dispute, in litigation. The chain produces evidence that does not depend on the vendor's continued cooperation.

**"What is the ongoing cost?"**

The dominant cost is the hardware security module (HSM). Mid-size banks typically run $80k–$200k/year all-in. Smaller banks can run as little as $25k/year with a relaxed cadence. The cost picture is in `docs/cost-model.md`.

**"What if our vendor changes their implementation?"**

The chain is open-standards. Past data the bank captured is verifiable forever; future data uses whatever conforming implementation the bank chooses. Vendor lock-in is structurally minimized.

**"Do we have to share our master key with the regulator?"**

No. The verifier works against the bank's published public key, which the regulator already has. The master HMAC key stays with the bank.

## Next steps

After this examination:

1. Confirm any management responses required by the report
2. Schedule the next chain verification (typically aligned with the next examination cycle or the bank's audit cadence)
3. Update the bank's control description to reflect any changes to the chain implementation
4. Continue to monitor the operational events the chain produces

For questions about the report, the bank's contact is the examiner who issued it. For questions about the chain itself, the bank's CISO or chain-operations lead.
