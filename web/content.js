// content.js — common data file. Edit this to add/modify documents and stakeholder mappings.
// All text data lives here; the app templates render it.

const CATEGORIES = [
  { id: 'spec',       name: 'Specifications',          icon: 'mdi-file-document-multiple-outline', color: 'primary' },
  { id: 'design',     name: 'Design Rationale',        icon: 'mdi-pencil-ruler',                   color: 'primary' },
  { id: 'regulator',  name: 'Regulator Pack',          icon: 'mdi-shield-account-outline',         color: 'info' },
  { id: 'controls',   name: 'Control Map',             icon: 'mdi-checkbox-marked-circle-outline', color: 'info' },
  { id: 'soc',        name: 'SOC Pack',                icon: 'mdi-clipboard-check-outline',        color: 'info' },
  { id: 'audience',   name: 'Audience Summaries',      icon: 'mdi-account-group-outline',          color: 'success' },
  { id: 'operations', name: 'Operations',              icon: 'mdi-cog-outline',                    color: 'success' },
  { id: 'adoption',   name: 'Adoption & Onboarding',   icon: 'mdi-rocket-launch-outline',          color: 'success' },
  { id: 'incident',   name: 'Incident & Legal',        icon: 'mdi-alert-circle-outline',           color: 'warning' },
  { id: 'audit',      name: 'Audit Support',           icon: 'mdi-magnify-scan',                   color: 'warning' },
  { id: 'special',    name: 'Specialized Scenarios',   icon: 'mdi-shape-plus',                     color: 'secondary' }
];

const STAKEHOLDERS = [
  {
    id: 'all',
    name: 'All / Browse',
    icon: 'mdi-view-grid-outline',
    tagline: 'Browse the complete corpus, organized by category.',
    description: 'See every document grouped by topic. Use this view when you want a complete picture or are unsure which role to pick.',
    primaryDocs: []
  },
  {
    id: 'ceo',
    name: 'Bank Executive (CEO)',
    icon: 'mdi-domain',
    tagline: 'Plain-English overview; cost; examiner conversation.',
    description: 'You want to know what the chain does, what it costs, and what to expect when an examiner walks in.',
    primaryDocs: ['management-summary', 'cost-model', 'mrm-committee-brief', 'audit-committee-summary']
  },
  {
    id: 'audit-committee',
    name: 'Audit Committee Chair',
    icon: 'mdi-clipboard-account-outline',
    tagline: 'Oversight role; what to expect in committee materials.',
    description: 'You oversee bank controls. The chain is one such control. You read summaries, ask oversight questions, sign off on remediation plans.',
    primaryDocs: ['audit-committee-summary', 'management-summary', 'cost-model']
  },
  {
    id: 'mrm',
    name: 'MRM Committee Chair',
    icon: 'mdi-chart-bell-curve',
    tagline: 'SR 11-7 model lifecycle; chain as a model-risk control.',
    description: 'You evaluate whether the chain helps your validators do effective challenge under SR 11-7.',
    primaryDocs: ['mrm-committee-brief', 'customer-dispute-procedures', 'ai-policy-alignment']
  },
  {
    id: 'chain-ops',
    name: 'Chain Operations Team',
    icon: 'mdi-server-network',
    tagline: 'Runtime ops, IR, DR, deployment, at-scale guidance.',
    description: 'You run the ledger, configure the SDK, manage the HSM, and respond to chain-detected events.',
    primaryDocs: ['operator-guide', 'cloud-hsm-guide', 'byoc-deployment', 'incident-response-playbook', 'dr-and-resilience', 'at-scale-operations', 'first-engagement-guide', 'edge-and-federated-ai']
  },
  {
    id: 'adopter',
    name: 'Bank Chain Adopter',
    icon: 'mdi-rocket-launch-outline',
    tagline: 'Deciding whether and how to adopt; pilot ramp-up.',
    description: 'You are evaluating whether to adopt and planning the rollout.',
    primaryDocs: ['management-summary', 'cost-model', 'minimum-viable-deployment', 'design-overview', 'm-and-a-handoff']
  },
  {
    id: 'vendor-mgmt',
    name: 'Vendor Management',
    icon: 'mdi-handshake-outline',
    tagline: 'Vendor-hosted, BYOC, supply chain, M&A.',
    description: 'You oversee the bank\'s relationships with chain implementation vendors.',
    primaryDocs: ['vendor-hosted-controls', 'byoc-deployment', 'supply-chain', 'm-and-a-handoff']
  },
  {
    id: 'privacy',
    name: 'Privacy / GDPR Team',
    icon: 'mdi-account-lock-outline',
    tagline: 'Privacy-by-design, GDPR/CCPA, customer rights.',
    description: 'You ensure the chain composes with the bank\'s privacy program.',
    primaryDocs: ['privacy-by-design', 'tsc-mapping', 'customer-dispute-procedures']
  },
  {
    id: 'legal-ir',
    name: 'Legal / IR Team',
    icon: 'mdi-gavel',
    tagline: 'IR scenarios, court-ordered disclosure, disputes.',
    description: 'You handle the chain\'s involvement in legal proceedings, IR escalations, and customer disputes.',
    primaryDocs: ['incident-response-playbook', 'legal-disclosure', 'customer-dispute-procedures']
  },
  {
    id: 'internal-audit',
    name: 'Internal Audit',
    icon: 'mdi-magnify-scan',
    tagline: 'Independent verification, testing procedures.',
    description: 'You run independent assessments of the chain controls within the bank.',
    primaryDocs: ['audit-procedures', 'cuecs', 'tsc-mapping', 'portfolio-comparison-procedures']
  },
  {
    id: 'soc',
    name: 'SOC Engagement Team',
    icon: 'mdi-clipboard-check-outline',
    tagline: 'SOC 1 / SOC 2 attestation; criterion mapping.',
    description: 'You issue or rely on SOC opinions on chain implementations.',
    primaryDocs: ['soc-section-4', 'control-evidence-events', 'tsc-mapping', 'cuecs', 'audit-procedures', 'anomaly-template', 'vendor-hosted-controls', 'user-entity-summary']
  },
  {
    id: 'ffiec-it',
    name: 'FFIEC IT Examiner',
    icon: 'mdi-shield-account-outline',
    tagline: 'IT examination; Handbook alignment.',
    description: 'You conduct IT examinations and use the chain output as examination evidence.',
    primaryDocs: ['examiner-quickstart', 'examiner-training', 'sample-report', 'finding-language', 'handbook-mapping', 'deployment-package', 'examiner-approval', 'portfolio-comparison-procedures']
  },
  {
    id: 'ffiec-cyber',
    name: 'FFIEC Cybersecurity Examiner',
    icon: 'mdi-shield-bug-outline',
    tagline: 'NIST CSF, threat model, IR, supply chain.',
    description: 'You examine cybersecurity controls including the chain.',
    primaryDocs: ['csf-2', 'threat-model', 'incident-response-playbook', 'cloud-hsm-guide', 'supply-chain', 'ai-policy-alignment', 'edge-and-federated-ai']
  },
  {
    id: 'ffiec-eic',
    name: 'FFIEC Examiner-in-Charge',
    icon: 'mdi-account-tie-outline',
    tagline: 'Examination logistics; bank-management communication.',
    description: 'You lead the examination and own the bank-management conversation.',
    primaryDocs: ['sample-report', 'finding-language', 'examiner-approval', 'first-engagement-guide', 'portfolio-comparison-procedures', 'management-summary']
  },
  {
    id: 'ffiec-cfpb',
    name: 'CFPB / Consumer Protection',
    icon: 'mdi-account-heart-outline',
    tagline: 'Customer disputes; AI policy; fair-lending evidence.',
    description: 'You examine consumer-facing AI behavior and customer protections.',
    primaryDocs: ['customer-dispute-procedures', 'ai-policy-alignment', 'csf-2']
  },
  {
    id: 'cryptographer',
    name: 'Cryptographic Expert',
    icon: 'mdi-key-chain-variant',
    tagline: 'Threat model; HMAC chain; HSM custody.',
    description: 'You provide cryptographic-witness testimony or independent cryptographic review.',
    primaryDocs: ['threat-model', 'design-chain', 'design-merkle', 'design-hsm', 'spec-main']
  },
  {
    id: 'implementer',
    name: 'Implementer',
    icon: 'mdi-code-tags',
    tagline: 'Building or porting an SDK, ledger, or verifier.',
    description: 'You implement a conforming SDK, ledger server, or verifier.',
    primaryDocs: ['spec-main', 'design-overview', 'design-test-vectors', 'design-chain', 'design-merkle', 'design-hsm', 'design-otlp']
  }
];

