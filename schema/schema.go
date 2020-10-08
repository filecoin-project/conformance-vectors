package schema

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
)

// Class represents the type of test vector this instance is.
type Class string

const (
	// ClassMessage tests the VM behaviour and resulting state over one or
	// many messages.
	ClassMessage Class = "message"
	// ClassTipset tests the VM behaviour and resulting state over one or many
	// tipsets and/or null rounds.
	ClassTipset Class = "tipset"
	// ClassBlockSeq tests the state of the system after the arrival of
	// particular blocks at concrete points in time.
	ClassBlockSeq Class = "blockseq"
)

const (
	// HintIncorrect is a standard hint to convey that a vector is knowingly
	// incorrect. Drivers may choose to skip over these vectors, or if it's
	// accompanied by HintNegate, they may perform the assertions as explained
	// in its godoc.
	HintIncorrect = "incorrect"

	// HintNegate is a standard hint to convey to drivers that, if this vector
	// is run, they should negate the postcondition checks (i.e. check that the
	// postcondition state is expressly NOT the one encoded in this vector).
	HintNegate = "negate"
)

// Selector is a predicate the driver can use to determine if this test vector
// is relevant given the capabilities/features of the underlying implementation
// and/or test environment.
type Selector map[string]string

// Metadata provides information on the generation of this test case
type Metadata struct {
	ID      string           `json:"id"`
	Version string           `json:"version,omitempty"`
	Desc    string           `json:"description,omitempty"`
	Comment string           `json:"comment,omitempty"`
	Gen     []GenerationData `json:"gen"`
	Tags    []string         `json:"tags,omitempty"`
}

// GenerationData tags the source of this test case.
type GenerationData struct {
	Source  string `json:"source,omitempty"`
	Version string `json:"version,omitempty"`
}

// StateTree represents a state tree within preconditions and postconditions.
type StateTree struct {
	RootCID cid.Cid `json:"root_cid"`
}

// Base64EncodedBytes is a base64-encoded binary value.
type Base64EncodedBytes []byte

// RandomnessKind specifies the type of randomness that is being requested.
type RandomnessKind string

const (
	RandomnessBeacon = RandomnessKind("beacon")
	RandomnessChain  = RandomnessKind("chain")
)

// RandomnessRule represents a rule to evaluate randomness matches against.
// This encodes to JSON as an array. See godocs on the Randomness type for
// more info.
type RandomnessRule struct {
	Kind                RandomnessKind
	DomainSeparationTag int64
	Epoch               int64
	Entropy             Base64EncodedBytes
}

func (rm RandomnessRule) MarshalJSON() ([]byte, error) {
	array := [4]interface{}{
		rm.Kind,
		rm.DomainSeparationTag,
		rm.Epoch,
		rm.Entropy,
	}
	return json.Marshal(array)
}

func (rm *RandomnessRule) UnmarshalJSON(v []byte) error {
	var (
		arr [4]json.RawMessage
		out RandomnessRule
		err error
	)
	if err = json.Unmarshal(v, &arr); err != nil {
		return err
	}
	if err = json.Unmarshal(arr[0], &out.Kind); err != nil {
		return err
	}
	if err = json.Unmarshal(arr[1], &out.DomainSeparationTag); err != nil {
		return err
	}
	if err = json.Unmarshal(arr[2], &out.Epoch); err != nil {
		return err
	}
	if err = json.Unmarshal(arr[3], &out.Entropy); err != nil {
		return err
	}
	*rm = out
	return nil
}

// Randomness encodes randomness the VM runtime should return while executing
// this vector. It is encoded as a list of ordered rules to match on.
//
// The json serialized form is:
//
//  "randomness": [
//    { "on": ["beacon", 12, 49327, "yxpTbzLhr4uaj7bK0Hl4Vw=="], "ret": "iKyZ2N83N8IoiK2tNJ/H9g==" },
//    { "on": ["chain", 8, 61002, "aacQWICNcMJWtuwTnU+1Hg=="], "ret": "M6HqmihwZ5fXcbQQHhbtsg==" }
//  ]
//
// The four positional values of the `on` array field are:
//
//  1. Kind of randomness (json string; values: beacon, chain).
//  2. Domain separation tag (json number).
//  3. Epoch (json number).
//  4. Entropy (json string; base64 encoded bytes).
//
// When no rules are matched, the driver should return the raw bytes of
// utf-8 string 'i_am_random_____i_am_random_____'
type Randomness []RandomnessMatch

// RandomnessMatch specifies a randomness match. When the implementation
// requests randomness that matches the RandomnessRule in On, Return will
// be returned.
type RandomnessMatch struct {
	On     RandomnessRule     `json:"on"`
	Return Base64EncodedBytes `json:"ret"`
}

