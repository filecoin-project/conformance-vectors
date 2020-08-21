package main

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	. "github.com/filecoin-project/test-vectors/gen/builders"
)

func createActorUnparsableParams(v *Builder) {
	v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

	sender := v.Actors.Account(address.SECP256K1, abi.NewTokenAmount(1_000_000_000_000))
	receiver := v.Actors.Account(address.SECP256K1, abi.NewTokenAmount(1_000_000_000_000))

	v.CommitPreconditions()

	balanceBefore := v.Actors.Balance(sender.Robust)

	// Valid message for construction of a payment channel
	createMsg := v.Messages.Sugar().CreatePaychActor(sender.Robust, receiver.Robust, Value(abi.NewTokenAmount(10_000)), Nonce(0))

	// Form an invalid CBOR payload
	createMsg.Message.Params = append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, createMsg.Message.Params...)

	v.CommitApplies()

	v.Assert.EveryMessageResultSatisfies(ExitCode(exitcode.SysErrSenderInvalid))             // make sure that we get SysErrSenderInvalid error code
	v.Assert.BalanceEq(sender.Robust, big.Sub(balanceBefore, CalculateDeduction(createMsg))) // make sure that gas is deducted from senders's account
}
