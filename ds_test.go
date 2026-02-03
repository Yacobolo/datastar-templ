package ds_test

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ds "github.com/Yacobolo/datastar-templ"
)

// ---------------------------------------------------------------------------
// Modifier helpers
// ---------------------------------------------------------------------------

func TestDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want ds.Modifier
	}{
		{"zero", 0, ".0ms"},
		{"1ms", time.Millisecond, ".1ms"},
		{"500ms", 500 * time.Millisecond, ".500ms"},
		{"1s", time.Second, ".1000ms"},
		{"rounds 500us up to 1ms", 500 * time.Microsecond, ".1ms"},
		{"rounds 100us down to 0ms", 100 * time.Microsecond, ".0ms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Duration(tt.d))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Duration(-1) })
	})
}

func TestMs(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want ds.Modifier
	}{
		{"zero", 0, ".0ms"},
		{"300", 300, ".300ms"},
		{"1000", 1000, ".1000ms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Ms(tt.n))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Ms(-1) })
	})
}

func TestSeconds(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want ds.Modifier
	}{
		{"zero", 0, ".0s"},
		{"1", 1, ".1s"},
		{"5", 5, ".5s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Seconds(tt.n))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Seconds(-1) })
	})
}

func TestThreshold(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want ds.Modifier
	}{
		{"0.25", 0.25, ".25"},
		{"0.50", 0.50, ".50"},
		{"0.75", 0.75, ".75"},
		{"1.0 (100%)", 1.0, ".100"},
		{"rounds 0.333 to .33", 0.333, ".33"},
		{"rounds 0.335 to .34", 0.335, ".34"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Threshold(tt.val))
		})
	}

	t.Run("panics on zero", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(0) })
	})
	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(-0.1) })
	})
	t.Run("panics on > 1", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(1.1) })
	})
}

// ---------------------------------------------------------------------------
// Merge
// ---------------------------------------------------------------------------

func TestMerge(t *testing.T) {
	t.Run("combines multiple", func(t *testing.T) {
		result := ds.Merge(
			ds.Show("$visible"),
			ds.OnClick("$open = false"),
		)
		assert.Equal(t, "$visible", result["data-show"])
		assert.Equal(t, "$open = false", result["data-on:click"])
		assert.Len(t, result, 2)
	})

	t.Run("later overrides earlier", func(t *testing.T) {
		result := ds.Merge(
			ds.Show("$first"),
			ds.Show("$second"),
		)
		assert.Equal(t, "$second", result["data-show"])
		assert.Len(t, result, 1)
	})

	t.Run("empty merge", func(t *testing.T) {
		result := ds.Merge()
		assert.Empty(t, result)
	})
}

// ---------------------------------------------------------------------------
// DOM Event Functions – data-on:{event}
// ---------------------------------------------------------------------------

func TestOnClick(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		mods    []ds.Modifier
		wantKey string
		wantVal string
	}{
		{
			name:    "simple",
			expr:    "$foo = true",
			wantKey: "data-on:click",
			wantVal: "$foo = true",
		},
		{
			name:    "with debounce",
			expr:    "$foo = true",
			mods:    []ds.Modifier{ds.ModDebounce, ds.Ms(300)},
			wantKey: "data-on:click__debounce.300ms",
			wantVal: "$foo = true",
		},
		{
			name:    "with window and throttle",
			expr:    "$x = 1",
			mods:    []ds.Modifier{ds.ModWindow, ds.ModThrottle, ds.Ms(100), ds.NoLeading},
			wantKey: "data-on:click__window__throttle.100ms.noleading",
			wantVal: "$x = 1",
		},
		{
			name:    "with prevent",
			expr:    "doStuff()",
			mods:    []ds.Modifier{ds.ModPrevent},
			wantKey: "data-on:click__prevent",
			wantVal: "doStuff()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := ds.OnClick(tt.expr, tt.mods...)
			require.Len(t, attrs, 1)
			assert.Equal(t, tt.wantVal, attrs[tt.wantKey])
		})
	}
}

// TestAllDOMEvents verifies every typed DOM event function produces the
// correct "data-on:{event}" attribute name.
func TestAllDOMEvents(t *testing.T) {
	type eventFunc = func(string, ...ds.Modifier) templ.Attributes

	tests := []struct {
		name    string
		fn      eventFunc
		wantKey string
	}{
		// Mouse
		{"OnClick", ds.OnClick, "data-on:click"},
		{"OnDblClick", ds.OnDblClick, "data-on:dblclick"},
		{"OnMouseDown", ds.OnMouseDown, "data-on:mousedown"},
		{"OnMouseUp", ds.OnMouseUp, "data-on:mouseup"},
		{"OnMouseOver", ds.OnMouseOver, "data-on:mouseover"},
		{"OnMouseOut", ds.OnMouseOut, "data-on:mouseout"},
		{"OnMouseMove", ds.OnMouseMove, "data-on:mousemove"},
		{"OnMouseEnter", ds.OnMouseEnter, "data-on:mouseenter"},
		{"OnMouseLeave", ds.OnMouseLeave, "data-on:mouseleave"},
		{"OnContextMenu", ds.OnContextMenu, "data-on:contextmenu"},
		// Keyboard
		{"OnKeyDown", ds.OnKeyDown, "data-on:keydown"},
		{"OnKeyUp", ds.OnKeyUp, "data-on:keyup"},
		{"OnKeyPress", ds.OnKeyPress, "data-on:keypress"},
		// Focus
		{"OnFocus", ds.OnFocus, "data-on:focus"},
		{"OnBlur", ds.OnBlur, "data-on:blur"},
		{"OnFocusIn", ds.OnFocusIn, "data-on:focusin"},
		{"OnFocusOut", ds.OnFocusOut, "data-on:focusout"},
		// Form
		{"OnSubmit", ds.OnSubmit, "data-on:submit"},
		{"OnReset", ds.OnReset, "data-on:reset"},
		{"OnInput", ds.OnInput, "data-on:input"},
		{"OnChange", ds.OnChange, "data-on:change"},
		{"OnInvalid", ds.OnInvalid, "data-on:invalid"},
		{"OnSelect", ds.OnSelect, "data-on:select"},
		// Drag
		{"OnDrag", ds.OnDrag, "data-on:drag"},
		{"OnDragStart", ds.OnDragStart, "data-on:dragstart"},
		{"OnDragEnd", ds.OnDragEnd, "data-on:dragend"},
		{"OnDragOver", ds.OnDragOver, "data-on:dragover"},
		{"OnDragEnter", ds.OnDragEnter, "data-on:dragenter"},
		{"OnDragLeave", ds.OnDragLeave, "data-on:dragleave"},
		{"OnDrop", ds.OnDrop, "data-on:drop"},
		// Touch
		{"OnTouchStart", ds.OnTouchStart, "data-on:touchstart"},
		{"OnTouchEnd", ds.OnTouchEnd, "data-on:touchend"},
		{"OnTouchMove", ds.OnTouchMove, "data-on:touchmove"},
		{"OnTouchCancel", ds.OnTouchCancel, "data-on:touchcancel"},
		// Pointer
		{"OnPointerDown", ds.OnPointerDown, "data-on:pointerdown"},
		{"OnPointerUp", ds.OnPointerUp, "data-on:pointerup"},
		{"OnPointerMove", ds.OnPointerMove, "data-on:pointermove"},
		{"OnPointerOver", ds.OnPointerOver, "data-on:pointerover"},
		{"OnPointerOut", ds.OnPointerOut, "data-on:pointerout"},
		{"OnPointerEnter", ds.OnPointerEnter, "data-on:pointerenter"},
		{"OnPointerLeave", ds.OnPointerLeave, "data-on:pointerleave"},
		{"OnPointerCancel", ds.OnPointerCancel, "data-on:pointercancel"},
		{"OnGotPointerCapture", ds.OnGotPointerCapture, "data-on:gotpointercapture"},
		{"OnLostPointerCapture", ds.OnLostPointerCapture, "data-on:lostpointercapture"},
		// Scroll/Wheel
		{"OnScroll", ds.OnScroll, "data-on:scroll"},
		{"OnWheel", ds.OnWheel, "data-on:wheel"},
		// Animation/Transition
		{"OnAnimationStart", ds.OnAnimationStart, "data-on:animationstart"},
		{"OnAnimationEnd", ds.OnAnimationEnd, "data-on:animationend"},
		{"OnAnimationIteration", ds.OnAnimationIteration, "data-on:animationiteration"},
		{"OnTransitionEnd", ds.OnTransitionEnd, "data-on:transitionend"},
		// Media
		{"OnLoad", ds.OnLoad, "data-on:load"},
		{"OnError", ds.OnError, "data-on:error"},
		// Clipboard
		{"OnCopy", ds.OnCopy, "data-on:copy"},
		{"OnCut", ds.OnCut, "data-on:cut"},
		{"OnPaste", ds.OnPaste, "data-on:paste"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := tt.fn("expr")
			require.Len(t, attrs, 1)
			assert.Equal(t, "expr", attrs[tt.wantKey])
		})
	}
}

