package command_test

import (
	"encoding/json"
	"testing"

	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCache(t *testing.T) cache.Cache {
	t.Helper()

	c := cache.NewMemoryCache()
	// A nil callback reads without populating, so seed with a real one.
	c.Get("alpha").Get("one", func() any { return "secret-value" })
	c.Get("beta").Get("two", func() any { return 42 })

	return c
}

func TestCache_ListAgent_KeysOnly(t *testing.T) {
	c := testCache(t)

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCache(l, c).Execute(t.Context(), parse(t, l, "cache list")))
	})

	var got command.CacheList
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	require.Len(t, got.Namespaces, 2)

	// Sorted by namespace.
	assert.Equal(t, "alpha", got.Namespaces[0].Namespace)
	assert.Equal(t, []string{"one"}, got.Namespaces[0].Keys)
	assert.Equal(t, "beta", got.Namespaces[1].Namespace)

	// Keys only: no cached value may appear anywhere in the payload.
	assert.NotContains(t, out, "secret-value")
}

func TestCache_GetAgent(t *testing.T) {
	c := testCache(t)

	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCache(l, c).Execute(t.Context(), parse(t, l, "cache get alpha one")))
	})

	var got command.CacheGet
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	assert.Equal(t, "alpha", got.Namespace)
	assert.Equal(t, "one", got.Key)
	assert.Equal(t, "secret-value", got.Value)
}

func TestCache_GetAgent_MissingKey(t *testing.T) {
	out := captureAgent(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCache(l, testCache(t)).Execute(t.Context(), parse(t, l, "cache get alpha nope")))
	})

	var got command.CacheGet
	require.NoError(t, json.Unmarshal([]byte(out), &got))

	assert.Nil(t, got.Value, "a missing key yields null rather than an error")
}

// TestCache_ListHuman verifies the non-agent path renders via PTerm rather than
// emitting JSON. PTerm writes through its own cached writer instead of the
// os.Stdout we swap, so its output does not reach the pipe - the assertion is
// that no JSON envelope was produced.
func TestCache_ListHuman(t *testing.T) {
	c := testCache(t)

	out := captureHuman(t, func() {
		l := newTestLogger()
		require.NoError(t, command.NewCache(l, c).Execute(t.Context(), parse(t, l, "cache list")))
	})

	// Results are emitted bare, so there is no envelope to key on: assert the
	// output is not the JSON payload at all.
	assert.NotContains(t, out, `"namespaces"`)
}
