package wizard

import (
	"github.com/consensys/zkevm-monorepo/prover/utils"
	"github.com/consensys/zkevm-monorepo/prover/utils/collection"
)

/*
In a nutshell, an item is an abstract type that
accounts for the fact that CompiledProtocol
registers various things for different rounds
*/
type ByRoundRegister[ID comparable, DATA any] struct {
	// All the data for each key
	mapping collection.Mapping[ID, DATA]
	// All the IDs for a given round
	byRounds collection.VecVec[ID]
	// Gives the round ID of an entry
	byRoundsIndex collection.Mapping[ID, int]
	// Marks an entry as ignorable (but does not delete it)
	ignored collection.Set[ID]
}

/*
Construct a new round register
*/
func NewRegister[ID comparable, DATA any]() ByRoundRegister[ID, DATA] {
	return ByRoundRegister[ID, DATA]{
		mapping:       collection.NewMapping[ID, DATA](),
		byRounds:      collection.NewVecVec[ID](),
		byRoundsIndex: collection.NewMapping[ID, int](),
		ignored:       collection.NewSet[ID](),
	}
}

/*
Insert for a given round. Will panic if an item
with the same ID has been registered first
*/
func (r *ByRoundRegister[ID, DATA]) AddToRound(round int, id ID, data DATA) {
	r.mapping.InsertNew(id, data)
	r.byRounds.AppendToInner(round, id)
	r.byRoundsIndex.InsertNew(id, round)
}

/*
Returns the list of all the keys ever. The result is returned in
Deterministic order.
*/
func (r *ByRoundRegister[ID, DATA]) AllKeys() []ID {
	res := []ID{}
	for roundID := 0; roundID < r.NumRounds(); roundID++ {
		ids := r.AllKeysAt(roundID)
		res = append(res, ids...)
	}
	return res
}

/*
Returns the list of all keys for a given round. Result has deterministic
order (order of insertion)
*/
func (r *ByRoundRegister[ID, DATA]) AllKeysAt(round int) []ID {
	// Reserve up to the desired length just in case.
	// It is absolutely legitimate to query "too far"
	// this can happens for queries for instance.
	// However, it should not happen for coins.
	r.byRounds.Reserve(round + 1)
	return r.byRounds.MustGet(round)
}

/*
Returns the data for associated to an ID. Panic if not found
*/
func (r *ByRoundRegister[ID, DATA]) Data(id ID) DATA {
	return r.mapping.MustGet(id)
}

/*
Find
*/
func (r *ByRoundRegister[ID, DATA]) Round(id ID) int {
	round, ok := r.byRoundsIndex.TryGet(id)
	if !ok {
		utils.Panic("Could not find entry %v", id)
	}
	return round
}

/*
Panic if the name is not found at the given round
*/
func (r *ByRoundRegister[ID, DATA]) MustBeInRound(round int, id ID) {
	round_, ok := r.byRoundsIndex.TryGet(id)
	if !ok {
		utils.Panic("entry `%v` is not found at all. Was expecting to find it at round %v", id, round)
	}
	if round_ != round {
		utils.Panic("Wrong round, the entry %v was expected to be in round %v but found it in round %v", id, round, round_)
	}
}

/*
Panic if the name is not found at all
*/
func (r *ByRoundRegister[ID, DATA]) MustExists(id ...ID) {
	r.mapping.MustExists(id...)
}

/*
Returns true if all the entry exist
*/
func (r *ByRoundRegister[ID, DATA]) Exists(id ...ID) bool {
	return r.mapping.Exists(id...)
}

/*
Returns the number of rounds
*/
func (r *ByRoundRegister[ID, DATA]) NumRounds() int {
	return r.byRounds.Len()
}

/*
Make sure enough rounds are allocated up to the given length
No-op if enough rounds have been allocated, otherwise, will
reserve as many as necessary.
*/
func (r *ByRoundRegister[ID, DATA]) ReserveFor(newLen int) {
	if r.byRounds.Len() < newLen {
		r.byRounds.Reserve(newLen)
	}
}

/*
Returns all the keys that are not marked as ignored in the structure
*/
func (s *ByRoundRegister[ID, DATA]) AllUnignoredKeys() []ID {
	res := []ID{}
	for r := 0; r < s.NumRounds(); r++ {
		allKeys := s.AllKeysAt(r)
		for _, k := range allKeys {
			if s.IsIgnored(k) {
				continue
			}
			res = append(res, k)
		}
	}
	return res
}

/*
Marks an entry as compiled. Panic if the key is missing from the register.
Returns true if the item was already ignored.
*/
func (r *ByRoundRegister[ID, DATA]) MarkAsIgnored(id ID) bool {
	r.mapping.MustExists(id)
	return r.ignored.Insert(id)
}

/*
Returns if the entry is ignored. Panics if the entry is missing from the
map.
*/
func (r *ByRoundRegister[ID, DATA]) IsIgnored(id ID) bool {
	r.mapping.MustExists(id)
	return r.ignored.Exists(id)
}