func TestOnEvent(t *testing.T) {
	t.Run("custom event", func(t *testing.T) {
		attrs := ds.OnEvent("table-select", "$selected = evt.detail.ids")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$selected = evt.detail.ids", attrs["data-on:table-select"])
	})

	t.Run("custom event with modifiers", func(t *testing.T) {
		attrs := ds.OnEvent("my-event", "handle()", ds.ModOnce, ds.ModCapture)
		require.Len(t, attrs, 1)
		assert.Equal(t, "handle()", attrs["data-on:my-event__once__capture"])
	})
}

func TestOnEventModifiers(t *testing.T) {
	t.Run("debounce with duration and leading", func(t *testing.T) {
		attrs := ds.OnInput("doSearch()", ds.ModDebounce, ds.Duration(500*time.Millisecond), ds.Leading)
		require.Len(t, attrs, 1)
		assert.Equal(t, "doSearch()", attrs["data-on:input__debounce.500ms.leading"])
	})

	t.Run("window outside", func(t *testing.T) {
		attrs := ds.OnClick("close()", ds.ModWindow, ds.ModOutside)
		require.Len(t, attrs, 1)
		assert.Equal(t, "close()", attrs["data-on:click__window__outside"])
	})

	t.Run("once passive capture", func(t *testing.T) {
		attrs := ds.OnScroll("track()", ds.ModOnce, ds.ModPassive, ds.ModCapture)
		require.Len(t, attrs, 1)
		assert.Equal(t, "track()", attrs["data-on:scroll__once__passive__capture"])
	})

	t.Run("stop prevent", func(t *testing.T) {
		attrs := ds.OnSubmit("save()", ds.ModStop, ds.ModPrevent)
		require.Len(t, attrs, 1)
		assert.Equal(t, "save()", attrs["data-on:submit__stop__prevent"])
	})

	t.Run("view transition", func(t *testing.T) {
		attrs := ds.OnClick("navigate()", ds.ModViewTransition)
		require.Len(t, attrs, 1)
		assert.Equal(t, "navigate()", attrs["data-on:click__viewtransition"])
	})

	t.Run("case modifier", func(t *testing.T) {
		attrs := ds.OnEvent("myEvent", "handle()", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, "handle()", attrs["data-on:myEvent__case.kebab"])
	})
}

// ---------------------------------------------------------------------------
// Datastar Attributes
// ---------------------------------------------------------------------------

func TestShow(t *testing.T) {
	attrs := ds.Show("$visible")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$visible", attrs["data-show"])
}

func TestText(t *testing.T) {
	attrs := ds.Text("$count + ' items'")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$count + ' items'", attrs["data-text"])
}

func TestEffect(t *testing.T) {
	attrs := ds.Effect("$total = $price * $qty")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$total = $price * $qty", attrs["data-effect"])
}

func TestInit(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Init("$count = 1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count = 1", attrs["data-init"])
	})

	t.Run("with delay", func(t *testing.T) {
		attrs := ds.Init("$count = 1", ds.ModDelay, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count = 1", attrs["data-init__delay.500ms"])
	})

	t.Run("with view transition", func(t *testing.T) {
		attrs := ds.Init("navigate()", ds.ModViewTransition)
		require.Len(t, attrs, 1)
		assert.Equal(t, "navigate()", attrs["data-init__viewtransition"])
	})
}

func TestBind(t *testing.T) {
	t.Run("keyed syntax", func(t *testing.T) {
		attrs := ds.Bind("name")
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:name"])
	})

	t.Run("nested signal", func(t *testing.T) {
		attrs := ds.Bind("table.search")
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:table.search"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Bind("my-field", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:my-field__case.kebab"])
	})
}

func TestBindExpr(t *testing.T) {
	attrs := ds.BindExpr("name")
	require.Len(t, attrs, 1)
	assert.Equal(t, "name", attrs["data-bind"])
}

func TestSignals(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1})
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals"])
	})

	t.Run("with ifmissing modifier", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1}, ds.ModIfMissing)
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals__ifmissing"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1}, ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals__case.kebab"])
	})
}