const DOCUMENTS = [
  // Specifications
  {
    id: 'spec-main',
    title: 'Chain-of-Custody Specification v1.0',
    category: 'spec',
    path: '../spec/chain-of-custody-v1.md',
    summary: 'The normative specification. Defines the four primitives, conformance keywords, optional attributes, operational requirements, and stakeholder navigation.',
    analogy: 'Think of it as the building code for AI-decision integrity — every conforming implementation builds the same shape, even if from different bricks.',
    keyPoints: [
      'The four primitives: HMAC chain, daily Merkle seal, HSM signature, OTLP wire',
      'Mandatory attributes plus six optional v1.0-final attributes',
      'Session-key handshake security floor (auth, conf, no-cache)',
      'Operational requirements: reconciliation, append-only, time sync, HSM custody',
      'Stakeholder navigation at §13'
    ],
    lift: 'normative'
  },

  // Design rationale (10 docs)
  {
    id: 'design-overview',
    title: 'Design — System Overview',
    category: 'design',
    path: '../docs/design/00-overview.md',
    summary: 'The big-picture design: components, data flow, trust boundaries, deployment topologies.',
    analogy: 'The architectural blueprint. Where the rooms are; which walls hold the building up; where the lockable doors are.',
    keyPoints: [
      'Three trust zones: application, ledger, examiner',
      'Three deployment topologies: self-hosted, BYOC, vendor-hosted',
      'Eight-step data flow from agent to verifier',
      'Per-tier proportionality (community, mid-size, Tier-1)'
    ]
  },
  {
    id: 'design-primitives',
    title: 'Design — The Four Primitives',
    category: 'design',
    path: '../docs/design/01-primitives-spec.md',
    summary: 'Why four primitives, what each defends against, alternatives considered, and the auditor\'s-lens defense.',
    analogy: 'Like a security system with four locks, each guarding a different attacker. Remove any one, and the building has an open door.',
    keyPoints: [
      'HMAC chain defends against external attacker on the wire',
      'Daily Merkle seal defends against insider with database access',
      'HSM signature defends against insider with operational access',
      'OTLP wire prevents vendor lock-in for verification'
    ]
  },
  {
    id: 'design-chain',
    title: 'Design — Chain Construction (Hot Path)',
    category: 'design',
    path: '../docs/design/02-chain-construction.md',
    summary: 'The HMAC chain construction in detail: canonicalization, persistence ordering, handshake, performance budget, multi-process patterns.',
    analogy: 'The wax seal on a chain of letters. Each letter wax-seals a hash of the letter before it. Try to remove a letter — the seals stop fitting.',
    keyPoints: [
      'JCS canonicalization for byte-for-byte determinism',
      'Write-then-disclose persistence ordering for crash safety',
      'Three handshake-floor properties (auth, conf, no-cache)',
      'Per-process / shared-run / DAG multi-process patterns',
      'Hot-path budget under 2 ms p99'
    ]
  },
  {
    id: 'design-merkle',
    title: 'Design — Daily Merkle Seal',
    category: 'design',
    path: '../docs/design/03-merkle-seal.md',
    summary: 'RFC 6962 Merkle construction; deterministic ordering; cadence proportionality; time-stamp authority.',
    analogy: 'Like notarizing every transaction in a day with a single notary stamp that depends on every transaction. Change one — the stamp doesn\'t match.',
    keyPoints: [
      'RFC 6962 binary Merkle (Certificate Transparency precedent)',
      'Empty days, late-binding events, out-of-order handling',
      'Streaming Merkle for billions/day scale',
      'Cadence: hourly / daily / weekly with examiner approval',
      'Day boundary = ledger receive time, not host clock'
    ]
  },
  {
    id: 'design-hsm',
    title: 'Design — HSM Custody',
    category: 'design',
    path: '../docs/design/04-hsm-custody.md',
    summary: 'HSM signing, key management, cloud HSM provisioning, master-key rotation window, notification thresholds.',
    analogy: 'A vault with a stamp inside. The stamp signs documents but never leaves the vault. Even the bank president can\'t take the stamp out.',
    keyPoints: [
      'Ed25519 in FIPS 140-2 L3+ HSM',
      'Cloud HSM matrix (AWS / Azure / Google) with conformant tiers',
      'Master-key rotation window with multi-version seal records',
      'PIN rotation procedure (quarterly)',
      'Cosign + GPG dual-signing for verifier binary'
    ]
  },
  {
    id: 'design-otlp',
    title: 'Design — OTLP Wire',
    category: 'design',
    path: '../docs/design/05-otlp-wire.md',
    summary: 'OpenTelemetry Protocol wire format, attribute namespace, no-key-material guarantee, optional attributes, TLS requirement.',
    analogy: 'Like the FedEx tracking system. The package contents (chain attributes) ride along with the institution\'s normal observability traffic.',
    keyPoints: [
      'OTLP layered: protobuf wire + JCS canonical hash',
      'ffiec.chain.* namespace (8 required + 6 optional attributes)',
      'No key material in attributes — safe to route anywhere',
      'TLS 1.3 minimum (1.2 sunset 2028-01-01)',
      'Backend interop: SIEM, FinOps, observability'
    ]
  },
  {
    id: 'design-ledger',
    title: 'Design — Ledger Server',
    category: 'design',
    path: '../docs/design/06-ledger-server-design.md',
    summary: 'Ledger server architecture: pipeline, storage, daily seal subsystem, BYOC, operational events, DR.',
    analogy: 'The post office for chain events. Receives mail, validates the postmark, files in the permanent archive, and notarizes the day\'s incoming.',
    keyPoints: [
      'Single-binary, multi-component (receiver, pipeline, writer, seal pools)',
      'WAL → hot store → cold store storage tier',
      'Per-tenant signing with master_version recorded',
      'Operational events for SOC and examination evidence',
      'RPO/RTO targets and sync-replication trade-offs'
    ]
  },
  {
    id: 'design-verifier',
    title: 'Design — Verifier',
    category: 'design',
    path: '../docs/design/07-verifier-design.md',
    summary: 'Standalone offline verifier: deterministic output, working-paper bundle, supply chain, examiner deployment.',
    analogy: 'The lighthouse keeper checking the ship\'s log. Independent, offline, just one binary; produces a report the keeper signs.',
    keyPoints: [
      'No network calls, single static binary, fail-closed',
      'PDF + JSON + bundle output for working papers',
      'Cosign + GPG + reproducible build trust path',
      'Strict mode for examiners; non-strict for SOC',
      'Recovery scenarios for WAL-from-older-backup'
    ]
  },
  {
    id: 'design-test-vectors',
    title: 'Design — Conformance Test Vectors',
    category: 'design',
    path: '../docs/design/08-test-vectors.md',
    summary: 'The conformance corpus: structure, taxonomy, cross-implementation conformance, regulatory significance.',
    analogy: 'The driver\'s test for implementations. Pass the test, you\'re a conforming driver. Two licensed drivers see the same road the same way.',
    keyPoints: [
      'Numbered test cases with self-contained inputs/expected',
      'Negative tests for false-positive prevention',
      'Cross-day, DAG, and master-rotation vectors (v1.0-final)',
      'Append-only corpus growth across spec versions',
      'Cross-implementation fuzzing for divergence detection'
    ]
  },
  {
    id: 'threat-model',
    title: 'Design — Threat Model',
    category: 'design',
    path: '../docs/design/09-threat-model.md',
    summary: 'Adversaries A through H, properties claimed, residual risks, identity-attribution composability, quantum readiness.',
    analogy: 'The "who would attack this and how do we stop them" analysis. For each attacker, what they can do, what stops them, what residual risk remains.',
    keyPoints: [
      'Eight adversaries: external wire, DBA insider, ops insider, HSM physical, vendor compromise, app compromise, master-key compromise, crypto break',
      'Three integrity properties: authentic capture, tamper evidence, independent verifiability',
      'Backup-integrity boundary, identity-attribution composability',
      'Quantum-readiness roadmap with 30-day emergency-patch SLA',
      'Forgery-allegation evidence path'
    ]
  },
  {
    id: 'design-glossary',
    title: 'Design — Glossary',
    category: 'design',
    path: '../docs/design/10-glossary.md',
    summary: 'Definitions of all domain terms: cryptographic primitives, chain components, deployment patterns, audit conventions.',
    analogy: 'The dictionary. Open when a term is unfamiliar; close when it is not.',
    keyPoints: [
      '110+ terms organized alphabetically',
      'Cross-references back to design docs',
      'Definitions in plain prose, not formulas',
      'Includes both cryptographic and operational vocabulary'
    ]
  },

  // Regulator pack
  {
    id: 'finding-language',
    title: 'Examination Finding Language',
    category: 'regulator',
    path: '../docs/regulator-pack/finding-language.md',
    summary: 'Standard examination-report language for verifier-detected findings, including repeat findings and public disclosure.',
    keyPoints: [
      'Severity guidance per failure mode',
      'Sample finding paragraphs',
      'Repeat-finding escalation language',
      'Public-disclosure language for consent orders',
      'Enforcement-action lifecycle'
    ]
  },
  {
    id: 'csf-2',
    title: 'NIST CSF 2.0 Mapping',
    category: 'regulator',
    path: '../docs/regulator-pack/CSF-2.0.md',
    summary: 'Maps chain primitives to NIST Cybersecurity Framework 2.0 functions and subcategories.',
    analogy: 'A translation guide between the chain and the NIST CSF: what each primitive supports in CSF terms.',
    keyPoints: [
      'PR.DS-06 is the headline (data integrity verification)',
      'Coverage across GOVERN, IDENTIFY, PROTECT, DETECT, RESPOND, RECOVER',
      'Explicit list of subcategories the chain does NOT satisfy',
      'Composes with broader cybersecurity posture'
    ]
  },
  {
    id: 'handbook-mapping',
    title: 'FFIEC IT Handbook Mapping',
    category: 'regulator',
    path: '../docs/regulator-pack/handbook-mapping.md',
    summary: 'Maps chain to specific control objectives in IS, AIO, Audit, and Outsourcing booklets.',
    keyPoints: [
      'IS booklet II.C.10 (logging) is the headline',
      'AIO booklet sections on resilience, capacity, change management',
      'Audit booklet alignment with independent assessment',
      'OTS booklet for vendor-hosted topology'
    ]
  },
  {
    id: 'deployment-package',
    title: 'Verifier Deployment Package',
    category: 'regulator',
    path: '../docs/regulator-pack/deployment-package.md',
    summary: 'Information the regulator\'s IT shop needs to allowlist and deploy the verifier.',
    keyPoints: [
      'Binary specs (size, type, network/FS behavior)',
      'Per-release trust artifacts (cosign, GPG manifest, SBOM)',
      'Two allowlisting postures (SHA-256 vs cosign signature)',
      'Sample command lines',
      'Validate-before-run wrapper script'
    ]
  },
  {
    id: 'sample-report',
    title: 'Sample Verifier Report',
    category: 'regulator',
    path: '../docs/regulator-pack/sample-report.md',
    summary: 'Completed verifier report from a fictional tenant for examiner training.',
    analogy: 'The model answer key. Examiners learn what passing looks like by reading a real example.',
    keyPoints: [
      'Cover page summary (pass/fail, counts)',
      'Per-day detail with anomaly notes',
      'Methodology section showing algorithms used',
      'How an examiner reads the report',
      'Sample finding language for the report'
    ]
  },
  {
    id: 'examiner-training',
    title: 'Examiner Training (30-Minute)',
    category: 'regulator',
    path: '../docs/regulator-pack/examiner-training.md',
    summary: 'Onboarding a new examiner from "I have a snapshot" to "I have a defensible report" in 30 minutes.',
    keyPoints: [
      'Module 1: concepts (5 min)',
      'Module 2: inputs the examiner needs (5 min)',
      'Module 3: run the verifier (10 min hands-on)',
      'Module 4: read the report (5 min)',
      'Module 5: common patterns (5 min)'
    ]
  },
  {
    id: 'examiner-approval',
    title: 'Examiner Approval Template',
    category: 'regulator',
    path: '../docs/regulator-pack/examiner-approval-template.md',
    summary: 'Template for an institution\'s request to relax seal cadence (daily → weekly).',
    keyPoints: [
      'Institution identification + tenant_ids',
      'Current vs proposed cadence',
      'Rationale and compensating controls',
      'Reversal commitment',
      'Regulator response format (approval / approval with mods / denial)'
    ]
  },
  {
    id: 'ai-policy-alignment',
    title: 'AI Policy Alignment',
    category: 'regulator',
    path: '../docs/regulator-pack/ai-policy-alignment.md',
    summary: 'How the chain aligns with U.S. and international AI policy: NIST AI RMF, EO 14110, EU AI Act, EBA, DORA, NIS2, FCA, MAS, JFSA.',
    keyPoints: [
      'EU AI Act Article 12 record-keeping is the headline international fit',
      'NIST AI RMF Manage 4.x alignment',
      'DORA / NIS2 timing notes',
      'State-level (CCPA, VCDPA, CPA, CTDPA, TDPSA)',
      'Multi-jurisdictional adaptation'
    ]
  },

  // Control map
  {
    id: 'cuecs',
    title: 'Complementary User Entity Controls',
    category: 'controls',
    path: '../docs/control-map/CUECs.md',
    summary: 'The 21 controls the institution must operate for the chain\'s claims to hold, organized into three tiers.',
    analogy: 'The bank\'s share of the responsibility. The chain provides the lock; CUECs are the bank\'s practices for using the lock correctly.',
    keyPoints: [
      'Tier 1 — Critical (operate from day 1)',
      'Tier 2 — Important (within 6 months)',
      'Tier 3 — Operational hygiene (within first year)',
      'Mapping back to chain primitives',
      'Ramp-up timeline aligned with maturity'
    ]
  },
  {
    id: 'tsc-mapping',
    title: 'SOC 2 TSC Mapping',
    category: 'controls',
    path: '../docs/control-map/TSC-mapping.md',
    summary: 'Maps chain primitives to AICPA Trust Services Criteria including criterion-level Privacy breakdown.',
    keyPoints: [
      'Common Criteria (CC1-CC9) coverage',
      'Processing Integrity (PI1.1, PI1.2) headline',
      'Availability for institutions claiming A',
      'Privacy criterion-level breakdown (P1-P8)',
      'What the chain does NOT satisfy'
    ]
  },

  // SOC pack
  {
    id: 'soc-section-4',
    title: 'SOC Section 4 Template',
    category: 'soc',
    path: '../docs/soc-pack/section-4-template.md',
    summary: 'Starter text for Service Organization Description (Section 4) of a SOC report on a chain implementation.',
    keyPoints: [
      'Overview of services',
      'Principal commitments and system requirements',
      'Components (infrastructure, software, people, procedures, data)',
      'System boundaries (in scope vs out)',
      'CUECs cross-reference',
      'Pre-issuance review checklist'
    ]
  },
  {
    id: 'control-evidence-events',
    title: 'Control-Evidence Operational Events',
    category: 'soc',
    path: '../docs/soc-pack/control-evidence-events.md',
    summary: 'Schema for operational events emitted as control evidence; SOC and examination teams consume mechanically.',
    analogy: 'The factory\'s machine logs. Every gear that turns leaves a record. SOC auditors read the logs to prove the gears actually turned.',
    keyPoints: [
      'Standard event names and field schemas',
      'Ledger lifecycle, seal job, chain integrity',
      'HSM operations, configuration, master-key rotation',
      'Reconciliation events',
      'Retention matches chain events'
    ]
  },

  // Audience summaries
  {
    id: 'management-summary',
    title: 'Management Summary',
    category: 'audience',
    path: '../docs/management-summary.md',
    summary: 'Plain-English overview an examiner can hand to a bank CEO. What the chain does, what it doesn\'t, what to expect.',
    analogy: 'The flight-simulator briefing for a CEO. You don\'t need to fly the plane; you need to understand what the cockpit instruments say.',
    keyPoints: [
      'What was tested and what was found',
      'Pass / pass-with-anomalies / fail interpretation',
      'What the CEO should know in each case',
      'What the chain does NOT do',
      'FAQs CEOs typically ask'
    ]
  },
  {
    id: 'audit-committee-summary',
    title: 'Audit Committee Summary',
    category: 'audience',
    path: '../docs/audit-committee-summary.md',
    summary: 'Three-page brief for the bank\'s audit committee chair on chain oversight responsibilities.',
    keyPoints: [
      'What the committee oversees and why',
      'Quarterly and annual reporting expectations',
      'Oversight questions to ask',
      'Where chain meets broader committee responsibilities',
      'Decisions the committee may take'
    ]
  },
  {
    id: 'mrm-committee-brief',
    title: 'MRM Committee Brief',
    category: 'audience',
    path: '../docs/MRM-COMMITTEE-BRIEF.md',
    summary: 'Two-to-four page brief the bank\'s MRM committee chair reads instead of full design docs.',
    keyPoints: [
      'What the chain does for SR 11-7',
      'How it changes MRM evidence',
      'Specific scenarios and committee responses',
      'Standard questions for the chain owner',
      'Cost summary and decisions'
    ]
  },
  {
    id: 'user-entity-summary',
    title: 'User-Entity Summary',
    category: 'audience',
    path: '../docs/user-entity-summary.md',
    summary: 'For SOC 1 user entities reading a SOC report on a chain-of-custody implementation.',
    keyPoints: [
      'Audience: external auditor, downstream partner, regulator, audit firm',
      'What the SOC report attests to',
      'CUEC verification scope',
      'How user-entity uses the report',
      'Boundary: what chain does NOT do for the user entity'
    ]
  },

  // Operations
  {
    id: 'operator-guide',
    title: 'Operator Guide',
    category: 'operations',
    path: '../docs/operator-guide.md',
    summary: 'Bank-operator-facing guide consolidating runtime operations, deployment, daily ops, periodic ops.',
    keyPoints: [
      'Topology, configuration, HSM provisioning',
      'Daily operations (health, metrics, logs)',
      'Quarterly / annual / as-needed periodic operations',
      'Master-key rotation procedure',
      'Spec version migration'
    ]
  },
  {
    id: 'cloud-hsm-guide',
    title: 'Cloud HSM Guide',
    category: 'operations',
    path: '../docs/cloud-hsm-guide.md',
    summary: 'Per-cloud guidance for HSM provisioning that satisfies FIPS 140-2 L3 requirement.',
    keyPoints: [
      'Conformant matrix: AWS / Azure / Google',
      'Disqualified options (KMS default tier, Key Vault Standard)',
      'Cross-region replication notes',
      'Cost ranges per provider',
      'Common pitfalls'
    ]
  },
  {
    id: 'byoc-deployment',
    title: 'BYOC Deployment',
    category: 'operations',
    path: '../docs/byoc-deployment.md',
    summary: 'Bring-Your-Own-Cloud deployment: bank runs vendor\'s image; IAM and network boundaries enforce separation.',
    analogy: 'The bank rents the vendor\'s recipe but cooks in the bank\'s own kitchen. Vendor never gets keys to the kitchen.',
    keyPoints: [
      'IAM permission matrix (bank vs vendor vs runtime)',
      'Network controls: HSM in private subnet, mirror registry',
      'Image-pull egress and signature verification',
      'Vendor support telemetry routing with redaction',
      'Audit considerations'
    ]
  },
  {
    id: 'supply-chain',
    title: 'Supply Chain',
    category: 'operations',
    path: '../docs/supply-chain.md',
    summary: 'Verifier and ledger-server build pipeline, signing artifacts, project-side governance.',
    keyPoints: [
      'Per-release artifacts (binaries, cosign, GPG manifest, SBOM, scan)',
      'Cosign + GPG dual signing',
      'Reproducible builds (deterministic flags)',
      'Project-side governance: maintainer controls, corpus integrity',
      'CycloneDX 1.5 SBOM format'
    ]
  },
  {
    id: 'dr-and-resilience',
    title: 'DR and Resilience',
    category: 'operations',
    path: '../docs/dr-and-resilience.md',
    summary: 'Reference RPO/RTO targets, sync-vs-async replication trade-offs, multi-region workaround.',
    keyPoints: [
      'Reference targets per component',
      'Synchronous (sub-second RPO, +latency) vs async',
      'Failure scenarios and recovery',
      'Multi-region per-region tenant_id workaround',
      'DR exercise procedures'
    ]
  },
  {
    id: 'at-scale-operations',
    title: 'At-Scale Operations',
    category: 'operations',
    path: '../docs/at-scale-operations.md',
    summary: 'Operational guidance for billions/day deployments: throughput, sharding, emergency rotation, reconciliation at scale.',
    keyPoints: [
      'Per-host throughput ceilings',
      'Sharding patterns (per-process, per-tenant, per-business-line)',
      'Emergency rotation forced-handshake mechanism',
      'Reconciliation baseline establishment',
      'Long-tail retention scenarios',
      'Holding-company complexity'
    ]
  },
  {
    id: 'edge-and-federated-ai',
    title: 'Edge and Federated AI',
    category: 'operations',
    path: '../docs/edge-and-federated-ai.md',
    summary: 'Deployment guidance for edge devices, federated learning, on-device inference, multi-agent autonomous coordination.',
    analogy: 'The chain works the same whether the AI is on a server in a data center or on a doorbell camera. The substrate scales down without changing.',
    keyPoints: [
      'Edge AI: SDK on edge devices, local-first persistence',
      'Federated learning: chain captures inference, not training',
      'Multi-agent autonomous coordination via DAG semantics',
      'On-device inference and privacy-by-design',
      'Operational considerations per pattern'
    ]
  },
  {
    id: 'vendor-hosted-controls',
    title: 'Vendor-Hosted Controls',
    category: 'operations',
    path: '../docs/vendor-hosted-controls.md',
    summary: 'Control distribution between institution and vendor when the vendor operates the chain implementation.',
    keyPoints: [
      'Institution-side vs vendor-side controls',
      'CUEC redistribution in vendor-hosted',
      'Vendor SOC report consumption',
      'Vendor relationship transition',
      'Common vendor-hosted issues'
    ]
  },

  // Adoption & onboarding
  {
    id: 'cost-model',
    title: 'Cost Model',
    category: 'adoption',
    path: '../docs/cost-model.md',
    summary: 'All-in cost picture by tier (community, mid-size, Tier-1) with cost trajectory over years.',
    keyPoints: [
      'Per-tier annual cost ranges',
      'HSM dominates 60-80% of total',
      'Cost trajectory year 1 → year 7',
      'Cost reduction options',
      'Worked examples per tier'
    ]
  },
  {
    id: 'minimum-viable-deployment',
    title: 'Minimum-Viable Deployment',
    category: 'adoption',
    path: '../docs/minimum-viable-deployment.md',
    summary: 'Smallest conformant deployment shape for community-bank pilots and first-time adopters.',
    analogy: 'The starter kit. Conformant, defensible, sized for the smallest viable use case. Add the larger kit as the bank grows.',
    keyPoints: [
      'When to use minimum-viable',
      'Tier 1 CUECs from day 1',
      'Cost target ~$37k/year',
      '12-month operational ramp-up timeline',
      'When to upgrade to mid-size'
    ]
  },
  {
    id: 'first-engagement-guide',
    title: 'First Engagement Guide',
    category: 'adoption',
    path: '../docs/first-engagement-guide.md',
    summary: 'Bank-side and EIC-side guide for the institution\'s first chain examination.',
    keyPoints: [
      'Pre-examination preparation (both sides)',
      'During-examination workflow',
      'Common first-engagement issues',
      'Post-engagement learning loop',
      'Resources during the engagement'
    ]
  },
  {
    id: 'examiner-quickstart',
    title: 'Examiner Quickstart',
    category: 'adoption',
    path: '../docs/examiner-quickstart.md',
    summary: '5-minute orientation for an examiner who has never run the chain verifier.',
    keyPoints: [
      'What the chain does in three sentences',
      'What you need (5 inputs)',
      'What you do (4 commands)',
      'What the report tells you',
      'Common failure modes'
    ]
  },

  // Incident & legal
  {
    id: 'incident-response-playbook',
    title: 'Incident Response Playbook',
    category: 'incident',
    path: '../docs/incident-response-playbook.md',
    summary: 'IR playbook for chain-detected events: scenarios, severity classification, roles, notification frameworks.',
    analogy: 'The fire-drill manual. When the chain alarm sounds, here\'s exactly who does what, in what order, with what notification.',
    keyPoints: [
      'Six common scenarios (chain hash mismatch through software-key fallback)',
      'CIRCIA, FFIEC 36-hour, state, sealing-delay notification paths',
      'HSM tamper-detection integration',
      'Long-dwell adversary considerations',
      'Opaque agent integration patterns'
    ]
  },
  {
    id: 'legal-disclosure',
    title: 'Legal Disclosure Procedures',
    category: 'incident',
    path: '../docs/legal-disclosure.md',
    summary: 'Institution\'s posture on court-ordered key disclosure, subpoenas, FRE 901/902 admissibility, expert testimony.',
    keyPoints: [
      'Court-ordered master-key disclosure → comply, then rotate',
      'HSM-protected signing key is non-extractable',
      'FRE 901 / 902 / 803(6) for chain output as evidence',
      'Cryptographic-expert testimony alignment',
      'Privilege and long-tail retention'
    ]
  },
  {
    id: 'customer-dispute-procedures',
    title: 'Customer Dispute Procedures',
    category: 'incident',
    path: '../docs/customer-dispute-procedures.md',
    summary: 'Customer-side path for disputing AI-driven decisions and the institution\'s response.',
    keyPoints: [
      'Dispute paths: direct, regulator complaint, litigation',
      'Bank\'s initial response with chain walk',
      'How the customer can verify integrity independently',
      'ECOA adverse-action notice flow',
      'CFPB-specific procedures'
    ]
  },

  // Audit support
  {
    id: 'audit-procedures',
    title: 'Audit Procedures',
    category: 'audit',
    path: '../docs/audit-procedures.md',
    summary: '21 sample testing procedures for SOC and FFIEC examiners testing chain controls.',
    keyPoints: [
      'Procedures by control category (IAM, crypto, ops, verifier, IR, vendor, config)',
      'Sampling guidance per population size',
      'Testing approach (claim → evidence → test → execute → document)',
      'Working-paper preservation',
      'Coordination with FFIEC examination'
    ]
  },
  {
    id: 'anomaly-template',
    title: 'Anomaly Documentation Template',
    category: 'audit',
    path: '../docs/anomaly-documentation-template.md',
    summary: 'Template for documenting verifier-reported anomalies during a SOC reporting period.',
    keyPoints: [
      'When to fill in the record',
      'Severity assessment guidance',
      '"Operationally explained" criteria',
      'SOC engagement use',
      'FFIEC examination use'
    ]
  },
  {
    id: 'portfolio-comparison-procedures',
    title: 'Portfolio Comparison Procedures',
    category: 'audit',
    path: '../docs/portfolio-comparison-procedures.md',
    summary: 'Manual procedures for cross-bank portfolio comparison and multi-region cross-comparison.',
    analogy: 'The supervisory bird\'s-eye view. Each bank produces its own report; the regulator compares across the flock to spot outliers.',
    keyPoints: [
      'Cross-bank portfolio comparison metrics',
      'Multi-region cross-comparison within institution',
      'Cross-period comparison (longitudinal)',
      'Patterns and their implications',
      'Working-paper preservation'
    ]
  },
  {
    id: 'm-and-a-handoff',
    title: 'M&A and Corporate Transactions',
    category: 'audit',
    path: '../docs/m-and-a-handoff.md',
    summary: 'Procedures for chain operation during mergers, acquisitions, divestitures, spin-offs, vendor changes.',
    keyPoints: [
      'Past public keys remain valid via tenant key registry',
      'Master-key transfer via HSM-vendor procedures',
      'Pattern A (combined tenant_id forward) vs Pattern B (renaming)',
      'SOC reporting during transactions',
      'Examination during transition periods'
    ]
  },

  // Specialized
  {
    id: 'privacy-by-design',
    title: 'Privacy by Design',
    category: 'special',
    path: '../docs/privacy-by-design.md',
    summary: 'Privacy patterns: tokenize PII before chain capture; maintain a separate privacy-store the institution can erase.',
    analogy: 'The chain captures the seat number, not the passenger\'s name. The airline\'s passenger list (the privacy store) holds the name and is erasable.',
    keyPoints: [
      'Tokenization at SDK time before canonicalization',
      'GDPR Article 17 right to erasure',
      'CCPA / state law alignment',
      'TSC Privacy criterion P1-P8 breakdown',
      'Cross-border data transfers'
    ]
  }
];

