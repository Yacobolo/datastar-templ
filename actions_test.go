package ds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ds "github.com/Yacobolo/datastar-templ"
)

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet(t *testing.T) {
	t.Run("simple url", func(t *testing.T) {
		assert.Equal(t, "@get('/api/updates')", ds.Get("/api/updates"))
	})

	t.Run("format args", func(t *testing.T) {
		assert.Equal(t, "@get('/api/todos/42')", ds.Get("/api/todos/%d", 42))
	})

	t.Run("multiple format args", func(t *testing.T) {
		assert.Equal(t, "@get('/api/users/5/todos/42')", ds.Get("/api/users/%d/todos/%d", 5, 42))
	})

	t.Run("string format arg", func(t *testing.T) {
		assert.Equal(t, "@get('/api/search?q=hello')", ds.Get("/api/search?q=%s", "hello"))
	})

	t.Run("single opt", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled'})",
			ds.Get("/api/updates", ds.Opt("requestCancellation", "disabled")),
		)
	})

	t.Run("multiple opts", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled', contentType: 'json'})",
			ds.Get("/api/updates",
				ds.Opt("requestCancellation", "disabled"),
				ds.Opt("contentType", "json"),
			),
		)
	})

	t.Run("format args and opts", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/todos/42',{openWhenHidden: 'true'})",
			ds.Get("/api/todos/%d", 42, ds.Opt("openWhenHidden", "true")),
		)
	})

	t.Run("raw opt", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{openWhenHidden: true})",
			ds.Get("/api/updates", ds.OptRaw("openWhenHidden", "true")),
		)
	})

	t.Run("mixed opt types", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled', openWhenHidden: true})",
			ds.Get("/api/updates",
				ds.Opt("requestCancellation", "disabled"),
				ds.OptRaw("openWhenHidden", "true"),
			),
		)
	})

	t.Run("raw opt with number", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{retryMaxCount: 10})",
			ds.Get("/api/updates", ds.OptRaw("retryMaxCount", "10")),
		)
	})

	t.Run("raw opt with object", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{filterSignals: {include: /^foo/}})",
			ds.Get("/api/updates", ds.OptRaw("filterSignals", "{include: /^foo/}")),
		)
	})

	t.Run("format args interleaved with opts", func(t *testing.T) {
		// Options can appear anywhere in the variadic args; they're partitioned by type
		assert.Equal(t,
			"@get('/api/users/5/todos/42',{requestCancellation: 'disabled'})",
			ds.Get("/api/users/%d/todos/%d", 5, ds.Opt("requestCancellation", "disabled"), 42),
		)
	})
}

// ---------------------------------------------------------------------------
// All verbs
// ---------------------------------------------------------------------------

func TestAllVerbs(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string, ...any) string
		verb string
	}{
		{"Get", ds.Get, "get"},
		{"Post", ds.Post, "post"},
		{"Put", ds.Put, "put"},
		{"Patch", ds.Patch, "patch"},
		{"Delete", ds.Delete, "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/simple", func(t *testing.T) {
			assert.Equal(t, "@"+tt.verb+"('/api/foo')", tt.fn("/api/foo"))
		})

		t.Run(tt.name+"/format_args", func(t *testing.T) {
			assert.Equal(t, "@"+tt.verb+"('/api/foo/42')", tt.fn("/api/foo/%d", 42))
		})

		t.Run(tt.name+"/with_opt", func(t *testing.T) {
			assert.Equal(t,
				"@"+tt.verb+"('/api/foo',{key: 'val'})",
				tt.fn("/api/foo", ds.Opt("key", "val")),
			)
		})
	}
}

// ---------------------------------------------------------------------------
// Post, Put, Patch, Delete specific tests
// ---------------------------------------------------------------------------

func TestPost(t *testing.T) {
	assert.Equal(t, "@post('/api/workcenters')", ds.Post("/api/workcenters"))
}

func TestPut(t *testing.T) {
	assert.Equal(t, "@put('/api/todos/42')", ds.Put("/api/todos/%d", 42))
}

func TestPatch(t *testing.T) {
	assert.Equal(t, "@patch('/api/workcenters/pagesize')", ds.Patch("/api/workcenters/pagesize"))
}

func TestDelete(t *testing.T) {
	assert.Equal(t, "@delete('/api/todos/42')", ds.Delete("/api/todos/%d", 42))
}