func TestSignalsJSON(t *testing.T) {
	attrs := ds.SignalsJSON(`{foo: {bar: 1}}`)
	require.Len(t, attrs, 1)
	assert.Equal(t, `{foo: {bar: 1}}`, attrs["data-signals"])
}

func TestSignalKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.SignalKey("foo", "1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "1", attrs["data-signals:foo"])
	})

	t.Run("with ifmissing", func(t *testing.T) {
		attrs := ds.SignalKey("foo", "1", ds.ModIfMissing)
		require.Len(t, attrs, 1)
		assert.Equal(t, "1", attrs["data-signals:foo__ifmissing"])
	})
}

func TestComputed(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		attrs := ds.Computed("total", "$price * $qty")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'total': () => $price * $qty}", attrs["data-computed"])
	})

	t.Run("multiple", func(t *testing.T) {
		attrs := ds.Computed("total", "$price * $qty", "tax", "$total * 0.1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'total': () => $price * $qty, 'tax': () => $total * 0.1}", attrs["data-computed"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Computed("name") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestComputedKey(t *testing.T) {
	attrs := ds.ComputedKey("total", "$price * $qty")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$price * $qty", attrs["data-computed:total"])
}

func TestClass(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Class("hidden", "$isHidden")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'hidden': $isHidden}", attrs["data-class"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Class("hidden", "$isHidden", "font-bold", "$isBold")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'hidden': $isHidden, 'font-bold': $isBold}", attrs["data-class"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Class("hidden") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestClassKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.ClassKey("font-bold", "$isBold")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$isBold", attrs["data-class:font-bold"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.ClassKey("myClass", "$flag", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$flag", attrs["data-class:myClass__case.camel"])
	})
}

func TestAttr(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Attr("title", "$tooltip")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'title': $tooltip}", attrs["data-attr"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Attr("title", "$tooltip", "disabled", "$loading")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'title': $tooltip, 'disabled': $loading}", attrs["data-attr"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Attr("title") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestAttrKey(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		attrs := ds.AttrKey("disabled", "$loading")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$loading", attrs["data-attr:disabled"])
	})

	t.Run("title", func(t *testing.T) {
		attrs := ds.AttrKey("title", "'Theme: ' + $theme")
		require.Len(t, attrs, 1)
		assert.Equal(t, "'Theme: ' + $theme", attrs["data-attr:title"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.AttrKey("myAttr", "$val", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$val", attrs["data-attr:myAttr__case.camel"])
	})
}

func TestStyle(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Style("display", "$hiding && 'none'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'display': $hiding && 'none'}", attrs["data-style"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Style("display", "'none'", "color", "'red'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'display': 'none', 'color': 'red'}", attrs["data-style"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Style("display") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestStyleKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.StyleKey("background-color", "$red ? 'red' : 'blue'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$red ? 'red' : 'blue'", attrs["data-style:background-color"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.StyleKey("myProp", "$val", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$val", attrs["data-style:myProp__case.camel"])
	})
}

func TestRef(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Ref("myEl")
		require.Len(t, attrs, 1)
		assert.Equal(t, "myEl", attrs["data-ref"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Ref("myEl", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, "myEl", attrs["data-ref__case.kebab"])
	})
}

func TestIndicator(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Indicator("fetching")
		require.Len(t, attrs, 1)
		assert.Equal(t, "fetching", attrs["data-indicator"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Indicator("fetching", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "fetching", attrs["data-indicator__case.camel"])
	})
}

func TestIgnore(t *testing.T) {
	t.Run("no modifiers", func(t *testing.T) {
		attrs := ds.Ignore()
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-ignore"])
	})

	t.Run("with self", func(t *testing.T) {
		attrs := ds.Ignore(ds.ModSelf)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-ignore__self"])
	})
}

func TestIgnoreMorph(t *testing.T) {
	attrs := ds.IgnoreMorph()
	require.Len(t, attrs, 1)
	assert.Equal(t, true, attrs["data-ignore-morph"])
}

func TestJSONSignals(t *testing.T) {
	t.Run("no filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{})
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-json-signals"])
	})

	t.Run("with terse modifier", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{}, ds.ModTerse)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-json-signals__terse"])
	})

	t.Run("with include filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/user/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /user/}", attrs["data-json-signals"])
	})

	t.Run("with exclude filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Exclude: "/temp$/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{exclude: /temp$/}", attrs["data-json-signals"])
	})

	t.Run("with both filters", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/^app/", Exclude: "/password/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /^app/, exclude: /password/}", attrs["data-json-signals"])
	})

	t.Run("with filter and terse", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/counter/"}, ds.ModTerse)
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /counter/}", attrs["data-json-signals__terse"])
	})
}

func TestPreserveAttr(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		attrs := ds.PreserveAttr("open")
		require.Len(t, attrs, 1)
		assert.Equal(t, "open", attrs["data-preserve-attr"])
	})

	t.Run("multiple", func(t *testing.T) {
		attrs := ds.PreserveAttr("open", "class")
		require.Len(t, attrs, 1)
		assert.Equal(t, "open class", attrs["data-preserve-attr"])
	})
}

// ---------------------------------------------------------------------------
// Plugin watchers
// ---------------------------------------------------------------------------

func TestOnIntersect(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect"])
	})

	t.Run("once full", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true", ds.ModOnce, ds.ModFull)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect__once__full"])
	})

	t.Run("with threshold", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true", ds.ModThreshold, ds.Threshold(0.25))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect__threshold.25"])
	})

	t.Run("with half", func(t *testing.T) {
		attrs := ds.OnIntersect("$half = true", ds.ModHalf)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$half = true", attrs["data-on-intersect__half"])
	})

	t.Run("with exit", func(t *testing.T) {
		attrs := ds.OnIntersect("$exited = true", ds.ModExit)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$exited = true", attrs["data-on-intersect__exit"])
	})
}

func TestOnInterval(t *testing.T) {
	t.Run("default interval", func(t *testing.T) {
		attrs := ds.OnInterval("$count++")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval"])
	})

	t.Run("custom duration", func(t *testing.T) {
		attrs := ds.OnInterval("$count++", ds.ModDuration, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval__duration.500ms"])
	})

	t.Run("with leading", func(t *testing.T) {
		attrs := ds.OnInterval("$count++", ds.ModDuration, ds.Ms(500), ds.Leading)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval__duration.500ms.leading"])
	})
}

func TestOnSignalPatch(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.OnSignalPatch("console.log(patch)")
		require.Len(t, attrs, 1)
		assert.Equal(t, "console.log(patch)", attrs["data-on-signal-patch"])
	})

	t.Run("with debounce", func(t *testing.T) {
		attrs := ds.OnSignalPatch("doStuff()", ds.ModDebounce, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "doStuff()", attrs["data-on-signal-patch__debounce.500ms"])
	})
}

