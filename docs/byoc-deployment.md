# BYOC deployment guide

> **What this doc is.** Concrete deployment guide for the Bring-Your-Own-Cloud topology. The bank operates the vendor's image in the bank's cloud account; the boundary between bank and vendor is enforced by cloud IAM and network controls.

## When to use BYOC

Pick BYOC when:

- The institution wants vendor implementation expertise but cannot have the vendor hold the data or the keys
- Compliance posture requires the institution to hold its own master key
- The institution operates in a multi-cloud environment where vendor-hosted is impractical
- The institution's procurement framework prefers software licensing over service licensing

If the institution wants a fully managed service, use vendor-hosted instead. If the institution wants full control with no vendor relationship, use self-hosted.

## High-level architecture

```
[ Bank cloud account ]
    └── [ Bank-controlled VPC ]
        ├── [ Vendor's ledger image, running under bank IAM ]
        ├── [ HSM in private subnet, bank-controlled ]
        ├── [ Postgres for WAL + hot store, bank-controlled ]
        ├── [ S3 bucket for cold store, bank-controlled ]
        └── [ Bank-controlled secret manager for HSM PIN ]
[ Bank-controlled mirror registry ]
    └── [ Vendor signed images, pulled through this mirror ]
[ Vendor's cloud account ]
    └── [ Read-only support telemetry, bank-mediated ]
```

## IAM permission matrix

The control: **the vendor must not be able to reach the bank's HSM, master key, or event-payload data.**

Below: who can do what.

| Resource | Bank's chain-ops role | Vendor's support role | Vendor's image (running in bank account) |
|---|---|---|---|
| HSM signing operations | Yes (via seal-job role only) | No | Yes (via seal-job role only; runs as the image's IAM identity) |
| HSM administration (extract, import, delete) | Yes (via HSM admin role; separate person) | No | No |
| Bank master-key custodian | Yes | No | No |
| Postgres WAL/hot store, write | No (read for ops) | No | Yes (the image is the writer) |
| Postgres WAL/hot store, read | Yes | No | Yes |
| S3 cold store, write | No (read for ops) | No | Yes |
| S3 cold store, read | Yes | No | Yes |
| Bank secret manager (HSM PIN) | Yes (chain-ops role) | No | Yes (read-only at process start; one secret only) |
| Bank's event payload data | Yes (read for ops/audit) | **No** | Yes (the image processes it) |
| Vendor support telemetry endpoint | No | Yes (receive) | Yes (send, after redaction) |
| Bank's IAM control plane | Yes | No | No |

The vendor's support role has access only to the redacted telemetry stream the bank chooses to share. The vendor's image runs under a bank-controlled IAM identity that the bank can revoke unilaterally.

## Deployment steps

### Step 1 — Provision bank-side resources

1. Create a dedicated VPC for the chain workload
2. Provision the HSM in a private subnet with private endpoints; no internet egress
3. Provision Postgres (RDS, Aurora, etc.) in a private subnet
4. Provision the S3 bucket for cold storage with bucket policy denying access from outside the bank's IAM principals
5. Configure the bank's secret manager (Secrets Manager, Key Vault) with the HSM PIN

### Step 2 — Provision the bank-controlled mirror registry

The mirror registry holds vendor-published images, signed and scanned, ready for deployment:

1. Provision a registry (ECR, ACR, GAR, or Harbor)
2. Configure replication-from-vendor pull rules with cosign signature verification
3. Configure vulnerability scanning (Trivy, Snyk, or equivalent) to run on every replicated image
4. Configure the bank's deployment pipeline to pull from this mirror, never from the vendor's registry

### Step 3 — Provision the vendor image's IAM role

1. Create a new IAM role: `chain-ledger-runtime`
2. Attach permission policies for:
   - HSM signing (seal-job operations only)
   - Postgres connection (chain-ledger user only)
   - S3 cold-store bucket (read/write to specific prefix)
   - Bank secret manager (read HSM PIN secret only)
   - Egress to vendor-support telemetry endpoint (post-redaction)
3. Set a strict trust policy: only the bank-controlled compute (EKS service account, EC2 instance profile) can assume this role

### Step 4 — Deploy the vendor image

1. Pull the vendor image from the bank's mirror
2. Verify the image's cosign signature against the cached project public key
3. Deploy to the bank-controlled compute, attaching `chain-ledger-runtime` as the runtime role
4. Configure the image with bank-controlled endpoints (Postgres DSN, HSM endpoint, S3 prefix, secret manager path)
5. Start the workload; verify health probes pass

### Step 5 — Configure support telemetry routing

The vendor needs operational telemetry. Telemetry flows through a bank-controlled router:

1. The vendor image emits operational events and metrics to a bank-controlled OTel Collector
2. The collector applies the bank's documented redaction policy (no event payload data; control-plane events only)
3. The collector forwards the redacted stream to the vendor's support endpoint
4. The bank's privacy-impact assessment for the support relationship documents the data flow
5. The bank can revoke the egress unilaterally

## Network controls

| Control | Implementation |
|---|---|
| Vendor cannot reach HSM | HSM is in a private subnet with private endpoints; no internet egress; security group denies all inbound except from `chain-ledger-runtime` workload |
| Vendor cannot read event data from outside bank | Postgres and S3 are in private network; vendor support role has no IAM permissions on either |
| Image-pull cannot be in-flight-tampered | All pulls go through bank's mirror registry, which validates cosign signatures before storing |
| Vendor cannot exfiltrate event data via vendor-support telemetry | Bank-controlled OTel Collector is the only egress for support telemetry; redaction is enforced at the collector |

## Audit considerations

The bank's SOC team and FFIEC examiner test the BYOC controls:

- **CC6.1 (Logical access).** IAM permission matrix; vendor cannot assume the runtime role.
- **CC6.7 (Restricts movement of data).** Bank-controlled OTel Collector; redaction policy.
- **CC6.8 (Prevents/detects unauthorized software).** Mirror-registry signature validation; cosign on every pull.
- **CC9.2 (Vendor management).** Bank-controlled IAM; bank can revoke vendor access unilaterally.

## Common pitfalls

- **Vendor pulls credentials are compromised.** Defense: bank-controlled mirror registry; bank pulls from vendor's registry; vendor never has credentials to push directly to the bank's compute.
- **Vendor's image is over-permissioned.** Defense: principle of least privilege on `chain-ledger-runtime`; remove permissions the runtime doesn't need.
- **HSM PIN leaks via vendor support telemetry.** Defense: strict redaction policy; the PIN is never in any log or trace; the redaction is tested.
- **Vendor's support role drifts to broader access.** Defense: reviews of the vendor support role on the bank's standard third-party-access cadence (typically quarterly).

## When the vendor relationship ends

If the bank decides to leave the vendor:

1. Revoke `chain-ledger-runtime` permissions for vendor support role
2. Stop pulling new images from the vendor's registry
3. Migrate to a different vendor's image (bank's data and keys remain bank-controlled; only the runtime image changes)
4. Verify continuity: the new image processes events correctly under the same chain spec; past events remain verifiable under the original public key

The chain's open-source standard plus the bank-controlled keys and data make vendor switching tractable.