// ---------------------------------------------------------------------------
// Composition with attribute helpers
// ---------------------------------------------------------------------------

func TestActionComposition(t *testing.T) {
	t.Run("init with get and opts", func(t *testing.T) {
		attrs := ds.Init(ds.Get("/api/updates", ds.Opt("requestCancellation", "disabled")))
		require.Len(t, attrs, 1)
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled'})",
			attrs["data-init"],
		)
	})

	t.Run("onclick with post", func(t *testing.T) {
		attrs := ds.OnClick(ds.Post("/api/workcenters"))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@post('/api/workcenters')", attrs["data-on:click"])
	})

	t.Run("oninput with post and debounce", func(t *testing.T) {
		attrs := ds.OnInput(ds.Post("/api/search"), ds.ModDebounce, ds.Ms(300))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@post('/api/search')", attrs["data-on:input__debounce.300ms"])
	})

	t.Run("onchange with patch", func(t *testing.T) {
		attrs := ds.OnChange(ds.Patch("/api/pagesize"))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@patch('/api/pagesize')", attrs["data-on:change"])
	})

	t.Run("init with delay and get", func(t *testing.T) {
		attrs := ds.Init(ds.Get("/api/updates"), ds.ModDelay, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@get('/api/updates')", attrs["data-init__delay.500ms"])
	})
}

// ---------------------------------------------------------------------------
// BuildUpdatesInitURL equivalent
// ---------------------------------------------------------------------------

func TestBuildUpdatesInitURLEquivalent(t *testing.T) {
	// Reproduces the output of handlers.BuildUpdatesInitURL using ds helpers
	url := "/api/workcenters/updates?page=1&sortColumn=title&sortDir=asc"
	result := ds.Get(url, ds.Opt("requestCancellation", "disabled"))
	assert.Equal(t,
		"@get('/api/workcenters/updates?page=1&sortColumn=title&sortDir=asc',{requestCancellation: 'disabled'})",
		result,
	)
}

// ---------------------------------------------------------------------------
// Edge Cases - URL Handling
// ---------------------------------------------------------------------------