func TestOnSignalPatchFilter(t *testing.T) {
	t.Run("include only", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/^counter$/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /^counter$/}", attrs["data-on-signal-patch-filter"])
	})

	t.Run("both", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/user/", Exclude: "/password/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /user/, exclude: /password/}", attrs["data-on-signal-patch-filter"])
	})
}

// ---------------------------------------------------------------------------
// Real-world usage patterns (from this project's templates)
// ---------------------------------------------------------------------------

func TestRealWorldPatterns(t *testing.T) {
	t.Run("workcenters create button", func(t *testing.T) {
		attrs := ds.OnClick("$createSlideoutOpen = true")
		assert.Equal(t, "$createSlideoutOpen = true", attrs["data-on:click"])
	})

	t.Run("slideout overlay close", func(t *testing.T) {
		attrs := ds.OnClick("$createSlideoutOpen = false")
		assert.Equal(t, "$createSlideoutOpen = false", attrs["data-on:click"])
	})

	t.Run("search input with debounce", func(t *testing.T) {
		attrs := ds.OnInput("@post('/api/workcenters/search')", ds.ModDebounce, ds.Ms(300))
		assert.Equal(t, "@post('/api/workcenters/search')", attrs["data-on:input__debounce.300ms"])
	})

	t.Run("page size change", func(t *testing.T) {
		attrs := ds.OnChange("@patch('/api/workcenters/pagesize')")
		assert.Equal(t, "@patch('/api/workcenters/pagesize')", attrs["data-on:change"])
	})

	t.Run("sidebar toggle", func(t *testing.T) {
		attrs := ds.OnClick("$sidebarCollapsed = !$sidebarCollapsed")
		assert.Equal(t, "$sidebarCollapsed = !$sidebarCollapsed", attrs["data-on:click"])
	})

	t.Run("theme toggle with attr binding", func(t *testing.T) {
		click := ds.OnClick("$theme = window.themeToggle.cycle()")
		title := ds.AttrKey("title", "'Theme: ' + $theme")
		merged := ds.Merge(click, title)
		assert.Equal(t, "$theme = window.themeToggle.cycle()", merged["data-on:click"])
		assert.Equal(t, "'Theme: ' + $theme", merged["data-attr:title"])
	})

	t.Run("custom web component event", func(t *testing.T) {
		attrs := ds.OnEvent("table-select", "$table.selected = evt.detail.ids")
		assert.Equal(t, "$table.selected = evt.detail.ids", attrs["data-on:table-select"])
	})

	t.Run("init with SSE connection", func(t *testing.T) {
		attrs := ds.Init("@get('/api/workcenters/updates',{requestCancellation: 'disabled'})")
		assert.Equal(t,
			"@get('/api/workcenters/updates',{requestCancellation: 'disabled'})",
			attrs["data-init"],
		)
	})

	t.Run("toast delay init", func(t *testing.T) {
		attrs := ds.Init("el.style.opacity = '0'; setTimeout(() => el.remove(), 300)", ds.ModDelay, ds.Ms(3000))
		assert.Equal(t,
			"el.style.opacity = '0'; setTimeout(() => el.remove(), 300)",
			attrs["data-init__delay.3000ms"],
		)
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Signals
// ---------------------------------------------------------------------------

func TestSignalsEdgeCases(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{})
		assert.Equal(t, "{}", attrs["data-signals"])
	})

	t.Run("nested objects", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"foo": map[string]any{
				"bar": map[string]any{
					"baz": 1,
				},
			},
		})
		assert.Contains(t, attrs["data-signals"], `"foo"`)
		assert.Contains(t, attrs["data-signals"], `"bar"`)
		assert.Contains(t, attrs["data-signals"], `"baz"`)
	})

	t.Run("boolean values", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"isActive": true,
			"isHidden": false,
		})
		assert.Contains(t, attrs["data-signals"], `true`)
		assert.Contains(t, attrs["data-signals"], `false`)
	})

	t.Run("null values", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"value": nil,
		})
		assert.Contains(t, attrs["data-signals"], `null`)
	})

	t.Run("array values", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"items": []int{1, 2, 3},
		})
		assert.Contains(t, attrs["data-signals"], `[1,2,3]`)
	})

	t.Run("empty array", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"emptyList": []string{},
		})
		assert.Contains(t, attrs["data-signals"], `[]`)
	})

	t.Run("mixed array types", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"mixed": []any{"text", 42, true, nil},
		})
		signal := attrs["data-signals"]
		assert.Contains(t, signal, `"text"`)
		assert.Contains(t, signal, `42`)
		assert.Contains(t, signal, `true`)
		assert.Contains(t, signal, `null`)
	})

	t.Run("local signal with underscore", func(t *testing.T) {
		attrs := ds.SignalKey("_localVar", "private")
		assert.Equal(t, "private", attrs["data-signals:_localVar"])
	})

	t.Run("signal with single underscore in middle", func(t *testing.T) {
		attrs := ds.SignalKey("my_signal", "value")
		assert.Equal(t, "value", attrs["data-signals:my_signal"])
	})

	t.Run("very large number", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"bigNum": 9007199254740991, // Max safe integer in JS
		})
		assert.Contains(t, attrs["data-signals"], `9007199254740991`)
	})

	t.Run("float values", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"price": 19.99,
			"tax":   0.15,
		})
		signal := attrs["data-signals"]
		assert.Contains(t, signal, `19.99`)
		assert.Contains(t, signal, `0.15`)
	})

	t.Run("string with single quotes", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"message": "it's working",
		})
		// JSON encoding should escape the quote
		assert.Contains(t, attrs["data-signals"], `it's working`)
	})

	t.Run("string with double quotes", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"message": `he said "hello"`,
		})
		// JSON encoding should escape the quotes
		assert.Contains(t, attrs["data-signals"], `he said \"hello\"`)
	})

	t.Run("unicode in values", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"greeting": "你好",
			"emoji":    "👋",
		})
		signal := attrs["data-signals"]
		assert.Contains(t, signal, `你好`)
		assert.Contains(t, signal, `👋`)
	})

	t.Run("nested array of objects", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{
			"todos": []map[string]any{
				{"id": 1, "title": "Task 1", "done": false},
				{"id": 2, "title": "Task 2", "done": true},
			},
		})
		signal := attrs["data-signals"]
		assert.Contains(t, signal, `"todos"`)
		assert.Contains(t, signal, `"id"`)
		assert.Contains(t, signal, `"title"`)
		assert.Contains(t, signal, `"done"`)
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Filters
// ---------------------------------------------------------------------------

