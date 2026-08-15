package util_test

import (
	"testing"

	"github.com/foomo/posh/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestPick(t *testing.T) {
	assert.Equal(t, "a", util.Pick(true, "a", "b"))
	assert.Equal(t, "b", util.Pick(false, "a", "b"))

	// Instantiates over a non-comparable type, as the check.Check use does.
	called := ""
	fn := util.Pick(true,
		func() { called = "a" },
		func() { called = "b" },
	)
	fn()

	assert.Equal(t, "a", called)
}

func TestPickF(t *testing.T) {
	assert.Equal(t, "a", util.PickF(true, func() string { return "a" }, func() string { return "b" }))
	assert.Equal(t, "b", util.PickF(false, func() string { return "a" }, func() string { return "b" }))
}

func TestPickF_OnlyEvaluatesTakenBranch(t *testing.T) {
	var evaluated []string

	got := util.PickF(true,
		func() string { evaluated = append(evaluated, "a"); return "a" },
		func() string { evaluated = append(evaluated, "b"); return "b" },
	)

	assert.Equal(t, "a", got)
	assert.Equal(t, []string{"a"}, evaluated, "the untaken branch must not run")
}

// TestPickF_UnsafeForPick covers the case Pick cannot express: the false branch
// would panic if evaluated eagerly.
func TestPickF_UnsafeForPick(t *testing.T) {
	var empty []string

	assert.NotPanics(t, func() {
		got := util.PickF(len(empty) > 0,
			func() string { return empty[0] },
			func() string { return "" },
		)

		assert.Empty(t, got)
	})
}
