# AI safety and evaluation overlay

> **What this doc is.** The research-facing companion to the chain-of-custody spec. The chain was built for FFIEC compliance, but the artifacts it produces — tamper-evident captures of model inputs, outputs, sampling parameters, tool calls, and routing decisions — are exactly the substrate empirical alignment work needs. This overlay names how a safety researcher consumes the chain, what extensions to the institution's `audit.*` namespace make the chain a first-class research input, and where the line falls between FFIEC-compliance fields and research-only fields. The overlay is institution-side; the spec body is unchanged. An institution that wants to publish or contribute chains to consortium safety research opts into the patterns here.

> **What this doc is not.** A model-evaluation framework. The chain produces the evidence; the researcher's HELM, BIG-bench, ARENA, or institution-internal harness consumes it. This document defines the consumption interface and the data-quality posture, not the experimental methodology.

---

## 1. The chain as a research substrate

The chain was designed for compliance-grade integrity: a regulator with no access to the institution beyond a public key can confirm an AI decision was captured at the time it claims, with the inputs and parameters it claims, and that nothing was altered after capture. That property is also what empirical alignment work has been missing.

Production deployment logs are usually unsigned, redacted at the vendor's discretion, and aggregated at a granularity that loses the per-decision detail a research question needs. A chain entry, by contrast:

- Binds the model identifier, request and response side, into the canonical bytes the MAC covers.
- Binds the sampling parameters (temperature, top_p, seed, max_tokens) where the institution has configured `gen_ai_parameters` capture.
- Binds the institution's `audit.*` payload, which can carry decision class, customer-impact tier, fairness cohort indicators, and any institution-specific labels.
- Links to parent and child entries via `parent_run_id` / `parent_seq`, so multi-turn conversations and chained tool calls reconstruct cleanly.
- Sits inside a daily Merkle root signed by an HSM, so a researcher who pulls a chain extract can prove no entries were inserted, deleted, or reordered after capture.

The first three properties give the researcher a per-decision substrate. The fourth gives them conversation reconstruction. The fifth gives them defensible publication ground — a peer reviewer can ask "how do you know the institution did not edit the data before you saw it" and the answer is the integrity proof, not the institution's word.

The four research questions this overlay anchors:

1. **Behavioral drift.** Same model, same input, does output change over weeks?
2. **Red-team episode reconstruction.** Given a known failure mode, where else does it occur in production?
3. **Fairness audit replay.** Given a decision distribution, does it differ across demographic cohorts?
4. **Hallucination ground-truth labeling.** Can production fact-checks become labeled training data for hallucination classifiers?

Every section below is grounded in one of those four questions plus the supporting infrastructure (multimodal capture, embedding storage, evaluation-harness export, IRB compliance, differential privacy on aggregates).

---

## 2. Reproducibility — what the chain captures and what the institution decides

The chain's reproducibility posture is a per-deployment institution decision, not a spec mandate. The spec defines the shape (`gen_ai_parameters` is JCS-canonical JSON of model sampling and reproducibility parameters per spec §4.4); the institution decides what to put in it.

### 2.1 The recommended `gen_ai_parameters` schema for research consumption

Institutions that want their chains usable for drift-detection research populate `gen_ai_parameters` with the following fields. All are RECOMMENDED for research consumers; absence is operationally acceptable but limits the research questions that chain can answer.

| Field | Type | Notes |
|---|---|---|
| `temperature` | float | Sampling temperature. Required for stochastic-output reproducibility. |
| `top_p` | float | Nucleus-sampling cumulative-probability cutoff. Required when the deployment uses nucleus sampling. |
| `top_k` | int | Top-k sampling cutoff. Required when the deployment uses top-k sampling. |
| `seed` | string \| int | Vendor-specific seed representation. The verifier does not interpret the seed; it preserves the exact bytes the institution captured. |
| `max_tokens` | int | Token budget the request authorized. |
| `stop_sequences` | string[] | Stop-string list. |
| `presence_penalty` | float | Vendor-specific penalty. |
| `frequency_penalty` | float | Vendor-specific penalty. |
| `repetition_penalty` | float | Vendor-specific penalty (some open-weight models). |
| `sampler_impl` | string | Implementation identifier (e.g., `openai-python-sdk-1.3.5`, `anthropic-sdk-typescript-0.27.0`). Cross-version reproducibility hinges on this — the same seed under different SDK versions is not guaranteed to produce the same output. |
| `system_prompt_hash` | string | Lowercase hex SHA-256 of the system prompt. The full prompt may be too large to bind in every entry; the hash plus an institution-side prompt-version registry resolves the prompt content. |
| `system_prompt_version` | string | Institution-defined prompt version identifier (e.g., `prod-v3.4.1`). Pairs with `system_prompt_hash`. |
| `few_shot_examples_hash` | string \| null | Lowercase hex SHA-256 of the few-shot examples set. Null if no few-shot examples were supplied. |
| `retrieval_context_hashes` | string[] | Per-document SHA-256 hashes of any retrieved (RAG) context, in the order the context was injected. |