// Quick stats for the home page
const STATS = {
  totalDocs: DOCUMENTS.length,
  totalStakeholders: STAKEHOLDERS.length - 1, // excluding 'all'
  primitives: 4,
  cuecs: 21,
  categories: CATEGORIES.length
};

// ============================================================
// COST MODEL — components, per-tier breakdown, reduction options
// ============================================================

// What makes up the cost. Shows readers where their money actually goes.
const COST_COMPONENTS = [
  {
    name: 'HSM operations',
    share: '60–80%',
    icon: 'mdi-shield-key-outline',
    detail: 'FIPS 140-2 Level 3+ Hardware Security Module cluster annual fee. Cloud HSM (AWS CloudHSM, Azure Managed HSM, Google Cloud HSM) runs roughly $13,000–$15,000 per cluster member per year. Production deployments use at least two members for HA. The HSM is the integrity anchor — software costs nothing, but the HSM\'s tamper-resistance is what makes the daily seal undefiable.'
  },
  {
    name: 'Compute',
    share: '5–15%',
    icon: 'mdi-server',
    detail: 'Ledger server instances. Sizing scales with event volume. A community bank runs a single small instance (~$2,500/year). A mid-size bank runs 2–3 medium instances (~$8,000–$12,000/year). A Tier-1 runs 4–8 instances per region across multiple regions (~$40,000+/year).'
  },
  {
    name: 'Storage',
    share: '5–15%',
    icon: 'mdi-database-outline',
    detail: 'Hot store (Postgres + JSONB) plus cold store (Parquet + Iceberg or S3 archives). Sized to the institution\'s retention period (typically 7 years for U.S. financial services). Community: ~$500/year. Mid-size: ~$4,500/year. Tier-1: ~$45,000+/year for multi-TB cold-store with multi-region replication.'
  },
  {
    name: 'Operational headcount',
    share: '10–25%',
    icon: 'mdi-account-hard-hat-outline',
    detail: 'Chain operations team. Community bank typically allocates 0.1 FTE (a fraction of an existing IT operator). Mid-size: 0.25–0.5 FTE. Tier-1: 0.5–1 dedicated FTE plus a part-time IR specialist. Includes seal-job monitoring, incident response, master-key rotation procedures.'
  },
  {
    name: 'Software',
    share: '0% direct',
    icon: 'mdi-file-code-outline',
    detail: 'The reference implementation is Apache 2.0 open source. No license cost. The opportunity cost is implementation effort (4–12 week initial deployment project) and ongoing operations.'
  },
  {
    name: 'SOC engagement (separate)',
    share: 'separate',
    icon: 'mdi-clipboard-check-outline',
    detail: 'A SOC 2 Type II engagement covering the chain costs roughly $50,000–$200,000 from a Big Four firm, depending on scope and complexity. Not chain-specific — same as any other controls engagement. Not always required; institutions decide based on their broader audit program.'
  }
];