func TestFilterEdgeCases(t *testing.T) {
	t.Run("empty filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{})
		// Empty filter should render as boolean true per Datastar spec
		assert.Equal(t, true, attrs["data-json-signals"])
	})

	t.Run("include only", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/user/"})
		assert.Equal(t, "{include: /user/}", attrs["data-json-signals"])
	})

	t.Run("exclude only", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Exclude: "/password/"})
		assert.Equal(t, "{exclude: /password/}", attrs["data-json-signals"])
	})

	t.Run("both include and exclude", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{
			Include: "/^foo/",
			Exclude: "/bar/",
		})
		assert.Equal(t, "{include: /^foo/, exclude: /bar/}", attrs["data-json-signals"])
	})

	t.Run("regex with special characters", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/foo\\.bar/"})
		assert.Equal(t, "{include: /foo\\.bar/}", attrs["data-json-signals"])
	})

	t.Run("regex with anchors", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{
			Include: "/^user$/",
			Exclude: "/^_/",
		})
		assert.Equal(t, "{include: /^user$/, exclude: /^_/}", attrs["data-json-signals"])
	})

	t.Run("regex with character classes", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/[a-z]+/"})
		assert.Equal(t, "{include: /[a-z]+/}", attrs["data-json-signals"])
	})

	t.Run("complex regex patterns", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{
			Include: "/^(user|admin)\\./",
			Exclude: "/(password|secret)/",
		})
		filter := attrs["data-json-signals"].(string)
		assert.Contains(t, filter, "include: /^(user|admin)\\./")
		assert.Contains(t, filter, "exclude: /(password|secret)/")
	})

	t.Run("OnSignalPatchFilter with include", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/user/"})
		assert.Equal(t, "{include: /user/}", attrs["data-on-signal-patch-filter"])
	})

	t.Run("OnSignalPatchFilter with both", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/^foo/", Exclude: "/bar/"})
		assert.Equal(t, "{include: /^foo/, exclude: /bar/}", attrs["data-on-signal-patch-filter"])
	})

	t.Run("unicode in regex", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/用户/"})
		assert.Equal(t, "{include: /用户/}", attrs["data-json-signals"])
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Special Characters in Expressions
// ---------------------------------------------------------------------------

func TestSpecialCharactersInExpressions(t *testing.T) {
	t.Run("single quotes in expression", func(t *testing.T) {
		attrs := ds.OnClick("alert('hello')")
		assert.Equal(t, "alert('hello')", attrs["data-on:click"])
	})

	t.Run("double quotes in expression", func(t *testing.T) {
		attrs := ds.OnClick(`alert("hello")`)
		assert.Equal(t, `alert("hello")`, attrs["data-on:click"])
	})

	t.Run("template literal syntax", func(t *testing.T) {
		attrs := ds.Text("`Hello ${$name}`")
		assert.Equal(t, "`Hello ${$name}`", attrs["data-text"])
	})

	t.Run("multiple statements with semicolon", func(t *testing.T) {
		attrs := ds.OnClick("$foo = 1; @post('/endpoint')")
		assert.Equal(t, "$foo = 1; @post('/endpoint')", attrs["data-on:click"])
	})

	t.Run("signal reference", func(t *testing.T) {
		attrs := ds.Text("$count")
		assert.Equal(t, "$count", attrs["data-text"])
	})

	t.Run("nested signal reference", func(t *testing.T) {
		attrs := ds.Text("$user.profile.name")
		assert.Equal(t, "$user.profile.name", attrs["data-text"])
	})

	t.Run("method call on signal", func(t *testing.T) {
		attrs := ds.Text("$name.toUpperCase()")
		assert.Equal(t, "$name.toUpperCase()", attrs["data-text"])
	})

	t.Run("chained method calls", func(t *testing.T) {
		attrs := ds.Text("$text.trim().toLowerCase()")
		assert.Equal(t, "$text.trim().toLowerCase()", attrs["data-text"])
	})

	t.Run("unicode in signal name", func(t *testing.T) {
		attrs := ds.BindExpr("用户名")
		assert.Equal(t, "用户名", attrs["data-bind"])
	})

	t.Run("unicode in expression", func(t *testing.T) {
		attrs := ds.OnClick("$greeting = '你好'")
		assert.Equal(t, "$greeting = '你好'", attrs["data-on:click"])
	})

	t.Run("emoji in expression", func(t *testing.T) {
		attrs := ds.Text("'👋 ' + $name")
		assert.Equal(t, "'👋 ' + $name", attrs["data-text"])
	})

	t.Run("array literal", func(t *testing.T) {
		attrs := ds.OnClick("$items = [1, 2, 3]")
		assert.Equal(t, "$items = [1, 2, 3]", attrs["data-on:click"])
	})

	t.Run("object literal", func(t *testing.T) {
		attrs := ds.OnClick("$user = {name: 'John', age: 30}")
		assert.Equal(t, "$user = {name: 'John', age: 30}", attrs["data-on:click"])
	})

	t.Run("arrow function", func(t *testing.T) {
		attrs := ds.ComputedKey("double", "() => $count * 2")
		assert.Equal(t, "() => $count * 2", attrs["data-computed:double"])
	})

	t.Run("ternary operator", func(t *testing.T) {
		attrs := ds.Text("$count > 0 ? 'yes' : 'no'")
		assert.Equal(t, "$count > 0 ? 'yes' : 'no'", attrs["data-text"])
	})

	t.Run("logical operators", func(t *testing.T) {
		attrs := ds.Show("$isLoggedIn && !$isLoading")
		assert.Equal(t, "$isLoggedIn && !$isLoading", attrs["data-show"])
	})

	t.Run("comparison operators", func(t *testing.T) {
		attrs := ds.Show("$count >= 10 && $count <= 100")
		assert.Equal(t, "$count >= 10 && $count <= 100", attrs["data-show"])
	})

	t.Run("event object reference", func(t *testing.T) {
		attrs := ds.OnInput("$value = evt.target.value")
		assert.Equal(t, "$value = evt.target.value", attrs["data-on:input"])
	})

	t.Run("element reference", func(t *testing.T) {
		attrs := ds.OnClick("$height = el.offsetHeight")
		assert.Equal(t, "$height = el.offsetHeight", attrs["data-on:click"])
	})

	t.Run("complex expression", func(t *testing.T) {
		expr := "$total = $items.reduce((sum, item) => sum + item.price, 0)"
		attrs := ds.OnClick(expr)
		assert.Equal(t, expr, attrs["data-on:click"])
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Modifier Combinations
// ---------------------------------------------------------------------------

func TestModifierCombinations(t *testing.T) {
	t.Run("debounce with leading", func(t *testing.T) {
		attrs := ds.OnInput("search()", ds.ModDebounce, ds.Ms(500), ds.Leading)
		assert.Equal(t, "search()", attrs["data-on:input__debounce.500ms.leading"])
	})

	t.Run("debounce with notrailing", func(t *testing.T) {
		attrs := ds.OnInput("search()", ds.ModDebounce, ds.Ms(500), ds.NoTrailing)
		assert.Equal(t, "search()", attrs["data-on:input__debounce.500ms.notrailing"])
	})

	t.Run("throttle with trailing", func(t *testing.T) {
		attrs := ds.OnScroll("track()", ds.ModThrottle, ds.Seconds(1), ds.Trailing)
		assert.Equal(t, "track()", attrs["data-on:scroll__throttle.1s.trailing"])
	})

	t.Run("throttle with noleading", func(t *testing.T) {
		attrs := ds.OnScroll("track()", ds.ModThrottle, ds.Seconds(1), ds.NoLeading)
		assert.Equal(t, "track()", attrs["data-on:scroll__throttle.1s.noleading"])
	})

	t.Run("window with debounce and leading", func(t *testing.T) {
		attrs := ds.OnClick("handler()", ds.ModWindow, ds.ModDebounce, ds.Ms(500), ds.Leading)
		assert.Equal(t, "handler()", attrs["data-on:click__window__debounce.500ms.leading"])
	})

	t.Run("once passive capture", func(t *testing.T) {
		attrs := ds.OnClick("init()", ds.ModOnce, ds.ModPassive, ds.ModCapture)
		assert.Equal(t, "init()", attrs["data-on:click__once__passive__capture"])
	})

	t.Run("prevent and stop", func(t *testing.T) {
		attrs := ds.OnSubmit("handleSubmit()", ds.ModPrevent, ds.ModStop)
		assert.Equal(t, "handleSubmit()", attrs["data-on:submit__prevent__stop"])
	})

	t.Run("delay with viewtransition", func(t *testing.T) {
		attrs := ds.OnClick("toggle()", ds.ModDelay, ds.Ms(500), ds.ModViewTransition)
		assert.Equal(t, "toggle()", attrs["data-on:click__delay.500ms__viewtransition"])
	})

	t.Run("case camel modifier", func(t *testing.T) {
		attrs := ds.OnEvent("my-event", "handle()", ds.ModCase, ds.Camel)
		assert.Equal(t, "handle()", attrs["data-on:my-event__case.camel"])
	})

	t.Run("case kebab modifier", func(t *testing.T) {
		attrs := ds.OnEvent("myEvent", "handle()", ds.ModCase, ds.Kebab)
		assert.Equal(t, "handle()", attrs["data-on:myEvent__case.kebab"])
	})

	t.Run("case snake modifier", func(t *testing.T) {
		attrs := ds.OnEvent("myEvent", "handle()", ds.ModCase, ds.Snake)
		assert.Equal(t, "handle()", attrs["data-on:myEvent__case.snake"])
	})

	t.Run("case pascal modifier", func(t *testing.T) {
		attrs := ds.OnEvent("myEvent", "handle()", ds.ModCase, ds.Pascal)
		assert.Equal(t, "handle()", attrs["data-on:myEvent__case.pascal"])
	})

	t.Run("bind with case camel", func(t *testing.T) {
		attrs := ds.Bind("my-signal", ds.ModCase, ds.Camel)
		assert.Equal(t, true, attrs["data-bind:my-signal__case.camel"])
	})

	t.Run("signals with case kebab", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"mySignal": 1}, ds.ModCase, ds.Kebab)
		assert.Contains(t, attrs, "data-signals__case.kebab")
	})

	t.Run("class with case camel", func(t *testing.T) {
		attrs := ds.ClassKey("my-class", "$active", ds.ModCase, ds.Camel)
		assert.Equal(t, "$active", attrs["data-class:my-class__case.camel"])
	})

	t.Run("intersect with once and full", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true", ds.ModOnce, ds.ModFull)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect__once__full"])
	})

	t.Run("intersect with threshold", func(t *testing.T) {
		attrs := ds.OnIntersect("$partial = true", ds.ModThreshold, ds.Threshold(0.5))
		assert.Equal(t, "$partial = true", attrs["data-on-intersect__threshold.50"])
	})

	t.Run("intersect with exit and half", func(t *testing.T) {
		attrs := ds.OnIntersect("$gone = true", ds.ModExit, ds.ModHalf)
		assert.Equal(t, "$gone = true", attrs["data-on-intersect__exit__half"])
	})

	t.Run("interval with duration and leading", func(t *testing.T) {
		attrs := ds.OnInterval("tick()", ds.ModDuration, ds.Ms(500), ds.Leading)
		assert.Equal(t, "tick()", attrs["data-on-interval__duration.500ms.leading"])
	})

	t.Run("signal patch with debounce", func(t *testing.T) {
		attrs := ds.OnSignalPatch("refresh()", ds.ModDebounce, ds.Ms(300))
		assert.Equal(t, "refresh()", attrs["data-on-signal-patch__debounce.300ms"])
	})

	t.Run("init with delay and viewtransition", func(t *testing.T) {
		attrs := ds.Init("setup()", ds.ModDelay, ds.Ms(1000), ds.ModViewTransition)
		assert.Equal(t, "setup()", attrs["data-init__delay.1000ms__viewtransition"])
	})

	t.Run("multiple case-sensitive modifiers", func(t *testing.T) {
		attrs := ds.OnClick("action()", ds.ModWindow, ds.ModOnce, ds.ModPrevent, ds.ModStop)
		assert.Equal(t, "action()", attrs["data-on:click__window__once__prevent__stop"])
	})

	t.Run("json signals with terse", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/user/"}, ds.ModTerse)
		assert.Equal(t, "{include: /user/}", attrs["data-json-signals__terse"])
	})

	t.Run("ignore with self", func(t *testing.T) {
		attrs := ds.Ignore(ds.ModSelf)
		assert.Equal(t, true, attrs["data-ignore__self"])
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Merge Complex Scenarios
// ---------------------------------------------------------------------------

func TestMergeComplexScenarios(t *testing.T) {
	t.Run("merge 10+ attributes", func(t *testing.T) {
		merged := ds.Merge(
			ds.Show("$visible"),
			ds.Text("$message"),
			ds.OnClick("toggle()"),
			ds.OnInput("update()"),
			ds.BindExpr("value"),
			ds.ClassKey("active", "$isActive"),
			ds.StyleKey("color", "$textColor"),
			ds.AttrKey("title", "$tooltip"),
			ds.Ref("myElement"),
			ds.Indicator("loading"),
			ds.Init("setup()"),
			ds.Effect("$count++"),
		)
		assert.Len(t, merged, 12)
		assert.Equal(t, "$visible", merged["data-show"])
		assert.Equal(t, "$message", merged["data-text"])
		assert.Equal(t, "toggle()", merged["data-on:click"])
		assert.Equal(t, "update()", merged["data-on:input"])
		assert.Equal(t, "value", merged["data-bind"])
		assert.Equal(t, "$isActive", merged["data-class:active"])
		assert.Equal(t, "$textColor", merged["data-style:color"])
		assert.Equal(t, "$tooltip", merged["data-attr:title"])
		assert.Equal(t, "myElement", merged["data-ref"])
		assert.Equal(t, "loading", merged["data-indicator"])
		assert.Equal(t, "setup()", merged["data-init"])
		assert.Equal(t, "$count++", merged["data-effect"])
	})

	t.Run("merge with override - same attribute type", func(t *testing.T) {
		merged := ds.Merge(
			ds.Show("$first"),
			ds.Show("$second"),
			ds.Show("$third"),
		)
		assert.Len(t, merged, 1)
		assert.Equal(t, "$third", merged["data-show"], "last value should win")
	})

	t.Run("merge with override - same event different modifiers", func(t *testing.T) {
		merged := ds.Merge(
			ds.OnClick("first()", ds.ModDebounce, ds.Ms(100)),
			ds.OnClick("second()", ds.ModThrottle, ds.Ms(200)),
		)
		// Different modifiers create different attribute keys, so both exist
		assert.Len(t, merged, 2)
		assert.Equal(t, "first()", merged["data-on:click__debounce.100ms"])
		assert.Equal(t, "second()", merged["data-on:click__throttle.200ms"])
	})

	t.Run("merge multiple events", func(t *testing.T) {
		merged := ds.Merge(
			ds.OnClick("handleClick()"),
			ds.OnInput("handleInput()"),
			ds.OnChange("handleChange()"),
			ds.OnFocus("handleFocus()"),
			ds.OnBlur("handleBlur()"),
		)
		assert.Len(t, merged, 5)
		assert.Equal(t, "handleClick()", merged["data-on:click"])
		assert.Equal(t, "handleInput()", merged["data-on:input"])
		assert.Equal(t, "handleChange()", merged["data-on:change"])
		assert.Equal(t, "handleFocus()", merged["data-on:focus"])
		assert.Equal(t, "handleBlur()", merged["data-on:blur"])
	})

	t.Run("merge multiple classes", func(t *testing.T) {
		merged := ds.Merge(
			ds.ClassKey("active", "$isActive"),
			ds.ClassKey("disabled", "$isDisabled"),
			ds.ClassKey("hidden", "$isHidden"),
		)
		assert.Len(t, merged, 3)
		assert.Equal(t, "$isActive", merged["data-class:active"])
		assert.Equal(t, "$isDisabled", merged["data-class:disabled"])
		assert.Equal(t, "$isHidden", merged["data-class:hidden"])
	})

	t.Run("merge multiple styles", func(t *testing.T) {
		merged := ds.Merge(
			ds.StyleKey("color", "$textColor"),
			ds.StyleKey("background", "$bgColor"),
			ds.StyleKey("font-size", "$fontSize"),
		)
		assert.Len(t, merged, 3)
		assert.Equal(t, "$textColor", merged["data-style:color"])
		assert.Equal(t, "$bgColor", merged["data-style:background"])
		assert.Equal(t, "$fontSize", merged["data-style:font-size"])
	})

	t.Run("merge multiple attributes", func(t *testing.T) {
		merged := ds.Merge(
			ds.AttrKey("disabled", "$isDisabled"),
			ds.AttrKey("title", "$tooltip"),
			ds.AttrKey("aria-label", "$label"),
		)
		assert.Len(t, merged, 3)
		assert.Equal(t, "$isDisabled", merged["data-attr:disabled"])
		assert.Equal(t, "$tooltip", merged["data-attr:title"])
		assert.Equal(t, "$label", merged["data-attr:aria-label"])
	})

	t.Run("merge with empty attributes", func(t *testing.T) {
		merged := ds.Merge(
			ds.Show("$visible"),
			templ.Attributes{},
			ds.Text("$message"),
			templ.Attributes{},
		)
		assert.Len(t, merged, 2)
		assert.Equal(t, "$visible", merged["data-show"])
		assert.Equal(t, "$message", merged["data-text"])
	})

	t.Run("merge signals with other attributes", func(t *testing.T) {
		merged := ds.Merge(
			ds.Signals(map[string]any{"count": 0, "name": "test"}),
			ds.OnClick("$count++"),
			ds.Text("$name"),
		)
		assert.Len(t, merged, 3)
		assert.Contains(t, merged, "data-signals")
		assert.Equal(t, "$count++", merged["data-on:click"])
		assert.Equal(t, "$name", merged["data-text"])
	})

	t.Run("merge computed with signals", func(t *testing.T) {
		merged := ds.Merge(
			ds.Signals(map[string]any{"price": 10, "qty": 2}),
			ds.ComputedKey("total", "$price * $qty"),
			ds.Text("$total"),
		)
		assert.Len(t, merged, 3)
		assert.Contains(t, merged, "data-signals")
		assert.Equal(t, "$price * $qty", merged["data-computed:total"])
		assert.Equal(t, "$total", merged["data-text"])
	})

	t.Run("merge with conflicting class definitions", func(t *testing.T) {
		merged := ds.Merge(
			ds.ClassKey("active", "$foo"),
			ds.ClassKey("active", "$bar"), // Same class, different expression
		)
		assert.Len(t, merged, 1)
		assert.Equal(t, "$bar", merged["data-class:active"], "last value should win")
	})

	t.Run("real world form example", func(t *testing.T) {
		merged := ds.Merge(
			ds.BindExpr("email"),
			ds.OnInput("validateEmail()", ds.ModDebounce, ds.Ms(300)),
			ds.ClassKey("error", "$emailError"),
			ds.AttrKey("aria-invalid", "$emailError"),
			ds.Show("!$isLoading"),
		)
		assert.Len(t, merged, 5)
		assert.Equal(t, "email", merged["data-bind"])
		assert.Equal(t, "validateEmail()", merged["data-on:input__debounce.300ms"])
		assert.Equal(t, "$emailError", merged["data-class:error"])
		assert.Equal(t, "$emailError", merged["data-attr:aria-invalid"])
		assert.Equal(t, "!$isLoading", merged["data-show"])
	})

	t.Run("real world modal example", func(t *testing.T) {
		merged := ds.Merge(
			ds.Show("$isOpen"),
			ds.OnClick("$isOpen = false", ds.ModWindow),
			ds.ClassKey("active", "$isOpen"),
			ds.AttrKey("role", "'dialog'"),
			ds.AttrKey("aria-modal", "true"),
			ds.Init("$isOpen = false"),
		)
		assert.Len(t, merged, 6)
		assert.Equal(t, "$isOpen", merged["data-show"])
		assert.Equal(t, "$isOpen = false", merged["data-on:click__window"])
		assert.Equal(t, "$isOpen", merged["data-class:active"])
		assert.Equal(t, "'dialog'", merged["data-attr:role"])
		assert.Equal(t, "true", merged["data-attr:aria-modal"])
		assert.Equal(t, "$isOpen = false", merged["data-init"])
	})
}

// ---------------------------------------------------------------------------
// Edge Cases - Boundary Conditions
// ---------------------------------------------------------------------------

func TestBoundaryConditions(t *testing.T) {
	t.Run("threshold minimum 0.01", func(t *testing.T) {
		threshold := ds.Threshold(0.01)
		assert.Equal(t, ".01", string(threshold))
	})

	t.Run("threshold maximum 0.99", func(t *testing.T) {
		threshold := ds.Threshold(0.99)
		assert.Equal(t, ".99", string(threshold))
	})

	t.Run("threshold edge case 0.001", func(t *testing.T) {
		threshold := ds.Threshold(0.001)
		assert.Equal(t, ".00", string(threshold)) // Rounds to 2 decimal places
	})

	t.Run("threshold edge case 0.999", func(t *testing.T) {
		threshold := ds.Threshold(0.999)
		assert.Equal(t, "1.00", string(threshold)) // Rounds to 2 decimal places
	})

	t.Run("duration zero milliseconds", func(t *testing.T) {
		dur := ds.Ms(0)
		assert.Equal(t, ".0ms", string(dur))
	})

	t.Run("duration zero seconds", func(t *testing.T) {
		dur := ds.Seconds(0)
		assert.Equal(t, ".0s", string(dur))
	})

	t.Run("duration very large milliseconds", func(t *testing.T) {
		dur := ds.Ms(30000) // 30 seconds in ms
		assert.Equal(t, ".30000ms", string(dur))
	})

	t.Run("duration very large seconds", func(t *testing.T) {
		dur := ds.Seconds(3600) // 1 hour
		assert.Equal(t, ".3600s", string(dur))
	})

	t.Run("very long expression", func(t *testing.T) {
		// Simulate a complex expression
		expr := "$items.filter(item => item.active && item.price > 0).map(item => ({...item, total: item.price * item.qty})).reduce((sum, item) => sum + item.total, 0)"
		attrs := ds.OnClick(expr)
		assert.Equal(t, expr, attrs["data-on:click"])
	})

	t.Run("empty expression", func(t *testing.T) {
		attrs := ds.OnClick("")
		assert.Equal(t, "", attrs["data-on:click"])
	})

	t.Run("very long URL", func(t *testing.T) {
		url := "/api/very/long/path/with/many/segments/to/test/boundary/conditions/endpoint?param1=value1&param2=value2&param3=value3"
		result := ds.Get(url)
		assert.Contains(t, result, url)
	})

	t.Run("many format placeholders", func(t *testing.T) {
		result := ds.Get("/api/%s/%d/%s/%d/%s", "users", 1, "posts", 2, "comments")
		assert.Equal(t, "@get('/api/users/1/posts/2/comments')", result)
	})

	t.Run("zero duration from time.Duration", func(t *testing.T) {
		dur := ds.Duration(0)
		assert.Equal(t, ".0ms", string(dur))
	})

	t.Run("duration 500 microseconds rounds up", func(t *testing.T) {
		dur := ds.Duration(500 * time.Microsecond)
		assert.Equal(t, ".1ms", string(dur))
	})

	t.Run("duration 100 microseconds rounds down", func(t *testing.T) {
		dur := ds.Duration(100 * time.Microsecond)
		assert.Equal(t, ".0ms", string(dur))
	})

	t.Run("very small threshold", func(t *testing.T) {
		// Test that very small thresholds are handled
		threshold := ds.Threshold(0.005)
		// Should round to .01 or .00 depending on rounding
		assert.Contains(t, []string{".00", ".01"}, string(threshold))
	})

	t.Run("threshold at exactly 0.5", func(t *testing.T) {
		threshold := ds.Threshold(0.5)
		assert.Equal(t, ".50", string(threshold))
	})

	t.Run("long signal name", func(t *testing.T) {
		longName := "veryLongSignalNameThatMightBeUsedInSomeApplicationWithDescriptiveVariableNames"
		attrs := ds.SignalKey(longName, "value")
		assert.Equal(t, "value", attrs["data-signals:"+longName])
	})

	t.Run("long class name", func(t *testing.T) {
		longClass := "very-long-css-class-name-that-might-exist-in-utility-first-frameworks"
		attrs := ds.ClassKey(longClass, "$active")
		assert.Equal(t, "$active", attrs["data-class:"+longClass])
	})

	t.Run("many modifiers chained", func(t *testing.T) {
		attrs := ds.OnClick("action()",
			ds.ModWindow,
			ds.ModOnce,
			ds.ModPassive,
			ds.ModCapture,
			ds.ModPrevent,
			ds.ModStop,
			ds.ModDebounce,
			ds.Ms(500),
			ds.Leading,
		)
		expected := "data-on:click__window__once__passive__capture__prevent__stop__debounce.500ms.leading"
		assert.Equal(t, "action()", attrs[expected])
	})

	t.Run("empty signal map", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{})
		assert.Equal(t, "{}", attrs["data-signals"])
	})

	t.Run("single character expression", func(t *testing.T) {
		attrs := ds.Text("x")
		assert.Equal(t, "x", attrs["data-text"])
	})

	t.Run("expression with many semicolons", func(t *testing.T) {
		expr := "$a = 1; $b = 2; $c = 3; $d = 4; $e = 5"
		attrs := ds.OnClick(expr)
		assert.Equal(t, expr, attrs["data-on:click"])
	})

	t.Run("url with many query parameters", func(t *testing.T) {
		url := "/api/data?a=1&b=2&c=3&d=4&e=5&f=6&g=7&h=8&i=9&j=10"
		result := ds.Get(url)
		assert.Equal(t, "@get('"+url+"')", result)
	})
}
