package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/testkit/clitest"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

const (
	operationalRootFail = "Status: FAIL\nStep: 11\nReason: operational events log root mismatch — sibling-log contents do not produce sealed root\n"
	emptyTreeRoot       = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

// markOperational labels the last count events operational, as a sink would write them,
// and makes every tag and link again. The lines stay a true chain under the sample key.
func markOperational(t *testing.T, lines [][]byte, count int) [][]byte {
	t.Helper()
	first := len(lines) - count
	return restamp(t, lines, func(n int, whole map[string]any) {
		if n >= first {
			whole["chain_kind"] = "operational"
		}
	})
}

// sealRoots computes both roots from the lines, sharing nothing with the walk but the
// RFC 6962 hash functions: §4.2 leaves over payload_hash, §4.3 leaves over canonical bytes.
func sealRoots(t *testing.T, lines [][]byte) (captured, operational string) {
	t.Helper()
	var capturedLeaves, operationalLeaves [][]byte
	for _, line := range lines[1:] {
		whole := decodeLine(t, line)
		stamp := whole["audit_chain"].(map[string]any)
		delete(whole, "audit_chain")
		if whole["chain_kind"] == "operational" {
			canonical, err := jcs.Canonicalize(whole)
			if err != nil {
				t.Fatal(err)
			}
			operationalLeaves = append(operationalLeaves, verify.MerkleLeafHash(canonical))
			continue
		}
		tag, err := base64.StdEncoding.DecodeString(stamp["payload_hash"].(string))
		if err != nil {
			t.Fatal(err)
		}
		capturedLeaves = append(capturedLeaves, verify.MerkleLeafHash(tag))
	}
	return hex.EncodeToString(verify.MerkleTreeHash(capturedLeaves)), hex.EncodeToString(verify.MerkleTreeHash(operationalLeaves))
}

// v1_0cSeal signs a v1.0c seal over these roots, as the seal job would.
func v1_0cSeal(t *testing.T, lines [][]byte, captured, operational string) string {
	t.Helper()
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	for k, v := range goodV1_0b() {
		seal[k] = v
	}
	seal["sign_payload_version"] = "v1.0c"
	seal["merkle_root_hex"] = captured
	seal["operational_events_log_root_hex"] = operational
	delete(seal, "key_versions")
	return writeSeal(t, resign(t, seal, institutionKey()))
}

func TestV1_0c_CorrectSealPassesWithOperationalEntries(t *testing.T) {
	lines := markOperational(t, sampleLines(t), 3)
	captured, operational := sealRoots(t, lines)
	withKey(writeLines(t, lines), v1_0cSeal(t, lines, captured, operational), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutLacks(t, "not walked")
}

// The silent cut. Every remaining link is good, and step 10 excludes these entries, so
// only the root the seal bound can show them gone.
func TestV1_0c_CuttingEveryOperationalEntryAtTheTailFails(t *testing.T) {
	lines := markOperational(t, sampleLines(t), 3)
	captured, operational := sealRoots(t, lines)
	seal := v1_0cSeal(t, lines, captured, operational)
	cut := lines[:len(lines)-3]
	withKey(writeLines(t, cut), seal, samplePubKey).MustExit(t, 1).StdoutContains(t, operationalRootFail)
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, cut), "--seal", seal, "--root-key", samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, operationalRootFail)
}

func TestV1_0c_CuttingSomeOperationalEntriesFails(t *testing.T) {
	lines := markOperational(t, sampleLines(t), 3)
	captured, operational := sealRoots(t, lines)
	withKey(writeLines(t, lines[:len(lines)-1]), v1_0cSeal(t, lines, captured, operational), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, operationalRootFail)
}

func TestV1_0c_TrueEmptyDayPasses(t *testing.T) {
	lines := sampleLines(t)
	captured, operational := sealRoots(t, lines)
	if operational != emptyTreeRoot {
		t.Fatalf("empty tree root is %s", operational)
	}
	withKey(sampleFile, v1_0cSeal(t, lines, captured, operational), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutLacks(t, "not walked")
}

func TestV1_0c_BoundRootWithZeroEntriesFails(t *testing.T) {
	lines := sampleLines(t)
	captured, _ := sealRoots(t, lines)
	withKey(sampleFile, v1_0cSeal(t, lines, captured, strings.Repeat("cd", sha256.Size)), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, operationalRootFail)
}

// A key holder relabels one captured entry and makes new tags. Steps 1-9 pass. Its bytes
// leave one tree and join the other, so step 10 stops it.
func TestV1_0c_RelabelingAnEntryFailsAtTheRoot(t *testing.T) {
	lines := markOperational(t, sampleLines(t), 3)
	captured, operational := sealRoots(t, lines)
	seal := v1_0cSeal(t, lines, captured, operational)
	relabeled := markOperational(t, lines, 4)
	withKey(writeLines(t, relabeled), seal, samplePubKey).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

// Below v1.0c there is no second tree: the seal's root covers every entry.
func TestBelowV1_0c_RootCoversOperationalEntries(t *testing.T) {
	lines := markOperational(t, sampleLines(t), 3)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	withKey(writeLines(t, lines), writeSeal(t, seal), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n")
}