// Per-tier deployment shapes with worked examples explaining how the cost adds up.
const COST_TIERS = [
  {
    tier: 'Community Bank ($1B–$10B)',
    annual: '$25,000 – $60,000',
    icon: 'mdi-bank-outline',
    config: 'Single AWS CloudHSM cluster (2 members), one ledger instance, daily seal cadence, single region.',
    workedExample: {
      title: 'Worked example — $5B community bank using AI for customer-service routing',
      lines: [
        { label: 'AWS CloudHSM (2 members)', cost: '$26,000' },
        { label: 'Single ledger Postgres instance', cost: '$2,500' },
        { label: 'Storage (10 GB/year × 7 years retention)', cost: '$500' },
        { label: 'Operational overhead (0.1 FTE)', cost: '$20,000' },
        { label: 'Total', cost: '~$49,000/year', total: true }
      ]
    },
    drivers: [
      'Low AI volume — typically thousands of events per day',
      'Single AI use case (often customer-service or fraud-screening)',
      'Modest resilience needs — single region acceptable',
      'Examination cycle is 18-month, allowing weekly cadence with examiner approval'
    ]
  },
  {
    tier: 'Mid-Size Bank ($10B–$50B)',
    annual: '$80,000 – $200,000',
    icon: 'mdi-domain',
    config: 'Dedicated cloud HSM, 2–3 ledger instances, daily cadence, single region with cross-region cold-store backup.',
    workedExample: {
      title: 'Worked example — $25B bank using AI for fraud screening + customer service',
      lines: [
        { label: 'Azure Managed HSM partition', cost: '$15,000' },
        { label: 'Two ledger instances (8 vCPU each)', cost: '$8,000' },
        { label: 'Storage (200 GB/year × 7 years)', cost: '$4,500' },
        { label: 'Chain operations (0.3 FTE)', cost: '$60,000' },
        { label: 'Total', cost: '~$87,500/year', total: true }
      ]
    },
    drivers: [
      'Multiple AI use cases (fraud + customer service + emerging)',
      'Event volume in low millions/day',
      'Standard daily cadence',
      'Cost varies with operational headcount and number of distinct AI programs'
    ]
  },
  {
    tier: 'Tier-1 / G-SIB ($50B+)',
    annual: '$300,000 – $1,000,000+',
    icon: 'mdi-bank',
    config: 'Multi-region HSM cluster (4+ members across 4 regions), multi-region ledger compute, 5+ TB/year storage, dedicated chain ops team.',
    workedExample: {
      title: 'Worked example — $250B globally-systemic bank using AI across multiple business lines',
      lines: [
        { label: '4 cloud HSM clusters × 4 regions', cost: '$208,000' },
        { label: 'Multi-region ledger compute (16 instances)', cost: '$40,000' },
        { label: 'Storage (5 TB/year × multi-region replication)', cost: '$50,000' },
        { label: 'Chain operations (0.5 FTE)', cost: '$80,000' },
        { label: 'Total', cost: '~$378,000/year', total: true }
      ]
    },
    drivers: [
      'Billions of events per day across business lines',
      'Multi-region resilience required for operational continuity',
      'Multiple business lines warrant separate tenants for risk isolation',
      'May require multi-vendor implementation for diversity defense',
      'Larger institutions can exceed $1M with active-active multi-region and dedicated SOC engagements'
    ]
  }
];