func TestActionURLEdgeCases(t *testing.T) {
	t.Run("url with query parameters", func(t *testing.T) {
		result := ds.Get("/api/users?filter=active&sort=name")
		assert.Equal(t, "@get('/api/users?filter=active&sort=name')", result)
	})

	t.Run("url with query params and format args", func(t *testing.T) {
		result := ds.Get("/api/users/%d?filter=active", 42)
		assert.Equal(t, "@get('/api/users/42?filter=active')", result)
	})

	t.Run("url with multiple format args", func(t *testing.T) {
		result := ds.Get("/api/users/%d/posts/%d", 5, 10)
		assert.Equal(t, "@get('/api/users/5/posts/10')", result)
	})

	t.Run("url with string format arg", func(t *testing.T) {
		result := ds.Get("/api/search?q=%s", "test")
		assert.Equal(t, "@get('/api/search?q=test')", result)
	})

	t.Run("url with mixed format args", func(t *testing.T) {
		result := ds.Get("/api/%s/%d", "users", 42)
		assert.Equal(t, "@get('/api/users/42')", result)
	})

	t.Run("relative url", func(t *testing.T) {
		result := ds.Post("/endpoint")
		assert.Equal(t, "@post('/endpoint')", result)
	})

	t.Run("url with path segments", func(t *testing.T) {
		result := ds.Get("/api/v1/users/profile")
		assert.Equal(t, "@get('/api/v1/users/profile')", result)
	})

	t.Run("url with fragment", func(t *testing.T) {
		result := ds.Get("/page#section")
		assert.Equal(t, "@get('/page#section')", result)
	})

	t.Run("url with encoded characters", func(t *testing.T) {
		result := ds.Get("/api/search?q=%s", "hello world")
		assert.Equal(t, "@get('/api/search?q=hello world')", result)
	})

	t.Run("url with special path characters", func(t *testing.T) {
		result := ds.Get("/api/users/john.doe@example.com")
		assert.Equal(t, "@get('/api/users/john.doe@example.com')", result)
	})

	t.Run("absolute url", func(t *testing.T) {
		result := ds.Get("https://api.example.com/data")
		assert.Equal(t, "@get('https://api.example.com/data')", result)
	})

	t.Run("url with port", func(t *testing.T) {
		result := ds.Get("http://localhost:3000/api/data")
		assert.Equal(t, "@get('http://localhost:3000/api/data')", result)
	})

	t.Run("unicode in url", func(t *testing.T) {
		result := ds.Get("/api/用户/%d", 1)
		assert.Equal(t, "@get('/api/用户/1')", result)
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Options
// ---------------------------------------------------------------------------

func TestActionOptionEdgeCases(t *testing.T) {
	t.Run("contentType json", func(t *testing.T) {
		result := ds.Post("/api/data", ds.Opt("contentType", "json"))
		assert.Equal(t, "@post('/api/data',{contentType: 'json'})", result)
	})

	t.Run("contentType form", func(t *testing.T) {
		result := ds.Post("/api/data", ds.Opt("contentType", "form"))
		assert.Equal(t, "@post('/api/data',{contentType: 'form'})", result)
	})

	t.Run("selector option", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("selector", ".target"))
		assert.Equal(t, "@get('/api/data',{selector: '.target'})", result)
	})

	t.Run("selector null", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("selector", "null"))
		assert.Equal(t, "@get('/api/data',{selector: null})", result)
	})

	t.Run("headers option", func(t *testing.T) {
		result := ds.Post("/api/data", ds.OptRaw("headers", "{'X-Csrf-Token': 'abc123'}"))
		assert.Equal(t, "@post('/api/data',{headers: {'X-Csrf-Token': 'abc123'}})", result)
	})

	t.Run("openWhenHidden true", func(t *testing.T) {
		result := ds.Post("/api/data", ds.OptRaw("openWhenHidden", "true"))
		assert.Equal(t, "@post('/api/data',{openWhenHidden: true})", result)
	})

	t.Run("openWhenHidden false", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("openWhenHidden", "false"))
		assert.Equal(t, "@get('/api/data',{openWhenHidden: false})", result)
	})

	t.Run("retry auto", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("retry", "auto"))
		assert.Equal(t, "@get('/api/data',{retry: 'auto'})", result)
	})

	t.Run("retry error", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("retry", "error"))
		assert.Equal(t, "@get('/api/data',{retry: 'error'})", result)
	})

	t.Run("retry always", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("retry", "always"))
		assert.Equal(t, "@get('/api/data',{retry: 'always'})", result)
	})

	t.Run("retry never", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("retry", "never"))
		assert.Equal(t, "@get('/api/data',{retry: 'never'})", result)
	})

	t.Run("requestCancellation auto", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("requestCancellation", "auto"))
		assert.Equal(t, "@get('/api/data',{requestCancellation: 'auto'})", result)
	})

	t.Run("requestCancellation disabled", func(t *testing.T) {
		result := ds.Get("/api/data", ds.Opt("requestCancellation", "disabled"))
		assert.Equal(t, "@get('/api/data',{requestCancellation: 'disabled'})", result)
	})

	t.Run("retryInterval option", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("retryInterval", "2000"))
		assert.Equal(t, "@get('/api/data',{retryInterval: 2000})", result)
	})

	t.Run("retryScaler option", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("retryScaler", "1.5"))
		assert.Equal(t, "@get('/api/data',{retryScaler: 1.5})", result)
	})

	t.Run("retryMaxWaitMs option", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("retryMaxWaitMs", "30000"))
		assert.Equal(t, "@get('/api/data',{retryMaxWaitMs: 30000})", result)
	})

	t.Run("retryMaxCount option", func(t *testing.T) {
		result := ds.Get("/api/data", ds.OptRaw("retryMaxCount", "5"))
		assert.Equal(t, "@get('/api/data',{retryMaxCount: 5})", result)
	})

	t.Run("filterSignals with include", func(t *testing.T) {
		result := ds.Post("/api/data", ds.OptRaw("filterSignals", "{include: /^foo\\./}"))
		assert.Equal(t, "@post('/api/data',{filterSignals: {include: /^foo\\./}})", result)
	})

	t.Run("filterSignals with include and exclude", func(t *testing.T) {
		result := ds.Post("/api/data", ds.OptRaw("filterSignals", "{include: /user/, exclude: /password/}"))
		assert.Equal(t, "@post('/api/data',{filterSignals: {include: /user/, exclude: /password/}})", result)
	})

	t.Run("multiple options together", func(t *testing.T) {
		result := ds.Post("/api/data",
			ds.Opt("contentType", "form"),
			ds.Opt("selector", ".target"),
			ds.Opt("retry", "error"),
			ds.OptRaw("openWhenHidden", "true"),
			ds.OptRaw("retryMaxCount", "3"),
		)
		expected := "@post('/api/data',{contentType: 'form', selector: '.target', retry: 'error', openWhenHidden: true, retryMaxCount: 3})"
		assert.Equal(t, expected, result)
	})

	t.Run("duplicate option keys - last wins", func(t *testing.T) {
		result := ds.Get("/api/data",
			ds.Opt("retry", "auto"),
			ds.Opt("retry", "never"),
		)
		// Both are included - Datastar will handle precedence
		assert.Contains(t, result, "retry: 'auto'")
		assert.Contains(t, result, "retry: 'never'")
	})

	t.Run("mix of Opt and OptRaw", func(t *testing.T) {
		result := ds.Post("/api/data",
			ds.Opt("contentType", "json"),
			ds.OptRaw("openWhenHidden", "true"),
			ds.Opt("retry", "error"),
		)
		expected := "@post('/api/data',{contentType: 'json', openWhenHidden: true, retry: 'error'})"
		assert.Equal(t, expected, result)
	})

	t.Run("option with complex object", func(t *testing.T) {
		result := ds.Post("/api/data", ds.OptRaw("payload", "{user: {name: 'John', age: 30}}"))
		assert.Equal(t, "@post('/api/data',{payload: {user: {name: 'John', age: 30}}})", result)
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Format Args
// ---------------------------------------------------------------------------

func TestActionFormatArgsEdgeCases(t *testing.T) {
	t.Run("single int format arg", func(t *testing.T) {
		result := ds.Get("/api/users/%d", 42)
		assert.Equal(t, "@get('/api/users/42')", result)
	})

	t.Run("single string format arg", func(t *testing.T) {
		result := ds.Get("/api/users/%s", "john")
		assert.Equal(t, "@get('/api/users/john')", result)
	})

	t.Run("multiple int format args", func(t *testing.T) {
		result := ds.Get("/api/users/%d/posts/%d/comments/%d", 1, 2, 3)
		assert.Equal(t, "@get('/api/users/1/posts/2/comments/3')", result)
	})

	t.Run("mixed type format args", func(t *testing.T) {
		result := ds.Get("/api/%s/%d/%s", "users", 42, "profile")
		assert.Equal(t, "@get('/api/users/42/profile')", result)
	})

	t.Run("float format arg", func(t *testing.T) {
		result := ds.Get("/api/products?price=%f", 19.99)
		assert.Equal(t, "@get('/api/products?price=19.990000')", result)
	})

	t.Run("format args interleaved with options", func(t *testing.T) {
		result := ds.Get("/api/users/%d", 42, ds.Opt("retry", "error"))
		assert.Equal(t, "@get('/api/users/42',{retry: 'error'})", result)
	})

	t.Run("multiple format args with multiple options", func(t *testing.T) {
		result := ds.Post("/api/users/%d/posts/%d", 5, 10,
			ds.Opt("contentType", "json"),
			ds.Opt("retry", "always"),
		)
		expected := "@post('/api/users/5/posts/10',{contentType: 'json', retry: 'always'})"
		assert.Equal(t, expected, result)
	})

	t.Run("format args with special characters", func(t *testing.T) {
		result := ds.Get("/api/search?q=%s", "hello world")
		assert.Equal(t, "@get('/api/search?q=hello world')", result)
	})

	t.Run("format args with quotes", func(t *testing.T) {
		result := ds.Get("/api/search?q=%s", "it's")
		assert.Equal(t, "@get('/api/search?q=it's')", result)
	})

	t.Run("zero value format args", func(t *testing.T) {
		result := ds.Get("/api/users/%d", 0)
		assert.Equal(t, "@get('/api/users/0')", result)
	})

	t.Run("negative int format arg", func(t *testing.T) {
		result := ds.Get("/api/offset/%d", -10)
		assert.Equal(t, "@get('/api/offset/-10')", result)
	})

	t.Run("bool format arg", func(t *testing.T) {
		result := ds.Get("/api/toggle?active=%t", true)
		assert.Equal(t, "@get('/api/toggle?active=true')", result)
	})
}
