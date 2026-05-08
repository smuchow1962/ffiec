# Bonnie Hartwell — Plaintiff-side e-discovery review of FFIEC chain-of-custody v1.0a

**Date:** 2026-05-07
**Reviewer:** Bonnie Hartwell, Partner, plaintiffs' class-action practice (consumer financial fraud, RESPA/TILA, data-breach class actions)
**Posture:** First encounter with v1.0a. No prior iteration history. Reviewing as the opposing expert in a hypothetical class action where the defendant bank produces chain-of-custody artifacts to authenticate AI-driven adverse-action records against a class of consumer plaintiffs.

## How I read this

When a defendant bank tells me their AI logs are tamper-evident under a published cryptographic standard, I do not start by accepting the math. The math is usually fine. I start by asking who controls the production, who controls the verification, who controls the artifacts opposing counsel needs to mount a forgery argument, and where the chain interlocks with the bank's other systems. The seams are where the defense story breaks.

I read the spec, the verifier design, the FRE 902 template, the partial-disclosure document, the audit procedures, the customer-dispute procedures, and the threat model. The construction is more rigorous than what I have seen in production at any institution I have litigated against in the last decade. The cryptography is FIPS-grounded and the verifier procedure is testable. That cuts against my client and I would say so on the record at a Daubert hearing.

But the spec lives next to a litigation-support apparatus that the bank's IT custodian will sign certifications under. That is where the cross-examination targets are. Below are eight findings, framed as I would brief them to my expert and to the bench.

## Findings

### G-1 — The certification declarant has no specified independence from the verifier-run.

**Section/file:** `docs/templates/fre-902-certification.md` §A.1, §A.4; `spec/chain-of-custody-v1.md` §10.13

**Cross-examination angle.** The 902(13) template names a "records custodian" as the declarant. The declarant signs that the verifier ran on a specified date and produced a specified exit code. Nothing in the template — and nothing in §10.13 evidentiary artifacts — requires the declarant to be different from the person who ran the verifier, configured the master-key flag, controlled the files the verifier read, or had write access to the IKM registry. The `verifier.run_completed` operational event (§A.4 ¶14) is signed by the same HSM under the same institutional controls. There is no second-set-of-eyes requirement.

I would put the declarant on the stand and ask three questions. Did you run the verifier yourself? If not, who did, and were they reporting to you? Did anyone outside your reporting line independently re-run the verifier against the same bytes? The third question is where the certification falls apart, because the spec's defense-in-depth argument rests on independent verification by a regulator or auditor (§1.1, §5 of design 09), and the discovery production reaching me has not been touched by either. The certification is a self-certification. Under FRE 803(6)'s trustworthiness clause and *Lorraine v. Markel*'s authentication framework, that is a foundation challenge I can win, especially in a class-action context where the bank's interest in a passing verifier output is concrete.

**Motion-in-limine angle.** I would file a motion to exclude the certification under FRE 902(13) until the bank produces an independent verifier run by a non-conflicted party — the institution's external SOC 2 auditor, the OCC's examination team, or a court-appointed neutral. Alternatively, foundation testimony from an independent verifier-runner.

### G-2 — The "customer-side verification path" is performative; the meaningful path is opt-in by the defendant.

**Section/file:** `docs/customer-dispute-procedures.md` §"Customer-side verification path"; `docs/customer-dispute-procedures.md` §"IKM access for customer-side verification"

**Cross-examination angle.** The dispute-procedures document tells me the bank "should not resist" customer-side verification and that the design intent is exactly that customers can verify. That sentence reads well until I cross-examine it. The two acceptable shapes for IKM access are (1) a court-issued protective order — which means I have to litigate to one before I get the bytes — or (2) HSM-mediated verification through the bank's HSM API, which means my expert never holds the IKM and the bank's HSM is in the loop on every verification my expert runs.

Option 2 is the one the bank will offer first. It is also the one where the bank controls the verification infrastructure. If the bank's HSM rate-limits my expert's queries, refuses run identifiers as "expired," or dispatches the wrong session key for a contested entry, my expert has no independent path to detect that — the HSM is a black box on the defendant's side. Without a regulator-in-the-loop variant or a court-controlled IKM escrow, the bank holds the verification key and the verification infrastructure simultaneously. That is the textbook configuration for a one-sided integrity story.