// Cost reduction options. The institution's CFO uses this menu.
const COST_REDUCTION_OPTIONS = [
  {
    option: 'Relax cadence to weekly',
    saving: 'Reduces HSM signing operations by ~7×',
    tradeoff: 'Wider retroactive-tamper-detection window (1 day → 7 days). Requires examiner approval per the cadence-approval workflow.',
    bestFor: 'Community banks with 18-month examination cycle and modest AI use'
  },
  {
    option: 'Shared cloud HSM (vendor-hosted)',
    saving: 'Reduces HSM cost by ~50%',
    tradeoff: 'Multi-tenancy concerns documented in control description. Vendor-hosted topology shifts some controls to the vendor.',
    bestFor: 'Smaller institutions with limited risk-isolation needs'
  },
  {
    option: 'Cold-only retention beyond year 1',
    saving: 'Reduces hot-storage cost by 80%+',
    tradeoff: 'Slower query for older events; retrieval from cold-store takes minutes rather than milliseconds.',
    bestFor: 'Audit-only access pattern for older data'
  },
  {
    option: 'Single region (defer multi-region)',
    saving: 'Reduces HSM cost by region count factor',
    tradeoff: 'Lower resilience. Acceptable for institutions whose risk profile accepts single-region.',
    bestFor: 'Institutions whose existing DR program is single-region'
  },
  {
    option: 'Self-host on existing on-prem HSM',
    saving: 'Defers ongoing cloud-HSM annual fees',
    tradeoff: 'Higher upfront capex; ongoing ops headcount. Per-operation cost is lower at high volume.',
    bestFor: 'Institutions with mature HSM operations and significant volume'
  }
];

