package keccakf

import (
	"math/rand"
	"testing"

	"github.com/consensys/zkevm-monorepo/prover/crypto/keccak"
	"github.com/consensys/zkevm-monorepo/prover/protocol/compiler/dummy"
	"github.com/consensys/zkevm-monorepo/prover/protocol/ifaces"
	"github.com/consensys/zkevm-monorepo/prover/protocol/wizard"
	"github.com/stretchr/testify/assert"
)

func MakeTestCaseInputOutputModule(maxNumKeccakF int) (
	define wizard.DefineFunc,
	prover func(permTrace keccak.PermTraces) wizard.ProverStep,
) {
	round := 0
	mod := &Module{}
	mod.MaxNumKeccakf = maxNumKeccakF
	mod.state = [5][5]ifaces.Column{}
	define = func(builder *wizard.Builder) {
		comp := builder.CompiledIOP
		mod.lookups = newLookUpTables(comp, maxNumKeccakF)
		mod.declareColumns(comp, round, maxNumKeccakF)
		mod.theta.declareColumn(comp, round, maxNumKeccakF)
		mod.rho.declareColumns(comp, round, maxNumKeccakF)
		mod.piChiIota.declareColumns(comp, round, maxNumKeccakF)
		mod.IO.newInput(comp, round, maxNumKeccakF, *mod)
		mod.IO.newOutput(comp, round, maxNumKeccakF, *mod)
	}

	prover = func(permTrace keccak.PermTraces) wizard.ProverStep {
		return func(run *wizard.ProverRuntime) {
			mod.Assign(run, permTrace)
		}
	}
	return define, prover
}

func TestInputOutputModule(t *testing.T) {
	// #nosec G404 --we don't need a cryptographic RNG for testing purpose
	rng := rand.New(rand.NewSource(0))
	numCases := 15
	maxNumKeccakf := 128
	// The -1 is here to prevent the generation of a padding block
	maxInputSize := maxNumKeccakf*keccak.Rate - 1

	definer, prover := MakeTestCaseInputOutputModule(maxNumKeccakf)
	comp := wizard.Compile(definer, dummy.Compile)

	for i := 0; i < numCases; i++ {
		// Generate a random piece of data
		dataSize := rng.Intn(maxInputSize + 1)
		data := make([]byte, dataSize)
		rng.Read(data)

		// Generate permutation traces for the data
		traces := keccak.PermTraces{}
		keccak.Hash(data, &traces)

		proof := wizard.Prove(comp, prover(traces))
		assert.NoErrorf(t, wizard.Verify(comp, proof), "invalid proof")
	}
}