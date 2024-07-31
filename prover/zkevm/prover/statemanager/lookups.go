package statemanager

import (
	"github.com/consensys/zkevm-monorepo/prover/protocol/ifaces"
	"github.com/consensys/zkevm-monorepo/prover/protocol/wizard"
	mimcCodeHash "github.com/consensys/zkevm-monorepo/prover/zkevm/prover/statemanager/mimccodehash"
	"github.com/consensys/zkevm-monorepo/prover/zkevm/prover/statemanager/statesummary"
)

// lookupStateSummaryCodeHash adds the lookup constraints to ensure the MiMC
// CodeHash module and the StateSummary module contains consistent consistent
// pairs of (MiMCCodeHash, KeccakCodeHash)
func lookupStateSummaryCodeHash(comp *wizard.CompiledIOP,
	accountPeek *statesummary.AccountPeek,
	codeHash *mimcCodeHash.Module) {
	// All the state manager operations are performed in round zero
	round := 0

	// Lookup between code hashes (Keccak and MiMC) between state summary initial account and MiMC code hash module
	comp.InsertInclusionDoubleConditional(round,
		ifaces.QueryID("LOOKUP_MIMC_CODEHASH_INITIAL_ACCOUNT_INTO_STATE_SUMMARY"),
		[]ifaces.Column{codeHash.CodeHashHi, codeHash.CodeHashLo, codeHash.NewState, codeHash.CodeSize},
		[]ifaces.Column{accountPeek.Initial.KeccakCodeHash.Hi, accountPeek.Initial.KeccakCodeHash.Lo, accountPeek.Initial.MiMCCodeHash, accountPeek.Initial.CodeSize},
		codeHash.IsHashEnd,
		accountPeek.Initial.Exists)

	// Lookup between code hashes (Keccak and MiMC) between state summary final account and MiMC code hash module
	comp.InsertInclusionDoubleConditional(round,
		ifaces.QueryIDf("LOOKUP_MIMC_CODEHASH_FINAL_ACCOUNT_INTO_STATE_SUMMARY"),
		[]ifaces.Column{codeHash.CodeHashHi, codeHash.CodeHashLo, codeHash.NewState, codeHash.CodeSize},
		[]ifaces.Column{accountPeek.Final.KeccakCodeHash.Hi, accountPeek.Final.KeccakCodeHash.Lo, accountPeek.Final.MiMCCodeHash, accountPeek.Final.CodeSize},
		codeHash.IsHashEnd,
		accountPeek.Final.Exists)
}