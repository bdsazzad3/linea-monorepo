package innerproduct

import (
	"github.com/consensys/zkevm-monorepo/prover/protocol/coin"
	"github.com/consensys/zkevm-monorepo/prover/protocol/column"
	"github.com/consensys/zkevm-monorepo/prover/protocol/ifaces"
	"github.com/consensys/zkevm-monorepo/prover/protocol/query"
	"github.com/consensys/zkevm-monorepo/prover/protocol/wizard"
	"github.com/consensys/zkevm-monorepo/prover/symbolic"
)

// contextForSize stores the compilation context of a pass of the inner-product
// cryptographic compiler. In particular it stores the protocol items added to
// the protocol by the compilation pass (coins, columns, queries etc...)
// relating to a particular size of column.
type contextForSize struct {

	// Collapsed contains the linear combination of the product pairs covered
	// by the context. Says the compilation context compiles the inner-product
	// of the pairs: (a_i, b_i) for i=0..n, then Collapsed is assigned as
	//
	// Collapsed = \sum_i a_i * b_i * BatchinCoin^i
	Collapsed *symbolic.Expression

	// CollapsedBoard is as Collapsed and stores the ExpressionBoard
	// corresponding to the expression.
	CollapsedBoard symbolic.ExpressionBoard

	// Summation column is built by accumulating the sum of all the sub-product
	// terms.
	Summation ifaces.Column

	// SummationOpening stores the local opening query pointing to the last
	// entry of [Summation]. It is compared to the alleged inner-product values
	// by the verifier to finalize the compilation step.s
	SummationOpening query.LocalOpening
}

// compileForSize applies the compilation step on a range of queries such that
// they all relate to column of the same size. The function expects a non-empty
// list of queries.
//
// It returns the compilation context of the query
func compileForSize(
	comp *wizard.CompiledIOP,
	round int,
	queries []query.InnerProduct,
) *contextForSize {

	var (
		hasMoreThan1Pair = len(queries) > 1 || len(queries[0].Bs) > 1
		size             = queries[0].A.Size()
		ctx              = &contextForSize{}
		// batchingCoin is used to collapse all the inner-product queries
		// into a batched inner-product query so that we only need to
		// commit to a single `Summation` column for all theses.
		batchingCoin coin.Info
	)

	if hasMoreThan1Pair {
		round = round + 1
	}

	ctx.Summation = comp.InsertCommit(
		round+1,
		deriveName[ifaces.ColID]("SUMMATION", comp.SelfRecursionCount),
		size,
	)

	if hasMoreThan1Pair {

		var (
			pairProduct = []*symbolic.Expression{}
		)

		batchingCoin = comp.InsertCoin(
			round+1,
			deriveName[coin.Name]("BATCHING_COIN", comp.SelfRecursionCount),
			coin.Field,
		)

		for _, q := range queries {
			for _, b := range q.Bs {
				pairProduct = append(pairProduct, symbolic.Mul(q.A, b))
			}
		}

		ctx.Collapsed = symbolic.NewPolyEval(batchingCoin.AsVariable(), pairProduct)
		ctx.Collapsed.Board()
	}

	if !hasMoreThan1Pair {
		ctx.Collapsed = symbolic.Mul(queries[0].A, queries[0].Bs[0])
		ctx.CollapsedBoard = ctx.Collapsed.Board()
	}

	// This constraints set the recurrent property of summation
	comp.InsertGlobal(
		round+1,
		deriveName[ifaces.QueryID]("SUMMATION_CONSISTENCY", comp.SelfRecursionCount),
		symbolic.Sub(
			ctx.Summation,
			column.Shift(ctx.Summation, -1),
			ctx.Collapsed,
		),
	)

	// This constraint ensures that summation has the correct initial value
	comp.InsertLocal(
		round+1,
		deriveName[ifaces.QueryID]("SUMMATION_INIT", comp.SelfRecursionCount),
		symbolic.Sub(ctx.Collapsed, ctx.Summation),
	)

	// The opening of the final position of ctx.Summation should be equal to
	// the linear combinations of the alleged openings of the inner-products.
	ctx.SummationOpening = comp.InsertLocalOpening(
		round+1,
		deriveName[ifaces.QueryID]("SUMMATION_END", comp.SelfRecursionCount),
		column.Shift(ctx.Summation, -1),
	)

	comp.RegisterVerifierAction(round+1, &verifierForSize{
		Queries:          queries,
		SummationOpening: ctx.SummationOpening,
		BatchOpening:     batchingCoin,
	})

	return ctx
}