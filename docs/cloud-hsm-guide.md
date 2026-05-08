# Cloud HSM guide

> **What this doc is.** Per-cloud guidance for provisioning and operating an HSM that satisfies the chain's FIPS 140-2 Level 3 (or higher) requirement.

## The conformance bar

The chain spec requires:

- FIPS 140-2 Level 3 or higher (FIPS 140-3 L3 is also acceptable)
- Per-tenant key labels, enforced by HSM ACLs
- Non-extractable signing keys
- Authenticated operator role for signing

These requirements disqualify several common cloud key-management options. The matrix:

| Service | Conformant? | Notes |
|---|---|---|
| AWS CloudHSM Classic | **Yes** | FIPS 140-2 L3 by default |
| AWS KMS Custom Key Stores backed by CloudHSM | **Yes** | Operations performed by CloudHSM; FIPS validation flows through |
| AWS KMS default tier | **No** | FIPS 140-2 L2 only |
| Azure Managed HSM (Premium tier of Key Vault) | **Yes** | FIPS 140-2 L3 |
| Azure Dedicated HSM | **Yes** | Thales Luna under the hood; L3 |
| Azure Key Vault Standard | **No** | FIPS 140-2 L1 |
| Azure Key Vault Premium (legacy SKU) | **Partial** | FIPS 140-2 L2 default; L3 available; consult Azure docs |
| Google Cloud HSM | **Yes** | FIPS 140-2 L3 by default |
| Google Cloud KMS (software keys) | **No** | Software-backed; not L3 |
| On-prem (Thales Luna, Entrust nShield, Utimaco, Yubico YubiHSM 2) | **Yes** | All current models meet L3 |

## AWS CloudHSM

### Provisioning

1. Provision a CloudHSM cluster in the bank's VPC
2. Add at least 2 HSM members for high availability (additional members for higher tiers)
3. Initialize the cluster with the bank's CO (Crypto Officer) credentials
4. Create the seal-job CU (Crypto User) credentials for the chain's seal job

### Key creation

```bash
# Connect to CloudHSM with the seal-job credentials
$ /opt/cloudhsm/bin/cloudhsm-cli login --user seal-job
# Generate the tenant's Ed25519 keypair with non-extractable private key
$ /opt/cloudhsm/bin/cloudhsm-cli key generate-asymmetric-pair \
    --label "tenant-acme-seal" \
    --key-type EC \
    --curve EdDSA \
    --persistent
# Export the public key for the tenant key registry
$ /opt/cloudhsm/bin/cloudhsm-cli key get-public-key --label "tenant-acme-seal"
```

### Cross-region replication

CloudHSM uses backup-and-restore for cross-region key replication:

1. Initiate a backup of the source-region cluster
2. Restore into the destination-region cluster
3. Verify the destination cluster has the same keys with same labels
4. Confirm seal job in destination region can sign with the same keys

The bank's DR plan documents the backup/restore procedure and tests it at least annually.

### Cost

Approximately $13,000/year per HSM cluster member at on-demand pricing. Reserved-instance pricing is lower. A 2-member cluster: $26,000/year per region.

## Azure Managed HSM

### Provisioning

1. Provision a Managed HSM resource in the bank's subscription
2. Activate the HSM with at least 3 RSA security domain keys (the activation procedure)
3. Create the seal-job role assignment with `Crypto User` scope (signing only)

### Key creation

```bash
# Generate Ed25519 key (Managed HSM supports Ed25519 since 2023)
$ az keyvault key create \
    --hsm-name <hsm-name> \
    --name tenant-acme-seal \
    --kty OKP \
    --crv Ed25519 \
    --ops sign verify
# Export public key
$ az keyvault key show --hsm-name <hsm-name> --name tenant-acme-seal --query "key.x" -o tsv
```

### Cross-region

Managed HSM supports geo-redundant configuration. Configure the HSM for the institution's geo-redundancy requirements; the bank's DR plan tests failover.

### Cost

Approximately $15,000/year per HSM partition; first one is the most expensive. Multi-partition deployments scale less than linearly.

## Google Cloud HSM

### Provisioning

1. Enable the Cloud HSM service in the bank's GCP project
2. Create a key ring and a HSM-protected key

### Key creation

```bash
$ gcloud kms keys create tenant-acme-seal \
    --location us-east1 \
    --keyring chain-keys \
    --purpose asymmetric-signing \
    --default-algorithm ec-sign-ed25519 \
    --protection-level hsm
$ gcloud kms keys versions get-public-key 1 \
    --key tenant-acme-seal \
    --keyring chain-keys \
    --location us-east1
```

### Cross-region

Cloud HSM supports cross-region with explicit regional configuration. Bank chooses the regions.

### Cost

Pay-per-use; approximately $0.06 per signing operation, $1/month per key version. Daily-seal-cadence costs are dominated by the per-key-version monthly fee. For a single tenant, approximately $5,000–$15,000/year depending on operations volume.

## On-prem HSMs

### Provisioning

The institution's IT infrastructure team handles. Common products:

- **Thales Luna Network HSM** — typical for large banks; clusters of 2–6 appliances
- **Entrust nShield** — alternative to Thales; similar feature set
- **Utimaco SecurityServer** — common in European institutions
- **Yubico YubiHSM 2** — small-form-factor; suitable for non-production or low-volume

Provisioning, networking, and operational procedures are vendor-specific. The chain's spec is HSM-vendor-neutral; the institution provides the PKCS#11 module path and the seal-job credentials in the ledger config.

### Cost

Approximately $30k–$80k/year all-in: capital amortization + support contracts + ops headcount. Higher upfront cost than cloud; lower per-operation cost at high volume.

## Choosing between options

| Decision | Cloud HSM | On-prem HSM |
|---|---|---|
| Speed of provisioning | Hours | Weeks to months |
| Capital posture | Operating cost | Capital + ongoing support |
| Resilience | Multi-region native | Multi-DC requires explicit cluster |
| Compliance burden | Vendor handles FIPS validation | Institution handles |
| Cost at low volume | Lower | Higher |
| Cost at high volume | Higher | Lower |
| Skill requirements | Cloud admin | HSM specialist |
| Vendor lock-in | Cloud-provider lock-in | HSM-vendor lock-in |

Most institutions running on the cloud choose cloud HSM; institutions with mature HSM operations and significant volume run on-prem.

## Validation

Regardless of provider, the institution's control description records:

- The HSM provider and product
- The FIPS validation certificate number
- The key labels in use per tenant
- The seal-job and admin role separation

The verifier output's signature validation tests confirms the HSM is producing valid signatures; the institution's SOC and examination evidence confirms the FIPS validation.

## Failover and DR

The HSM is the integrity anchor. HSM unavailability stops the seal job, but events continue to be captured and chained (the SDK is independent of the HSM). When the HSM is restored, the seal job catches up.

For multi-region resilience, run an HSM cluster in each region. Per-region keys are isolated; cross-region signing is structurally prevented. The institution's DR plan covers each region.

## Pitfalls

- **Mixing HSM tiers.** A tenant with two HSM keys, one at L3 and one at L2, has the L2 as the weakest link. Don't mix.
- **Sharing keys across tenants.** Per-tenant keys are required by the spec; cross-tenant sharing is non-conformant.
- **Backup keys without replication.** A backup that is not tested for restore is not a backup.
- **PIN management drift.** Without rotation, the PIN ages; a leaked-but-not-rotated PIN is a vulnerability.