**Motion-in-limine angle.** I would seek a Rule 26(c) protective-order modification placing the IKM in court-controlled escrow with my expert receiving access under sealed conditions, OR seek an order requiring the bank to designate a regulator-side verification path (e.g., the institution's SOC engagement firm running the verifier under their National Office independence procedure per the audit-procedures document §3) as a parallel verification track. Without one of those, the customer-side verification path is theatrical.

### G-3 — The partial-disclosure mode is a defendant's selective-production weapon and the spec acknowledges the limit only in its own favor.

**Section/file:** `docs/design/07-verifier-design.md` §10; `docs/selective-production-and-sampling.md` §1, §6

**Cross-examination angle.** The partial-disclosure mode is a beautiful instrument — for the producing party. The mode authenticates inclusion of the entries the producing party chose to disclose. It does not assert that the producing party disclosed every responsive entry. The companion document is honest about this — `This report does NOT attest chain completeness for tenant_day {D}` — and prints the limit in the output banner.

That admission is the cross-examination opening I want. The selective-production document tells me directly that "the consumer plaintiff assesses completeness through the order's specificity and through the receiver's compliance certification, not through the cryptography" (§6). Completeness is a representation, not a cryptographic fact. The receiver picks which `(run_id, seq)` triples are responsive. The cryptography validates only the receiver's chosen subset.

In a 2703(d) criminal context the order's specificity bounds the risk because the order names the records. In class-action discovery I do not have an order with that specificity — I have a request for production and the defendant's privilege log. The defendant ships a partial bundle of records the defendant deems responsive, the verifier confirms inclusion, the bench sees a clean PASS, and the records the defendant chose to suppress are not represented in the bundle in any form. The seal commits to a Merkle root over a leaf set whose size is not revealed by the audit path alone (§10.4 of design 07). The defendant controls what is produced AND the count of what existed — and a sophisticated adversary can shape the production to obscure entries the seal nevertheless covers.

**Motion-in-limine angle.** I would object to any FRE 902(13) certification that authenticates a partial-disclosure bundle as a substitute for full production unless the bundle is accompanied by (a) a sworn statement of the count of total entries in the seal-day plus the count produced (so the proportion of withheld entries is on the record), and (b) the disclosure manifest's basis field naming the search criteria the defendant ran against the ledger to identify responsive entries. The selective-production document mentions the basis field under `--strict` but it is not load-bearing for the cryptography (§9 of selective-production, last paragraph) — which means it is procedural, which means I can challenge it.

### G-4 — Spoliation under FRCP 37(e) is hardest to prove against an opponent who controls the evidence of preservation.

**Section/file:** `spec/chain-of-custody-v1.md` §10.9 IKM registry retention; §10.13 evidentiary artifacts; `docs/customer-dispute-procedures.md` §"Documentation for dispute response"

**Cross-examination angle.** Rule 37(e) sanctions for failure to preserve ESI require me to show that the ESI should have been preserved, that it was lost because reasonable steps were not taken, and that the loss prejudices my client. The 2015 amendments raised the bar for adverse-inference instructions to require a finding of intent to deprive.

Here is the structural problem. The institution retains chain entries and seal records under §10.13. The institution retains the IKM under §10.9 as long as any chain entry under that key_version is retained. The verifier returns PASS, the IT witness testifies the chain is intact. What the institution does NOT have to retain — under any provision I read — is the surrounding evidence I need to argue that responsive entries were omitted: the custodian's query logs, the record-retention configuration during the period at issue, the customer-correlation index entries, and the data-classification policy that decides whether substantive prompt/response is captured or only the SHA-256 hash. The 902 template's §C.2 ¶3 explicitly contemplates the case where "the Institution's data classification policy required exclusion of customer PII from the recorded prompt and response, and only the SHA-256 hash of those values is retained."

That last clause is the one I litigate hardest. The institution's data-classification policy decides whether substantive content is captured. If the policy excluded substantive content during the period at issue, the chain proves only that the institution received some bytes whose hash matches; the substantive content is gone, and the §C certification authenticates only the hash. The content was lost by policy, not by accident. Under Rule 37(e), policy-driven exclusion of substantive content from a record-keeping system otherwise designed to preserve evidence is a candidate for an adverse-inference instruction, especially when the policy was set by the same institution that benefits from the exclusion.

**Motion-in-limine angle.** I would seek discovery of the institution's data-classification policy at the time of capture, its change history, and the substantive content that flowed through the chain at the time of capture but was excluded by policy. I would seek an adverse-inference instruction on the basis that the institution chose to exclude substantive AI prompts and responses from the chain while otherwise designing the chain to be a litigation-grade integrity record, and that the exclusion choice was made under the institution's control during a period when class-member harms were foreseeable.

### G-5 — The "hash-only" 902(11) certification is the defendant's escape hatch and it authenticates almost nothing useful to me.

**Section/file:** `docs/templates/fre-902-certification.md` §C; `spec/chain-of-custody-v1.md` §1.2 epistemic scope

**Cross-examination angle.** The §C variant is the load-bearing fall-back for cases where the institution cannot reproduce the LLM's input or output. It authenticates that "the Institution received the bytes whose hash is recorded" and that "no alteration of the recorded bytes occurred between capture and verification." The certification disclaims attestation of factual accuracy, policy compliance, and freedom from bias (§C.2 ¶4).

Read as a defense instrument, this is masterful. It admits the chain proves only what it proves. It anchors the claim under penalty of perjury per 28 USC §1746. It will be admitted under FRE 902(11) without live foundation if I do not challenge it. It leaves a class-action plaintiff arguing that the defendant's AI made discriminatory adverse-action decisions with no record of what the AI actually said — only a hash that the defendant claims matches bytes the defendant claims it received.

The §C certification is not dishonest. It is precisely scoped, which is what makes it dangerous. Under FRE 803(6)'s trustworthiness clause, the certification proves bytes-as-received but does not prove what the bytes were. Without an independent reproduction path — a provider-side log per §C.5 ¶12 — the certification's value to my client is approximately zero. The value to the bank is that it authenticates the bank's narrative that the chain was operating as designed.

**Motion-in-limine angle.** I would file under Daubert/702 to exclude expert testimony that interprets the §C certification as evidence of what the AI said. The certification's own scope-limitation paragraph (§C.2 ¶3) is the exhibit. I would also seek a Rule 26 order requiring the institution to disclose whether substantive content was retained under a different system (provider-side logs, downstream system mirrors, regulator-side examination archives) and to produce that content if it exists.

### G-6 — Witness mode + 902(13) is a cross-examinable foundation, and the spec admits it.

**Section/file:** `docs/templates/fre-902-certification.md` §G "Witness-mode verification in regulator hand-offs"; `spec/chain-of-custody-v1.md` §7 witness-verifier mode

**Cross-examination angle.** Witness mode is the verifier's posture when the IKM is not provided. It executes the structural walk, Merkle recomputation, and HSM signature verification, but skips the per-event MAC checks. The §G paragraph is candid: "Opposing counsel may challenge witness-mode authentication on the theory that the per-event MACs were not verified."

I am opposing counsel. I will challenge it. The per-event MAC is the primitive that binds the entry to the institution's HSM-derived session key. Skipping it leaves me with structural integrity, the Merkle root, and the HSM signature — none of which prove per-entry IKM binding. The spec rebuts on three theories: cross-tenant binding via HKDF info, seal-level MAC binding via the Merkle root, and HSM signature coverage of the root (§G ¶3). That rebuttal is correct as far as it goes. It relies on every entry's `payload_hash` being in the Merkle leaf set — which is what the seal asserts but not what witness mode independently verifies. Witness mode recomputes Merkle over the sealed file's entries; a forged entry inside a sealed file passes if the institution sealed the forged file.

The defense's answer is the §1.4 compositional security argument — three independent layers, none alone sufficient. That argument holds where the three layers are independently controlled. The institution controls all three. The HSM is the institution's. The seal job is the institution's. The IKM custody is the institution's. Independence in §1.4 is operational separation-of-duties between institution roles; from a class-action plaintiff's perspective, those roles are all on the same side of the wall.

**Motion-in-limine angle.** I would object to the introduction of any witness-mode verifier output as authentication of the chain entries' content. Witness mode authenticates structural integrity and seal validity; it does not authenticate that the entries' bodies are what the institution captured at the time. I would require full-key verification by an independent runner before the chain entries come in for content (as opposed to existence).

### G-7 — Confirmation: the verifier exit code and reason strings are normative and that constrains the defendant's framing.

**Section/file:** `spec/chain-of-custody-v1.md` §7 "Failure-reason strings are normative byte-for-byte"; §10.12 verifier CLI exit-code contract; `docs/templates/fre-902-certification.md` §F

**Cross-examination angle.** This is a confirmation. The spec normates the verifier's failure-reason strings byte-for-byte. The IT witness cannot soften a FAIL into "a minor inconsistency the team is investigating" — the verifier produces `Status: FAIL`, `Step: N`, `Reason: <text>`, and the strings are normative reference. The 902 template's §G commits the institution to disclosing the FAIL exit code, step number, and reason string under penalty of perjury per 28 USC §1746. Producing a misleading PASS certification is perjury.

That is good for me when the verifier returns FAIL. It is good for the defense when the verifier returns PASS. Both sides are pinned to the exit code. I cannot manufacture a FAIL from a PASS, and the bank cannot soften a FAIL into a PASS. The exit-code contract is a non-negotiable handle, which is exactly what an integrity standard should produce. In cases where the verifier returns FAIL on produced records, the institution's certification language has nowhere to hide.

**Litigation use.** I would reference §10.12 and the §F table in any motion seeking the institution's verifier-run output. If the institution claims the records pass and refuses to disclose the verifier output, I would seek an order compelling production of the raw verifier output (stdout, exit code, the operational event `verifier.run_completed` per §A.4 ¶14) so the FAIL/PASS determination is on the record rather than in the institution's narrative.

### G-8 — The institution's customer-correlation index is unbinding and unaudited.

**Section/file:** `docs/customer-dispute-procedures.md` §"Documentation for dispute response"; §"Practical guidance for the dispute desk"

**Cross-examination angle.** The chain captures `(tenant_id, run_id, seq)` triples. The institution maintains a customer-correlation index that maps a customer dispute to the affected runs. That index is not in the chain, not signed by the HSM, not covered by the per-event MAC. It is the institution's internal lookup table.

When the bank tells me they pulled "all the runs affecting Mrs. Garcia," I want the index. I want the change history. I want to know whether the index was rebuilt after the dispute was filed. None of that lives in a chain-bound record I can authenticate. It lives in the institution's CRM or compliance datastore, and it is what the dispute desk uses to assemble the production.

If the index is wrong — accidentally or otherwise — the chain entries the defendant produces are integrity-bearing for what they are, but they may not be the entries the index should have surfaced. The defendant's witness will testify the chain is intact. That testimony is true. It is also irrelevant to the question of whether the right entries reached production. Chain integrity is a property of the entries that arrived; it says nothing about the entries that should have arrived but did not.

**Motion-in-limine angle.** I would seek the institution's customer-correlation index as part of discovery, including the index's audit log, its change history during the period at issue, and the procedural document governing how a dispute reference resolves to a run set. I would also seek an order naming the index as a load-bearing artifact for production completeness, with adverse-inference exposure if the institution cannot produce a contemporaneous, integrity-bound version.

## Litigation posture

Would I file a motion to exclude under Daubert/702? Not against the cryptography itself — it is FIPS-grounded, peer-reviewed, and has a real Daubert posture per §1.1. I would lose. I would file under 702 against expert testimony that interprets witness-mode verification or the §C hash-only certification as proof of substantive AI behavior — those are scope-limited authentications and the institution's expert will overreach if not constrained.

Would I challenge authenticity under FRE 901? Yes — the certification's foundation under 901(b)(9) and 902(13) on the G-1 independence theory, the G-2 verification-asymmetry theory, and the G-6 witness-mode-coverage theory. I would not challenge the cryptography. I would challenge whether the records reaching me are the records that should have reached me — a question the cryptography does not answer.

Would I seek an adverse-inference instruction under Rule 37(e)? Yes, on the G-4 spoliation theory — the institution's data-classification policy excluded substantive content from the chain during a period when class-member harms were foreseeable. That is the theory most likely to land, because the institution's witness will be forced to admit on cross that the policy was institutionally controlled and that the substantive content is unrecoverable. Rule 37(e)(2) sanctions require intent; institutional knowledge of foreseeable litigation plus a policy choice to exclude is a viable factual basis a bench could credit.

Would I depose the IT custodian? At length. The questions write themselves from the spec: who configured the verifier-run, who held the IKM, who authored the customer-correlation index, who decided which entries were responsive, who authored the partial-disclosure bundle, and how many other entries the seal-day covered that were not produced. Some of the answers will reveal that the same handful of people held all the relevant levers — which is what I need for the trustworthiness-clause challenge under 803(6) and the foundation challenge under 902(13).

## What the defense is going to use against me

The harder honesty. The verifier exit-code contract per §10.12, the normative failure-reason strings per §7, the test-vector corpus, and the open-source reference implementation are evidence-of-process artifacts I cannot match in most legacy bank litigation. The institution's IT witness gets to walk the bench through a 12-step procedure with byte-exact failure modes and a signed FIPS attestation. The cryptography has a real Daubert posture, the test corpus is reproducible, and §1.2 names what the chain proves and what it does not. That epistemic transparency reads well to a bench.

My case lives in the seams — independence of the verifier-run, asymmetry of verification access, defendant control over what is produced, and policy-level exclusion of substantive content. The cryptography passes Daubert. The institutional procedures around it are where I have angles. The bank's counsel will read this same document and prepare witnesses to close those angles. That is the work, on both sides.
