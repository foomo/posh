package util

// Pick returns a when cond is true and b otherwise.
//
// It is the expression form of a two-branch if, for cases where the choice is
// a value rather than control flow:
//
//	check: util.Pick(agent.IsAgentMode(), check.AgentCheck, check.DefaultCheck),
//
// Both arguments are evaluated before the call, so never use it where one
// operand is only valid under the condition - Pick(a.HasIndex(i), a[i], "")
// panics on the very case the guard exists for. Use PickF there, and also when
// either side is expensive or has side effects.
func Pick[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}

// PickF returns a() when cond is true and b() otherwise, evaluating only the
// branch it takes.
//
// Use it where Pick would be unsafe or wasteful - when an operand is only valid
// under the condition, is expensive, or has side effects:
//
//	util.PickF(a.HasIndex(i), func() string { return a[i] }, func() string { return "" })
func PickF[T any](cond bool, a, b func() T) T {
	if cond {
		return a()
	}

	return b()
}
