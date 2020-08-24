package main

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	. "github.com/filecoin-project/test-vectors/gen/builders"
	. "github.com/filecoin-project/test-vectors/schema"
)

func main() {
	g := NewGenerator()
	defer g.Wait()

	g.MessageVectorGroup("addresses",
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "sequential-10",
				Version: "v1",
				Desc:    "actor addresses are sequential",
			},
			Func: sequentialAddresses,
		},
	)

	g.MessageVectorGroup("on_transfer",
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "ok-create-secp256k1",
				Version: "v1",
			},
			Func: actorCreationOnTransfer(actorCreationOnTransferParams{
				senderType:   address.SECP256K1,
				senderBal:    abi.NewTokenAmount(1_000_000_000_000_000),
				receiverAddr: MustNewSECP256K1Addr("publickeyfoo"),
				amount:       abi.NewTokenAmount(10_000),
				exitCode:     exitcode.Ok,
			}),
		},
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "ok-create-bls",
				Version: "v1",
			},
			Func: actorCreationOnTransfer(actorCreationOnTransferParams{
				senderType:   address.SECP256K1,
				senderBal:    abi.NewTokenAmount(1_000_000_000_000_000),
				receiverAddr: MustNewBLSAddr(1),
				amount:       abi.NewTokenAmount(10_000),
				exitCode:     exitcode.Ok,
			}),
		},
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "fail-secp256k1-insufficient-balance",
				Version: "v1",
			},
			Func: actorCreationOnTransfer(actorCreationOnTransferParams{
				senderType:   address.SECP256K1,
				senderBal:    abi.NewTokenAmount(9_999),
				receiverAddr: MustNewSECP256K1Addr("publickeyfoo"),
				amount:       abi.NewTokenAmount(10_000),
				exitCode:     exitcode.SysErrSenderStateInvalid,
			}),
		},
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "fail-bls-insufficient-balance",
				Version: "v1",
			},
			Func: actorCreationOnTransfer(actorCreationOnTransferParams{
				senderType:   address.SECP256K1,
				senderBal:    abi.NewTokenAmount(9_999),
				receiverAddr: MustNewBLSAddr(1),
				amount:       abi.NewTokenAmount(10_000),
				exitCode:     exitcode.SysErrSenderStateInvalid,
			}),
		},
	)

	g.MessageVectorGroup("params",
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "fails-unparsable-init-actor-exec-msg",
				Version: "v1",
				Desc:    "verifies that actor creation fails and gas is deducted when passing unparsable init exec message",
				Comment: "this should not return SysErrSenderInvalid; it should return something else, likely an SysErrSerialization",
			},
			Mode:  ModeLenientAssertions,
			Hints: []string{HintIncorrect, HintNegate},
			Func:  createActorInitExecUnparsableParams,
		},
		&MessageVectorGenItem{
			Metadata: &Metadata{
				ID:      "fails-unparsable-constructor-params-via-init-actor",
				Version: "v1",
				Desc:    "verifies that actor creation fails and gas is deducted when passing unparsable constructor params via init actor",
				Comment: "this should not return SysErrSenderInvalid; it should return something else, likely an ErrSerialization because the error is in actor space",
			},
			Mode:  ModeLenientAssertions,
			Hints: []string{HintIncorrect, HintNegate},
			Func:  createActorCtorUnparsableParamsViaInitExec,
		},
	)
}