// Preconditions contain the environment that needs to be set before the
// vector's applies are applied.
type Preconditions struct {
	// Epoch must be interpreted by the driver as an abi.ChainEpoch in Lotus, or
	// equivalent type in other implementations.
	Epoch     int64      `json:"epoch"`
	StateTree *StateTree `json:"state_tree,omitempty"`

	// BaseFee is an optional base fee to inject into the VM when feeding this
	// message. If absent, it defaults to 100 attoFIL.
	BaseFee *big.Int `json:"basefee,omitempty"`
	// CircSupply is optional. If specified, it is the value that will be
	// injected in the VM when feeding this message. If absent, the default
	// value will be injected (TotalFilecoin, the maximum supply of Filecoin
	// that will ever exist). It is usually odd to set it, and it's only here
	// for specialized vectors.
	CircSupply *big.Int `json:"circ_supply,omitempty"`
}

// Receipt represents a receipt to match against.
type Receipt struct {
	// ExitCode must be interpreted by the driver as an exitcode.ExitCode
	// in Lotus, or equivalent type in other implementations.
	ExitCode    int64              `json:"exit_code"`
	ReturnValue Base64EncodedBytes `json:"return"`
	GasUsed     int64              `json:"gas_used"`
}

// Postconditions contain a representation of VM state at th end of the test
type Postconditions struct {
	ApplyMessageFailures []int      `json:"apply_message_failures,omitempty"`
	StateTree            *StateTree `json:"state_tree"`
	Receipts             []*Receipt `json:"receipts"`
	ReceiptsRoots        []cid.Cid  `json:"receipts_roots,omitempty"`
}

func (b Base64EncodedBytes) String() string {
	return base64.StdEncoding.EncodeToString(b)
}

// MarshalJSON implements json.Marshal for Base64EncodedBytes
func (b Base64EncodedBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON implements json.Unmarshal for Base64EncodedBytes
func (b *Base64EncodedBytes) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}

	if len(s) == 0 {
		*b = nil
		return nil
	}

	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	*b = bytes
	return nil
}

// Diagnostics contain a representation of VM diagnostics
type Diagnostics struct {
	Format string             `json:"format"`
	Data   Base64EncodedBytes `json:"data"`
}

// TestVector is a single test case.
type TestVector struct {
	Class    `json:"class"`
	Selector `json:"selector,omitempty"`

	// Hints are arbitrary flags that convey information to the driver.
	// Use hints to express facts like this vector is knowingly incorrect
	// (e.g. when the reference implementation is broken), or that drivers
	// should negate the postconditions (i.e. test that they are NOT the ones
	// expressed in the vector), etc.
	//
	// Refer to the Hint* constants for common hints.
	Hints []string `json:"hints,omitempty"`

	Meta *Metadata `json:"_meta"`

	// CAR binary data to be loaded into the test environment, usually a CAR
	// containing multiple state trees, addressed by root CID from the relevant
	// objects.
	CAR Base64EncodedBytes `json:"car"`

	// Randomness encodes randomness to be replayed during the execution of this
	// test vector. See godocs on the Randomness type for more info.
	Randomness Randomness `json:"randomness,omitempty"`

	Pre *Preconditions `json:"preconditions"`

	ApplyMessages []Message `json:"apply_messages,omitempty"`
	ApplyTipsets  []Tipset  `json:"apply_tipsets,omitempty"`

	Post        *Postconditions `json:"postconditions"`
	Diagnostics *Diagnostics    `json:"diagnostics,omitempty"`
}

type Message struct {
	Bytes Base64EncodedBytes `json:"bytes"`
	// Epoch must be interpreted by the driver as an abi.ChainEpoch in Lotus, or
	// equivalent type in other implementations.
	Epoch *int64 `json:"epoch,omitempty"`
}

type Tipset struct {
	// Epoch must be interpreted by the driver as an abi.ChainEpoch in Lotus, or
	// equivalent type in other implementations.
	Epoch int64 `json:"epoch"`
	// BaseFee must be interpreted by the driver as an abi.TokenAmount in Lotus,
	// or equivalent type in other implementations.
	BaseFee big.Int `json:"basefee"`
	Blocks  []Block `json:"blocks,omitempty"`
}

type Block struct {
	MinerAddr address.Address      `json:"miner_addr"`
	WinCount  int64                `json:"win_count"`
	Messages  []Base64EncodedBytes `json:"messages"`
}

// Validate validates this test vector against the JSON schema, and applies
// further validation rules that cannot be enforced through JSON Schema.
func (tv TestVector) Validate() error {
	if tv.Class == ClassMessage {
		if len(tv.Post.Receipts) != len(tv.ApplyMessages) {
			return fmt.Errorf("length of postcondition receipts must match length of messages to apply")
		}
	}
	return nil
}

// MustMarshalJSON encodes the test vector to JSON and panics if it errors.
func (tv TestVector) MustMarshalJSON() []byte {
	b, err := json.Marshal(&tv)
	if err != nil {
		panic(err)
	}
	return b
}