// ============================================================
// THREAT MODEL — eight adversaries with detail and links
// ============================================================
const THREATS = [
  {
    id: 'A',
    name: 'External wire attacker',
    icon: 'mdi-wifi-strength-2-alert',
    capability: 'Network-level attacker who can inject, modify, or drop OTLP messages between the SDK and the ledger server. Likely landing point: man-in-the-middle on internal network, compromised proxy, or rogue OTel collector.',
    goal: 'Forge AI events into the captured stream that the institution did not produce. If successful, the institution\'s evidence base is corrupted.',
    defense: 'HMAC chain at capture. Each event\'s payload_hash is HMAC-SHA-256 over the previous event\'s hash plus the canonical payload, keyed with the per-process session key. An attacker without the key cannot produce an event whose payload_hash matches the chain. The ledger re-verifies HMAC on ingest and rejects mismatches; failed verifications become integrity-alert events routed to the institution\'s SIEM.',
    residual: 'An attacker who has compromised the application process and recovered its session key can produce valid-looking events. Detection comes from the institution\'s incident-response framework. The compromise window is bounded by the reconciliation cadence (weekly per spec; faster for institutions with elevated risk profile).',
    relatedDocs: ['threat-model', 'design-chain', 'design-otlp']
  },
  {
    id: 'B',
    name: 'Insider with database access',
    icon: 'mdi-database-edit-outline',
    capability: 'Direct read and write access to the ledger storage (DBA, infrastructure operator with elevated privileges). Bypasses the application-level append-only invariant.',
    goal: 'Rewrite past events to favor a different decision history — make a denied loan look approved, or a flagged transaction look unflagged.',
    defense: 'Daily Merkle seal under HSM signature. The Merkle root is committed under HSM-protected Ed25519 signature; modifying any past event changes the root, which the verifier catches by recomputing from the ledger and comparing to the signed root. The insider cannot re-sign because the HSM holds the signing key non-extractably.',
    residual: 'The insider can corrupt the database (denial of service against the verifier), but cannot silently rewrite history. The recomputed-versus-signed mismatch is the detection signal. RBAC defense-in-depth (INSERT and SELECT only on the events table) raises the operational floor.',
    relatedDocs: ['threat-model', 'design-merkle', 'design-hsm']
  },
  {
    id: 'C',
    name: 'Insider with operational access',
    icon: 'mdi-account-key-outline',
    capability: 'Authorized to invoke the daily seal job and supply payloads to the HSM for signing. Has the seal-job operator credentials.',
    goal: 'Construct a plausible alternate history and sign it under the institution\'s key. If successful, the seal verifies against the public key, and the chain shows authentic.',
    defense: 'Separation of duties. The seal-job operator role is separated from the database admin role and from the master-key custody role. To execute the attack, the insider needs the database admin\'s ability to alter ledger contents and the seal-job operator\'s ability to sign. Banks enforce this separation for other privileged roles (payment authorization, change management); the chain leverages the existing pattern.',
    residual: 'Collusion between the database admin and the seal-job operator. Detection: independent audit (internal or external) periodically samples ledger consistency; the regulator\'s independent verification is the long-cycle defense.',
    relatedDocs: ['threat-model', 'design-hsm', 'cuecs']
  },
  {
    id: 'D',
    name: 'HSM physical attack',
    icon: 'mdi-shield-alert-outline',
    capability: 'Physical possession of the HSM long enough to attempt key extraction. Requires physical access to the data center or cloud HSM facility.',
    goal: 'Extract the Ed25519 signing key. With the key, the attacker can forge daily roots indistinguishable from legitimate ones.',
    defense: 'FIPS 140-2 Level 3 (or higher) physical tamper resistance. The HSM is designed to detect physical compromise and either zeroize keys or refuse to operate. AWS CloudHSM, Azure Dedicated HSM, and on-prem appliances at L3+ have physically passed FIPS validation.',
    residual: 'A nation-state-level adversary with sufficient resources may compromise an HSM. The threat model accepts this as residual; the institution\'s incident response includes notifying the regulator of any HSM physical-security incident, plus immediate key rotation.',
    relatedDocs: ['threat-model', 'design-hsm', 'cloud-hsm-guide']
  },
  {
    id: 'E',
    name: 'Vendor compromise',
    icon: 'mdi-package-variant-closed-remove',
    capability: 'The vendor providing the chain implementation (SDK or ledger) is compromised — insider, supply-chain attack, or build-pipeline subversion.',
    goal: 'Subvert the chain at build time so the vendor (or a deeper attacker) can forge events on demand.',
    defense: 'Multi-layered. Apache 2.0 open source means every line is publicly auditable. Reproducible builds let independent parties rebuild from source and confirm the published binary matches. Cosign signatures with a GPG-signed manifest fallback. The conformance corpus catches behavioral deviation; a subverted build that produces incorrect output fails the corpus.',
    residual: 'A subversion that passes the corpus and produces deterministic-but-wrong output is theoretically possible (the attacker would have to compromise both the corpus and the implementation). Multi-party review and the security disclosure process are the defenses.',
    relatedDocs: ['threat-model', 'supply-chain', 'design-test-vectors']
  },
  {
    id: 'F',
    name: 'Application process compromise',
    icon: 'mdi-bug-outline',
    capability: 'The AI agent\'s host process is compromised — RCE, container escape, malicious dependency, supply-chain attack on the agent platform itself.',
    goal: 'Forge events going forward, claiming they reflect legitimate AI decisions.',
    defense: 'Bounded forward-only attack window. The compromised process holds the session key and can produce valid-looking events. Past events cannot be retroactively altered — they have already been chained, exported, and (eventually) sealed. The institution\'s incident response identifies the compromise; events from the compromise window are flagged.',
    residual: 'The compromise window between when the attacker gains access and when the institution detects it. Detection mechanisms: anomaly detection on the captured stream, out-of-band monitoring of the agent\'s behavior, intrusion detection.',
    relatedDocs: ['threat-model', 'design-chain', 'incident-response-playbook']
  },
  {
    id: 'G',
    name: 'Master-key exfiltration',
    icon: 'mdi-key-remove',
    capability: 'The tenant master HMAC key is exfiltrated from the master-key custodian.',
    goal: 'Derive arbitrary session keys and forge events without needing to compromise individual application processes.',
    defense: 'Rotation. On detection, the institution rotates the master. New session keys derive from the new master; events captured after rotation pass verification only against the new master. Past events remain verifiable against the previous master (the institution retains key version history). Weekly session-key-id reconciliation bounds the master-compromise detection window.',
    residual: 'Events captured during the compromise window between exfiltration and detection are repudiable — the institution cannot prove they were produced by legitimate processes versus by the attacker holding the leaked key. Forensic analysis (network logs, process audit logs, HSM operations log) is the supplementary evidence path.',
    relatedDocs: ['threat-model', 'design-hsm', 'incident-response-playbook']
  },
  {
    id: 'H',
    name: 'Cryptographic break',
    icon: 'mdi-key-alert-outline',
    capability: 'A practical break of HMAC-SHA-256, SHA-256, or Ed25519 emerges. Most concerning: a quantum computer capable of breaking Ed25519 in operationally feasible time.',
    goal: 'Forge events or seals without holding the original keys.',
    defense: 'Algorithm rotation. The seal record carries an explicit algorithm identifier. A future spec version can specify a new algorithm; institutions rotate to new keys under the new algorithm. The 30-day spec-patch SLA commits the working group to publishing an emergency spec patch within 30 days of credible demonstration.',
    residual: 'The window between when a break becomes practical and when institutions can rotate. The spec admits the new algorithm immediately, but bank-internal change management has its own cadence (180 days for signature break, 90 days for HMAC break per spec).',
    relatedDocs: ['threat-model', 'design-hsm', 'spec-main']
  }
];