**Reproducibility consequence of omission.** An institution that omits `seed` and `temperature` accepts that chain entries from that deployment cannot be reproduced — a researcher re-running the same prompt on the same model gets different outputs due to stochasticity. This is operationally acceptable when reproducibility is not a deployment requirement; it limits the research questions the chain can answer. The institution names its reproducibility posture in its MRM policy (per `regulator-pack/ai-policy-alignment.md`): **"Reproducibility parameters captured — drift analysis enabled"** or **"Reproducibility not required for this deployment; seed and temperature omitted."** Research consumers filter their input data to include only deployments with reproducibility parameters.

### 2.2 Vendor-specific extension fields

Vendors emit parameters not on the recommended list above (OpenAI's `logit_bias`, Anthropic's `sampling_preset`, Google's `safety_settings`). The institution captures them as additional keys in `gen_ai_parameters`; the verifier does not interpret them but preserves the bytes. A research consumer aggregating across institutions sees a heterogeneous schema and must adapt — but the per-institution chain is internally consistent, and the canonical-bytes binding means a chain entry's parameters cannot drift after capture.

---

## 3. Behavioral drift detection

The core safety question: given the same model and the same input, does the output change over weeks? The chain answers it by giving the researcher a per-decision substrate they can re-run.

### 3.1 The drift-detection workflow

A researcher with access to chain entries from `(model_id, deployment_id)` over a multi-week window:

1. Extract all entries for the (model_id, deployment_id) pair across the window.
2. Group by week (or rolling sliding window).
3. For each week, compute the empirical distribution of decisions across input cohorts (input similarity bucket, decision class, customer-impact tier).
4. Run a divergence test (KL, JS, or per-cohort chi-squared) week-over-week.
5. Flag weeks where the divergence exceeds a researcher-determined threshold.
6. For flagged weeks, sample entries from the flagged cohort and re-run the model under the captured `gen_ai_parameters` to measure same-input output change.

The chain provides the canonical inputs and outputs; the researcher provides the divergence test and the re-run infrastructure.

### 3.2 Silent vendor-side model swaps

A hosted-model deployment can experience a silent model swap when the vendor reroutes `gpt-4` from one snapshot version to another behind a stable model identifier. The chain captures `gen_ai.response.model` (the vendor's claimed model on the response), but if the vendor reports the same identifier for both snapshots, the chain alone does not surface the swap.

Researchers detect swaps statistically: token-count distribution shifts, latency distribution shifts, decision-distribution divergence on a held-out replayable corpus. The chain supports this when the institution emits:

- `gen_ai.response.model` — the vendor's claimed model.
- `gen_ai.response.usage.input_tokens` — actual request token count (not max_tokens budget).
- `gen_ai.response.usage.output_tokens` — actual response token count.
- `gen_ai.response.latency_ms` — wall-clock latency.
- `audit.deployment.intent` — set to `vendor_reroute_observed` when the institution's own monitoring flags a probable swap (per spec §4.4.2).

Institutions running hosted-model deployments are encouraged to capture all five fields. The token-count and latency fields are the minimum substrate for swap detection; without them a researcher's analysis is blind to the dimension where swaps most often surface.

---

## 4. Red-team episode reconstruction

A red-teamer who finds a failure mode (jailbreak prompt, prompt-injection signature, harmful-output trigger) wants to know whether the failure occurred in production and, if so, in what other forms.

### 4.1 The reconstruction workflow

The chain does not support similarity search natively. Similarity search is a research-layer tool the institution provides separately when it wants to support red-team work:

1. The researcher identifies a failure-mode entry (from their own red-team session) — captured prompt, captured output.
2. The researcher computes embeddings of the failure-mode prompt and output using a reference embedding model (e.g., `all-MiniLM-L6-v2`, `text-embedding-3-large`).
3. The institution maintains an embedding index over the chain's entries. The index is computed offline, stored separately from the chain, and not bound under the chain MAC. Institutions interested in supporting red-team research opt into maintaining the index.
4. The researcher queries the index for nearest neighbors under cosine similarity. The institution returns matching `(tenant_id, run_id, seq)` triples plus de-identified entry content.
5. The researcher analyzes the returned set for common patterns: input distribution, conversation context, token counts, model state, deployment intent.

### 4.2 Optional embedding capture in the chain

Some institutions prefer to bind embeddings under the chain MAC for tamper-evidence on the embedding itself. Optional fields under `audit.embedding.*`:

- `audit.embedding.input_embedding_model` — the embedding model identifier.
- `audit.embedding.input_embedding` — base64-encoded float32 array OR a content-addressed reference to an external embedding store.
- `audit.embedding.output_embedding_model` — the output-side embedding model identifier.
- `audit.embedding.output_embedding` — base64-encoded float32 array OR reference.

These fields are NOT load-bearing for the chain's integrity claim — embeddings are derived from the entry text the chain already binds. They are informational. An institution that emits them lets the researcher use the chain entry directly as a similarity-search input without a separate index lookup. Institutions that prefer to keep chain entries small omit the embedding fields and maintain the offline index.

### 4.3 Prompt-injection forensic recovery

Prompt-injection is outside the chain's threat model — the chain proves the prompt the model received, not whether the prompt was hostile. But the chain's record of the prompt is the substrate for post-hoc detection.

A researcher suspecting an injection occurred at time T on customer C:

1. Pulls the chain entry for the suspected decision.
2. Confirms the entry's MAC and seal-membership pass per spec §7.
3. Extracts the captured prompt text. The prompt's bytes are what the model received; if the prompt contains an injection signature (instruction override, prompt-boundary marker, token-jailbreak), the signature is in the captured bytes, not a sanitized version.
4. Runs an injection classifier (or manual analysis) on the captured prompt.
5. Cross-references with the decision outcome: did the injection-flagged prompt produce an unusual decision relative to non-flagged peers in the same customer cohort?

The chain's contribution is bytes the researcher can trust. The injection-detection logic is the institution's safety tooling or the researcher's own model — both downstream of the chain.

---

## 5. Hallucination ground-truth labeling

Production fact-checks, when chained as children of the AI's original output, become labeled training data for hallucination classifiers. The pattern is already in the customer-dispute procedure (`docs/customer-dispute-procedures.md` §"Hallucination cross-check") — institutions emit `chain_kind = 'audit'` entries with `audit.fact_verification.*` payload, parent-linked to the disputed decision.

### 5.1 The labeling workflow

For research consumption:

1. Extract chain entries with `chain_kind = 'model_call'` whose run is referenced as `parent_run_id` by a downstream `audit.fact_verification.*` entry.
2. For each, the parent entry is the AI's claim; the child entry is the institutional fact-check disposition.
3. Label the parent entry as one of `confirmed`, `contradicted`, `unverifiable` per the child's `audit.fact_verification.disposition`.
4. The labeled set is hallucination training data — model outputs the institution has fact-checked, with disposition labels grounded in the institution's authoritative records.

Tokenization (per `docs/privacy-by-design.md`) applies on export: customer-PII tokens remain tokens; the AI's reasoning text (the hallucination target) is the load-bearing research content and stays readable.

### 5.2 Post-hoc harm and toxicity classification

Institutions deploying content-moderation pipelines or fairness-monitoring systems often classify AI outputs after capture (toxicity score, bias category, sensitive-content tag). These classifications chain the same way fact-verification does, under `audit.output_classification.*`:

| Field | Type | Notes |
|---|---|---|
| `audit.output_classification.classifier` | string | Identifier of the classifier (e.g., `toxicity-v1.3`, `pii-leakage-detector-v0.2`). |
| `audit.output_classification.category` | string | The classifier's verdict (e.g., `toxic`, `safe`, `unsafe_bias`, `prompt_injection_suspected`). |
| `audit.output_classification.confidence` | float | Classifier confidence in [0.0, 1.0]. |
| `audit.output_classification.review_kind` | enum | `automated` \| `human` \| `human_reviewed_automated`. |
| `audit.output_classification.parent_run_id` | string | The classified entry's run_id (chain-link via `parent_run_id` / `parent_seq` per spec §4.4). |

Researchers consuming chains with these labels train classifiers on the institution's actual production decisions, not synthetic data. The chain integrity-binds the classification, so a downstream consumer cannot dispute "you re-labeled the data after the fact" — the seal date pins the classification's recording time.

---

## 6. Fairness audit replay

The fairness question — "do decisions distribute differently across demographic cohorts" — collides with universal redaction. The chain's tokenization removes direct identifiers before canonicalization, and most institutions also redact quasi-identifiers (ZIP code, age band, gender) under their privacy posture. A researcher cannot re-bin chain entries by demographic cohort directly.

The pattern that resolves the tension:

### 6.1 The institution-side cohort mapping

The institution maintains a separate, governance-controlled mapping table — **not in the chain** — that links customer tokens to demographic attributes used for fairness audits. The mapping is institution-side; it is not exposed to the regulator's verifier and is not part of the chain's integrity claim.

The mapping's role:

- Inputs: customer token (the same token bound into chain entries via `audit.*` payload).
- Outputs: demographic cohort indicators (age band, gender, ZIP-code-derived geography, ethnicity if the institution collects it, any other fairness-sensitive attribute).
- Access: the institution's fairness committee, internal model-risk-management team, and (under specific governance) the institution's external research partners.
- Retention: tied to the chain's retention period; deleted on the same schedule.

The mapping enables the query "show me decision distribution for customer tokens in demographic cohort C." The chain stays redacted; the mapping is the join key the institution holds outside the chain.

### 6.2 The fairness audit workflow

1. The researcher (institution-internal or external partner) defines the cohort and the test (logistic regression on decision outcomes; disparate-impact ratio; demographic-parity metric).
2. The institution's fairness committee resolves customer tokens for the cohort.
3. The chain is queried for entries belonging to the resolved tokens.
4. The aggregate is computed; per-cohort decision distributions are surfaced.
5. Statistical tests are applied; conclusions name the cohort definition, the chain coverage, and the observed disparity.

The mapping is governed by the institution's DPIA and RoPA per `docs/privacy-by-design.md`. External researchers do not see the mapping; they receive aggregate statistics only. The institution's fairness-audit governance documents what cohorts are sanctioned for analysis and on what basis.

### 6.3 Differential privacy for published aggregates

When fairness statistics are published — in a research paper, a regulator filing, or a public-disclosure report — they may be a privacy disclosure under differential-privacy theory. The publication options:

1. **Aggregate before publication.** Compute the statistic, add Laplace noise proportional to the query's sensitivity, then publish the noised aggregate. The DP budget is the institution's privacy posture's responsibility; the chain's role ends at producing the substrate.
2. **Use post-retention de-identification.** After the institution's retention period (typically 7 years), the privacy-store mapping is deleted and customer tokens become irreversibly anonymized. Research published on post-retention chains has no live re-identification path; the DP analysis is structurally simpler.
3. **Federated analysis without raw chain transfer.** Researchers run queries against an institution-hosted secure-computation service that returns noised aggregates. The raw chain never leaves the institution.

The institution's research-data-sharing policy names which option applies for which research product. Pattern (3) is the strongest privacy posture and the one consortium research is moving toward.

---

## 7. Multi-turn and tool-call reconstruction

Modern AI agents make multiple model calls, invoke external tools, retrieve from databases, and weave responses across turns. The chain captures each step as a separate entry; reconstruction is mechanical.

### 7.1 The reconstruction procedure

Given a `run_id` representing an agent flow:

1. Pull all entries with the matching `run_id` ordered by `seq`.
2. Each entry's `chain_kind` names what it represents: `model_call`, `tool_call`, `routing`, `translation`, `audit`.
3. Each entry's `parent_run_id` / `parent_seq` (when present) names which prior entry produced this entry's input.
4. Topological sort over the parent-child graph reconstructs the flow.

A typical multi-turn agent run produces, in order: a `model_call` for the agent's reasoning step, a `tool_call` for the database/web/calculator invocation, a `model_call` for the agent's processing of the tool result, possibly more `tool_call` and `model_call` entries, and a final `model_call` for the agent's response. The DAG over these entries IS the agent's reasoning trace, reconstructed from chain bytes alone.

### 7.2 Tool-call schema for safety analysis

Tool-call entries (`chain_kind = 'tool_call'`) carry the institution's `audit.tool.*` payload. RECOMMENDED schema for research consumption:

| Field | Type | Notes |
|---|---|---|
| `audit.tool.name` | string | Tool identifier (e.g., `database_query`, `web_search`, `calculator`, `email_send`). |
| `audit.tool.input` | object \| string | The tool's input. PII tokenized per privacy-by-design. |
| `audit.tool.output` | object \| string | The tool's output. Truncated or summarized when >10 KB; the truncation strategy is institution-defined and named in the entry. |
| `audit.tool.latency_ms` | int | Wall-clock latency of the tool invocation. |
| `audit.tool.status` | enum | `success` \| `timeout` \| `error`. |
| `audit.tool.error_message` | string \| null | Error text on failure; null on success. |
| `audit.tool.permission_scope` | string \| null | The institution's authorization scope the agent invoked the tool under (e.g., `read-only`, `customer-data-mutation`, `external-api`). |

Full input and output capture is a safety requirement, not just a research convenience. An agent that calls tools without the chain capturing what it called and what it received is an unauditable agent. An institution that runs agentic AI in regulated workflows captures the tool-call full input/output as part of its standard `audit.tool.*` payload; absence is a control gap, not just a research data gap.

### 7.3 Reasoning-trace capture

Modern reasoning models (Claude with extended thinking, OpenAI o-series, similar) emit internal reasoning traces sometimes 10K+ tokens long. Three institution postures for these traces:

1. **Full capture.** The entire reasoning trace lives in `audit.reasoning_trace` on the model_call entry. Storage-intensive (a single entry can carry 50KB+ of reasoning text); maximally useful for interpretability research.
2. **Summary capture.** A tokenized or compressed summary lives in `audit.reasoning_trace_summary`. Smaller; preserves the reasoning's structural shape and key decision points; loses the verbatim text.
3. **No capture.** Reasoning traces are not logged. Smallest storage; preserves vendor-confidentiality where the trace contains proprietary chain-of-thought; eliminates interpretability research on those entries.

Institutions running interpretability or alignment research opt for (1) or (2). Institutions prioritizing efficiency or vendor-confidentiality opt for (3). The institution's CC8.1 control description names the choice; research consumers filter to the entries that carry the substrate they need.

---

## 8. Multimodal AI representation

Vision models (Claude with vision, GPT-4V, Gemini multimodal), audio models (Anthropic audio, Whisper, similar), and video models bring inputs the chain's text-centric default does not directly accommodate. The chain handles them by representing the media via hash, reference, or embedding rather than raw bytes.

### 8.1 The three representation patterns

Choose the pattern that fits the deployment's storage and research-access posture:

1. **Media hash.** SHA-256 of the media bytes is bound in `audit.modality.input_hash` (or output_hash for generated media). The hash binds the modality content under the chain MAC; the bytes themselves live in the institution's content-addressable storage. Suitable for media >100KB where storing inline would bloat the chain.
2. **Media reference.** A storage reference (`s3://bucket/object_id`, `gs://bucket/object`, content-addressable URI) is bound in `audit.modality.input_uri`. The institution preserves the referenced object for the chain's retention period. Suitable for media stored in the institution's existing object store.
3. **Media embedding.** A vector representation in `audit.modality.input_embedding` (base64-encoded float32 array) lets researchers do similarity search on the media without retrieving the raw bytes. Suitable when the media itself is sensitive (customer voice recordings, identity documents) but the embedding is research-acceptable.

### 8.2 The decision output

The decision output (text) is captured normally on the same entry. The integrity binding covers the text output; the modality fields are informational on the same entry. A research consumer asking "what visual feature drove this decision" pulls the entry's `audit.modality.input_embedding` (pattern 3) or `audit.modality.input_hash` (pattern 1, then retrieves the bytes from content-addressable storage).

The institution's CC8.1 names which pattern applies for which deployment. A multimodal customer-service deployment might use pattern 2 (S3 references) for input audio, pattern 1 (hash) for output text, and pattern 3 (embedding) when the institution wants to support similarity search without exposing the raw audio.

---

## 9. Evaluation harness integration

Production chain entries can feed academic evaluation harnesses (HELM, BIG-bench, ARENA) when the institution exports them in a compatible format. The export is institution-side; the chain provides the substrate.

### 9.1 The minimum export schema

Per chain entry, the harness expects:

| Harness field | Chain source |
|---|---|
| `input` | The captured prompt text from the model_call entry (under the institution's `audit.*` namespace, typically `audit.model_call.prompt` or equivalent). |
| `output` | The captured response text (under `gen_ai.response.*` or institution-equivalent). |
| `model` | `gen_ai.response.model`. |
| `model_version` | The institution's model-version identifier from `audit.model_call.model_id` (per Q-5 below). |
| `task_tag` | An institution-defined task identifier (`mortgage_underwriting`, `fraud_triage`, `customer_chat`); the institution adds this to its `audit.*` namespace. |
| `ground_truth` | When available from the institution's downstream processing (a fact-verification disposition, an underwriter override, a customer outcome) — pulled from chained child entries. |

The institution publishes the export schema with its research-data-sharing policy. A researcher running HELM against the export sees how the institution's production performance compares to the harness's benchmark — answering "is the production deployment performing at the level the model card claims" with statistical rigor.

### 9.2 Model-card / evaluation-card linkage

Optional fields binding chain entries to published model documentation:

- `gen_ai.model_card_url` — the model card's published URL at the time of decision.
- `gen_ai.model_card_version` — the model-card version identifier.
- `gen_ai.evaluation_card_url` — link to the model's evaluation card or technical report.

A researcher can cross-reference: "the model card claims 92% accuracy on task Y; production chain data shows 91.3% — is there a gap?" The linkage is informational; the chain's integrity does not depend on it. Institutions consuming hosted models from providers that publish model cards (Anthropic, OpenAI, Google, Mistral) benefit from emitting these fields.

---

## 10. Out-of-distribution and capability evaluation

Chain attributes support capability and OOD detection at three layers.

### 10.1 OOD detection from `audit.*` payload

The institution's `audit.*` payload typically carries decision-class indicators (decision type, customer-impact tier, deployment intent). A research consumer can compute the joint distribution of `(audit.decision.class, gen_ai_parameters.temperature, audit.deployment.intent)` and flag entries that fall in low-density regions of the joint distribution. Those entries are candidate OOD events: edge cases the production deployment encounters rarely.

The chain provides the substrate. The OOD detector is the researcher's tool. The integrity-binding ensures the rare events haven't been redacted post-hoc to fit a narrative.

### 10.2 Capability-evaluation chains

When an institution runs internal capability evaluations (jailbreak resistance, mathematical reasoning, factual recall) against deployed models, the evaluation events themselves chain. The pattern:

- `chain_kind = 'model_call'` for the evaluation-prompt invocation.
- `audit.evaluation.benchmark_id` naming the evaluation suite (`mmlu_subset`, `internal_jailbreak_corpus_v3`).
- `audit.evaluation.expected_response` carrying the ground-truth answer (when the benchmark has one).
- `audit.evaluation.score` carrying the per-item score the institution computed.
- `audit.evaluation.aggregate_run_id` linking to a parent `chain_kind = 'audit'` entry that records the full evaluation-run summary.

The aggregate entry rolls up the per-item scores into the evaluation run's summary statistics, integrity-bound under the daily seal. A researcher analyzing capability drift over time consumes the aggregate entries; a researcher analyzing per-item failure patterns consumes the per-item entries.

### 10.3 Long-tail behavior detection

Most failure modes occur once per 100K decisions or rarer. If the institution's chain uses sampling (capturing every Nth event), rare events disappear. The institution's research-friendly posture is one of:

- **Full capture.** Every event lands in the chain. Storage-heavy; preserves all rare behavior.
- **Failure-stratified oversampling.** 100% capture when `audit.decision.class = 'deny'` or `audit.deployment.intent = 'vendor_reroute_observed'`; 10% capture when the decision is the modal `'approve'`. Rare events stay at full coverage; the institution accepts the bias-correction work researchers must do on the modal class.
- **Uniform sampling.** 10% capture across all events; documented as such. Researchers know to expect the sampling and adjust statistical claims.

The institution names its sampling policy in CC8.1. Research papers analyzing the institution's chain MUST state the sampling regime — "All decisions analyzed (no sampling)" or "Oversampled failure cases by 10x" — so peer reviewers can adjust for selection bias.

---

## 11. Statistical hypothesis testing on aggregated chains

A research claim like "Model X is 2.3% more accurate than Model Y on task Z, with p<0.05" requires N samples per (model, task) cell. The chain provides the per-decision substrate; the researcher provides the aggregation and the statistical machinery.

### 11.1 What the chain guarantees and what it does not

The chain guarantees the per-entry contents are byte-faithful to what the institution captured at decision time. The chain does NOT guarantee the institution's sample is representative, uniform, or unbiased relative to the population of decisions the institution made.

A research consumer pooling chains from multiple institutions MUST disclose each institution's sampling policy in the paper's data-availability statement. Chains from Institution A (uniform sampling) cannot be silently pooled with chains from Institution B (on-demand-only sampling for disputes); the resulting aggregate is a biased sample of the union, and the paper's statistical claims are invalid without a disclosed correction.

### 11.2 Cross-institution federation

A multi-institution research collaboration on chain data requires:

1. **A common schema.** The institutions agree on the required `gen_ai_parameters` fields, the required `audit.*` namespace fields, and the export format. This document's §2 is a candidate baseline.
2. **A common sampling policy.** Either all institutions sample uniformly, or each declares its policy and the collaboration's analysis adjusts.
3. **A common privacy posture.** Tokenization is consistent across institutions; cohort mappings are not exchanged; aggregates are noised before sharing per §6.3.
4. **A common verification stack.** Each institution runs the §7 witness-verifier on its export; the consortium accepts only PASS or PASS-STRUCTURALLY exports.

The consortium's research-data-sharing agreement codifies the four points. The chain's integrity property survives federation because each institution's seal record carries that institution's HSM signature; pooling chains from N institutions produces an aggregate where each entry's authenticity remains independently verifiable against its origin institution's public key.

---

## 12. IRB compliance and de-identification posture

Research using human-subjects data requires IRB approval. The chain's tokenization (per `docs/privacy-by-design.md`) supports IRB-approvable de-identification when the institution operates the post-retention erasure discipline.

### 12.1 The de-identification stages

| Stage | What is in the chain | What is in the privacy-store |
|---|---|---|
| Live (within retention window) | Tokenized customer references; AI inputs and outputs; institution `audit.*` payload. | Mapping table linking tokens to direct identifiers; demographic cohort mapping. |
| Post-retention | Tokenized references whose mapping is deleted; the AI inputs and outputs that survived per the institution's retention discipline. | Empty (mappings deleted; tokens become irreversibly anonymous). |

NIST SP 800-188 names de-identification as reducing identifiability through removal or transformation of direct identifiers and quasi-identifiers. The live-window chain meets the direct-identifier removal half (tokens replace identifiers); the quasi-identifier handling is the institution's privacy posture's responsibility. The post-retention chain meets both halves for any researcher who arrives after the retention period.

### 12.2 IRB approval pathways

For research on chain data:

1. **Post-retention research.** Once the retention period closes and the privacy-store mapping is deleted, tokens are irreversibly anonymous and the chain is automatically IRB-exempt under the minimal-risk anonymous-data exception. Most longitudinal studies on chain data fall here.
2. **In-retention research with IRB approval.** A researcher analyzing live-window chain data obtains IRB approval per the institution's host-IRB process. The approval names the cohort definition, the de-identification posture, the data-handling controls, and the publication-review path.
3. **In-retention research with IRB waiver.** Minimal-risk research with no direct contact with subjects may qualify for an IRB waiver under 45 CFR 46.116(f). The institution's IRB office reviews; the approval is documented.

The institution's research-data-sharing policy names which pathway applies to which research engagement. The chain's role is to produce an integrity-bound substrate; the IRB approves the research design that consumes it.

---

## 13. Provider attestation validation

Some vendors (Anthropic, OpenAI in select deployments) provide cryptographic attestations that the model produced the response. The spec captures the attestation byte-for-byte under `gen_ai.provider_attestation` (per spec §4.4), but the chain itself does not validate it — validation is institution-side or research-consumer-side.

### 13.1 The two-layer integrity claim

When the chain captures an attestation:

- **Layer 1: chain integrity.** The chain proves the institution captured the attestation as recorded. Spec §7's verifier output is the foundation.
- **Layer 2: vendor-attestation integrity.** A researcher with the vendor's public key validates the attestation, confirming the vendor signed the response the institution captured.

The two layers compose: Layer 1 says the institution's record is authentic; Layer 2 says the vendor produced what the institution recorded. Together, they close the trust path from "the model said this" to "the institution preserved the model's response without alteration."

### 13.2 The research-consumer validation procedure

A researcher consuming chains with attestations:

1. Obtains the vendor's public key from the vendor's documentation (Anthropic publishes its attestation key; OpenAI varies by deployment).
2. Obtains the vendor's validation procedure (signature algorithm, canonical-form definition, key rotation policy).
3. Runs Layer 1 verification per spec §7 — the chain entry passes.
4. Runs Layer 2 verification per the vendor's procedure — the attestation matches.
5. Cross-references the institution's `gen_ai.response.model` against the vendor's attested model (the vendor's attestation may name a more specific snapshot than the institution's request side).

Discrepancies in step 5 (institution claimed `gpt-4`; vendor attested `gpt-4-0613`) are a research signal — they expose silent vendor-side routing. The chain alone cannot surface the discrepancy; the chain plus vendor attestation does.

---

## 14. The research-publication checklist

A research paper analyzing chain data should include a data-availability statement covering:

1. **Source institution(s)** — named or anonymized per the data-sharing agreement.
2. **Sampling regime** — full capture, failure-stratified oversampling, or uniform sampling, with the relevant ratio.
3. **Reproducibility posture** — whether the source chains carried `gen_ai_parameters` for re-run reproducibility, or were captured without seed and temperature.
4. **De-identification stage** — live-window with IRB approval, live-window with IRB waiver, or post-retention.
5. **Verification posture** — the §7 verifier output the institution provided (PASS, PASS-STRUCTURALLY, or PASS-WITH-ANOMALY) on the analyzed corpus.
6. **Privacy controls** — differential-privacy noising of published aggregates (with the DP epsilon parameter), federated-analysis posture, or post-retention de-identification.
7. **Cohort definition** — for fairness research, the demographic cohort dimensions and the institution's mapping governance.
8. **Provider-attestation validation** — whether vendor attestations were captured and validated per §13.

A paper missing one or more of these statements is not necessarily wrong, but the peer reviewer cannot evaluate the research's privacy posture, statistical validity, or integrity foundation. The statements are the chain's contribution to research credibility.

---

## 15. Cross-references

| Topic | Where the substance lives |
|---|---|
| Per-entry MAC and canonical bytes | spec §4.1, §5 |
| Daily Merkle seal and HSM signature | spec §4.2, §4.3 |
| `gen_ai_parameters` placement | spec §4.4 |
| `audit.*` namespace | spec §4.4 |
| Witness-verifier mode | spec §7 |
| Hallucination cross-check pattern | `docs/customer-dispute-procedures.md` §"Hallucination cross-check" |
| Tokenization and privacy-store | `docs/privacy-by-design.md` |
| MRM policy framework | `docs/regulator-pack/ai-policy-alignment.md` |
| Threat model for prompt injection | `docs/design/09-threat-model.md` §"Adversarial AI scenarios" |
| Retention and post-retention erasure | spec §10.9, §10.13 |
| Multi-tenant federation | spec §10.15 (Pattern A and Pattern B) |
| Selective-production sampling for research extracts | `docs/selective-production-and-sampling.md` |

---

## 16. The consortium-research path — what the chain enables at scale

A single institution's chain produces useful research data; a consortium of institutions sharing chains across cooperating banks produces orders-of-magnitude more useful data. The consortium pattern is the long-term direction empirical alignment work is heading, and the chain's integrity property makes the cooperation tractable.

### 16.1 The consortium structure

A research consortium of N institutions agrees on:

1. A common required-fields list extending §2.1's recommended `gen_ai_parameters` schema.
2. A common `audit.*` namespace for fairness-cohort indicators, decision-class tags, and customer-impact-tier markers.
3. A common privacy posture — tokenization is consistent across institutions; cohort mappings are not exchanged across the consortium boundary; published aggregates are noised.
4. A common verification stack — each institution runs the §7 witness-verifier on its export; the consortium accepts only PASS or PASS-STRUCTURALLY exports.
5. A common research-data-sharing agreement codifying access rights, publication pathways, IRB posture, and cost allocation.

Each institution preserves its IKM custody. Each institution publishes only its own public key (per spec §10.5). The consortium's research products consume the federated dataset; no IKM is ever transferred.

### 16.2 Federated query patterns

A consortium-level research question — "how does Model Y's decision distribution shift across cohort C in the period after a vendor-side rerouting event" — answers via federated query:

1. The researcher submits a query specification to each institution's research-access endpoint.
2. Each institution runs the query against its own chain (institution-local computation).
3. Each institution returns aggregated, optionally noised results — never raw chain entries.
4. The researcher aggregates the per-institution results into a consortium-level finding.

The chain's integrity-binding survives the federation. Each per-institution aggregate is verifiable against that institution's seal records; the consortium-level finding is grounded in N independent integrity-bound substrates.

### 16.3 Why consortium scale matters for safety research

Single-institution chains cover the deployments and customer cohorts of one bank. Consortium chains cover the deployments and customer cohorts of many banks across the relevant market. For safety questions that turn on rare-event detection — the 1-in-100K failure mode, the silent vendor-side rerouting that reaches multiple consumers, the prompt-injection pattern that crosses institutional boundaries — single-institution chains are statistically underpowered. Consortium chains are not.

The chain's integrity-bound substrate is the reason the consortium is feasible. Without integrity-bound capture, the consortium's research products would be vulnerable to "Institution X edited their data to look better" challenges; with integrity-bound capture, those challenges are mechanically refutable. The chain's contribution to the consortium is exactly the property the consortium needs.

---

## 17. Operational notes for institutions opting into research consumption

The optional fields in this document add capture overhead — entry size, storage, and processing cost. Institutions opt in selectively per deployment based on the deployment's research-utility posture.

### 17.1 Per-deployment opt-in matrix

| Deployment class | `gen_ai_parameters` reproducibility fields | Embedding capture | Reasoning-trace capture | Multimodal embedding | Output classification |
|---|---|---|---|---|---|
| High-stakes credit-decision (ECOA / fair-lending) | RECOMMENDED | OPTIONAL | OPTIONAL | OPTIONAL | RECOMMENDED |
| Fraud-flagging (BSA / AML) | RECOMMENDED | OPTIONAL | OPTIONAL | OPTIONAL | RECOMMENDED |
| Customer-service routing | OPTIONAL | OPTIONAL | OPTIONAL | OPTIONAL | OPTIONAL |
| Internal task assignment / advisory | OPTIONAL | NO | OPTIONAL | OPTIONAL | OPTIONAL |
| Regulatory-narrative generation (BSA SAR drafts, adverse-action notices) | RECOMMENDED | OPTIONAL | RECOMMENDED | OPTIONAL | RECOMMENDED |
| Tool-using agentic deployments | RECOMMENDED | RECOMMENDED | RECOMMENDED (for interpretability) | OPTIONAL | RECOMMENDED |

The matrix is informative; the institution's MRM committee opines on the per-deployment opt-in pattern. RECOMMENDED settings reflect the deployments where research utility is highest and the regulatory scrutiny is deepest; OPTIONAL settings reflect deployments where the chain's compliance role dominates and research utility is secondary.

### 17.2 Storage and retention implications

Adding embedding capture to a deployment that runs 10M decisions per year produces roughly 10 TB of additional embedding data over the chain's 7-year retention period (assuming 1024-dim float32 embeddings). Reasoning-trace capture is substantially heavier — a single Claude-with-extended-thinking entry can carry 50KB+ of reasoning text, producing roughly 350 GB of additional data per year per deployment.

Institutions opting into research consumption budget for the additional storage. The institution's CC8.1 control description names the per-deployment storage footprint and the retention discipline. Storage minimization at the embedding layer is OK — institutions can retain embeddings for a shorter period than the chain itself, accepting that embedding-based queries answer only the more-recent window. The chain entries themselves are retained per spec §10.13; the optional research-utility fields can have shorter retention, named in the institution's MRM policy.

### 17.3 Privacy posture for opt-in fields

Embedding capture introduces a re-identification surface. An attacker with access to embeddings and an embedding model can sometimes recover content that was tokenized in the chain entry's `audit.*` payload. The institution's privacy assessment for opt-in fields:

- **Embeddings of input text containing tokenized PII.** The embedding may preserve enough signal to invert the tokenization if the embedding model is known. Mitigation: use a privacy-preserving embedding model (DP-trained), or restrict embedding access to within the institution's secure-computation boundary.
- **Reasoning-trace content.** May contain customer PII the AI mentioned in its reasoning. Mitigation: tokenize the reasoning trace using the same token vocabulary as the entry's `audit.*` payload, and treat the trace as PII for retention and access purposes.
- **Multimodal embeddings of customer media.** May preserve identifying features (face embeddings, voice embeddings). Mitigation: use task-specific embeddings (decision-relevant features only) rather than general-purpose embeddings, and treat them as PII.

The institution's DPIA covers each opt-in field's privacy posture; the institution's RoPA names the lawful basis for the additional processing. The opt-in is institutional, not customer-by-customer — but the institution's privacy notice MUST disclose the additional processing when the opt-in introduces re-identification surface beyond the baseline chain.

---

## 18. Stopping criterion for this overlay

The overlay covers the consumption interface for the four research questions named in §1, the multimodal and tool-call extensions, the IRB and DP postures for publication, the research-publication checklist, the consortium-research scaling path, and the per-deployment opt-in matrix. It does not name an experimental methodology — that is the researcher's framework. It does not enumerate every possible vendor parameter — vendors evolve and the institution captures what its deployment uses. It does not prescribe a specific embedding model — the researcher picks one suited to the question.

The overlay is institution-side. An institution that wants its chains usable for safety research opts into the patterns here; an institution that prioritizes minimal capture omits the optional fields and accepts that its chains support compliance but limit research consumption. Both are conformant to the spec; the spec scopes itself to integrity, not research utility.

The substrate is here. The research is the consumer's.