// ============================================================
// AUDIT PROCESS — timeline of actions from notice through follow-up
// Generic FFIEC examination + SOC engagement view. Specific cadences
// vary by regulator and engagement letter.
// ============================================================
const AUDIT_STAGES = [
  {
    id: 'engagement-notice',
    name: 'Engagement Notice',
    timing: 'T−90 to T−60 days',
    icon: 'mdi-email-outline',
    color: 'primary',
    summary: 'Examiner or audit firm formally notifies the institution. The clock starts.',
    auditor: {
      actions: [
        'Send engagement letter (SOC) or examination notice (FFIEC)',
        'Set scope: period, tenant_ids, business lines, criteria',
        'Identify lead auditor and team composition',
        'Plan resource allocation and timeline'
      ],
      lookFor: [
        'Whether the institution operates the chain at all',
        'Prior examination findings or open MRAs (Matters Requiring Attention)',
        'Material changes since last examination (vendor switch, spec version bump, M&A)',
        'Public-record incidents (cyber notifications, customer complaints)'
      ]
    },
    auditee: {
      actions: [
        'Acknowledge notice within the regulator\'s required timeline',
        'Identify primary contact and engagement liaison',
        'Brief the audit committee chair and CISO',
        'Begin compiling current control descriptions'
      ],
      provide: []
    },
    relatedDocs: ['first-engagement-guide']
  },
  {
    id: 'pre-prep',
    name: 'Pre-Examination Preparation',
    timing: 'T−60 to T−14 days',
    icon: 'mdi-clipboard-list-outline',
    color: 'primary',
    summary: 'Both sides prepare. Institution runs the verifier internally; auditor reviews documentation.',
    auditor: {
      actions: [
        'Review institution\'s control descriptions and SOC reports',
        'Identify sample population for testing',
        'Prepare audit-procedures checklist',
        'Pre-allocate examination time on examiner\'s laptop',
        'Validate verifier binary using verifier-validate.sh'
      ],
      lookFor: [
        'Completeness of the institution\'s control description',
        'Currency of the institution\'s SOC report',
        'Spec version the institution is operating',
        'Whether claimed cadence matches recent seal records',
        'Whether CUECs are operational at the declared tier'
      ]
    },
    auditee: {
      actions: [
        'Run pre-examination verifier on production ledger',
        'Prepare ledger snapshot for the examination period',
        'Confirm public key matches the registry',
        'Compile operational event logs for the period',
        'Pre-stage CUEC evidence (PIN rotation logs, reconciliation logs, RBAC grants)',
        'Send verifier output bundle to the auditor pre-examination'
      ],
      provide: [
        'Ledger snapshot (encrypted USB or read-only mount)',
        'Tenant public key (PEM)',
        'Tenant ID',
        'Pre-run verifier bundle (PDF + JSON + bundle.tar.gz)',
        'Current control description'
      ]
    },
    relatedDocs: ['first-engagement-guide', 'examiner-quickstart', 'audit-procedures']
  },
  {
    id: 'kickoff',
    name: 'Kickoff Meeting',
    timing: 'T−14 to T−0 days',
    icon: 'mdi-handshake-outline',
    color: 'info',
    summary: 'Walk through scope. Set expectations. Request any final pre-engagement materials.',
    auditor: {
      actions: [
        'Walk through scope with institution stakeholders',
        'Confirm sample population and testing approach',
        'Request final pre-engagement materials',
        'Schedule daily check-ins for field work week',
        'Set communication protocol for findings as they emerge'
      ],
      lookFor: [
        'Tenant_id naming scheme matches institution\'s declared pattern',
        'Institution\'s readiness to support the examination',
        'Any unresolved questions about scope',
        'Whether additional regulators are coordinating (Fed + OCC + CFPB)'
      ]
    },
    auditee: {
      actions: [
        'Brief the chain-operations team on examination scope',
        'Prepare workspace for the auditor (desk, network, snapshot access)',
        'Have CISO, MRM committee chair, and audit committee chair available for follow-up',
        'Stage the institution\'s IR playbook for review'
      ],
      provide: [
        'Final ledger snapshot if updated',
        'Workspace access credentials',
        'Schedule of available stakeholders'
      ]
    },
    relatedDocs: ['first-engagement-guide']
  },
  {
    id: 'fieldwork',
    name: 'Field Work — Physical Audit',
    timing: 'T0 to T+14 days (typical 2–10 days on-site)',
    icon: 'mdi-magnify-scan',
    color: 'warning',
    summary: 'The auditor is on-site (or remote-equivalent). Verifier runs, controls tested, anomalies reviewed.',
    auditor: {
      actions: [
        'Validate verifier binary (verifier-validate.sh)',
        'Run verifier independently against the institution\'s snapshot',
        'Compare to the institution\'s pre-run output',
        'Test CUECs from the institution\'s declared tier',
        'Sample operational events log entries',
        'Review the institution\'s anomaly evaluation records',
        'Walk specific runs (verifier walk) for any incidents in the period',
        'Compare two snapshots (verifier diff) if pre/post changes warrant'
      ],
      lookFor: [
        'Pass rate and anomaly rate per day in the period',
        'Sealing delays and notification compliance (72-hour SHOULD)',
        'Late-binding event rate',
        'Master-key reconciliation cadence and unmatched_count baseline',
        'PIN rotation evidence',
        'IR playbook activation records for chain-detected events',
        'Whether software-key fallback ever appeared in production',
        'Vendor SOC report consumption (BYOC / vendor-hosted)',
        'Institution\'s response to any findings during the field-work week'
      ]
    },
    auditee: {
      actions: [
        'Provide auditor access to systems they need to inspect',
        'Respond to evidence requests promptly',
        'Have chain-operations team available for technical questions',
        'Document any findings discussed during the engagement',
        'Review preliminary findings with the auditor before they finalize'
      ],
      provide: [
        'Operational event log queries on demand',
        'CUEC evidence samples (rotation logs, RBAC grants, IR records)',
        'IR playbook and any activation records',
        'Vendor SOC report (if BYOC or vendor-hosted)',
        'Customer-correlation index entries for any disputed runs'
      ]
    },
    relatedDocs: ['audit-procedures', 'sample-report', 'finding-language', 'cuecs', 'incident-response-playbook']
  },
  {
    id: 'findings',
    name: 'Findings & Discussion',
    timing: 'T+14 to T+45 days',
    icon: 'mdi-clipboard-text-outline',
    color: 'warning',
    summary: 'Auditor drafts findings. Institution responds with management responses and remediation plans.',
    auditor: {
      actions: [
        'Draft findings using examination-report language',
        'Apply severity guidance per failure mode',
        'Identify any repeat findings from prior examinations',
        'Share preliminary findings with institution for factual accuracy',
        'Adjust based on institution\'s factual corrections',
        'Coordinate cross-agency findings if applicable'
      ],
      lookFor: [
        'Institution\'s factual corrections vs interpretive disagreements',
        'Whether the institution remediated any findings during the engagement',
        'Whether findings warrant escalation to enforcement',
        'Pattern across the examination team\'s portfolio (other banks)'
      ]
    },
    auditee: {
      actions: [
        'Review preliminary findings for factual accuracy',
        'Draft management responses with remediation plans',
        'Brief the audit committee on material findings',
        'Engage legal counsel for findings that may have enforcement implications',
        'Coordinate with vendor (BYOC / vendor-hosted) if vendor-side action needed'
      ],
      provide: [
        'Factual corrections to preliminary findings',
        'Management response document',
        'Remediation plans with target dates',
        'Audit committee acknowledgment'
      ]
    },
    relatedDocs: ['finding-language', 'incident-response-playbook']
  },
  {
    id: 'report',
    name: 'Report Issuance',
    timing: 'T+45 to T+90 days',
    icon: 'mdi-file-document-check-outline',
    color: 'success',
    summary: 'Final examination report or SOC opinion issued. Distribution to stakeholders.',
    auditor: {
      actions: [
        'Finalize the examination report (or SOC opinion)',
        'Include findings, management responses, and remediation timelines',
        'For public disclosure (consent orders, formal agreements): coordinate with regulator\'s communications function',
        'Distribute the report through standard regulator channels'
      ],
      lookFor: [
        'Institution\'s consensus on the final report (no factual disagreements)',
        'Adequacy of management responses',
        'Whether the report supports next-cycle examination scoping'
      ]
    },
    auditee: {
      actions: [
        'Receive the final report',
        'Distribute to audit committee and CISO',
        'Update internal control descriptions to reflect findings',
        'Begin implementing remediation per the agreed timeline',
        'For SOC reports: distribute to user entities per engagement letter'
      ],
      provide: [
        'Confirmation of receipt',
        'Audit committee resolution acknowledging findings',
        'Customer-facing communications if findings are public'
      ]
    },
    relatedDocs: ['sample-report', 'finding-language', 'audit-committee-summary']
  },
  {
    id: 'remediation',
    name: 'Remediation & Follow-up',
    timing: 'T+90 days through next examination cycle',
    icon: 'mdi-progress-wrench',
    color: 'success',
    summary: 'Institution implements remediation; auditor tracks progress through next cycle.',
    auditor: {
      actions: [
        'Track open MRAs in the regulator\'s portfolio-management system',
        'Schedule follow-up examinations for material findings',
        'For institutions under enforcement action: monthly or quarterly verifier-output review',
        'Verify remediation effectiveness at next examination',
        'Recommend removal of enforcement action when consistent pass rate is sustained'
      ],
      lookFor: [
        'Verifier output during the remediation period (consistent pass = remediation effective)',
        'Repeat findings (would escalate)',
        'Institution\'s ongoing monitoring of the remediated control',
        'Whether the institution\'s control description was updated to reflect remediation'
      ]
    },
    auditee: {
      actions: [
        'Implement remediation per agreed timeline',
        'Run verifier monthly or quarterly during remediation period',
        'Document evidence of remediation effectiveness',
        'Update control description and IR playbook',
        'Brief audit committee on remediation progress',
        'For enforcement-action removal: assemble cumulative evidence of sustained operation'
      ],
      provide: [
        'Periodic remediation status reports',
        'Verifier output for each monitoring cycle',
        'Updated control description',
        'Audit committee resolutions on remediation milestones'
      ]
    },
    relatedDocs: ['finding-language', 'audit-procedures', 'portfolio-comparison-procedures']
  }
];

// Make available globally
window.CONTENT = {
  CATEGORIES, STAKEHOLDERS, DOCUMENTS, STATS,
  COST_COMPONENTS, COST_TIERS, COST_REDUCTION_OPTIONS,
  THREATS, AUDIT_STAGES
};
