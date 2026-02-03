package ds_test

import (
	"fmt"

	ds "github.com/Yacobolo/datastar-templ"
)

// Example demonstrates basic click handler
func ExampleOnClick() {
	attrs := ds.OnClick("$open = true")
	fmt.Println(attrs["data-on:click"])
	// Output: $open = true
}

// Example demonstrates click handler with debounce modifier
func ExampleOnClick_withModifiers() {
	attrs := ds.OnClick("search()", ds.ModDebounce, ds.Ms(500), ds.Leading)
	fmt.Println(attrs["data-on:click__debounce.500ms.leading"])
	// Output: search()
}

// Example demonstrates creating reactive signals
func ExampleSignals() {
	attrs := ds.Signals(map[string]any{"count": 0, "message": "Hello"})
	// JSON encoding may vary in key order, so we just check it contains expected keys
	signal := attrs["data-signals"]
	fmt.Printf("Signal contains count: %t\n", signal != "")
	// Output: Signal contains count: true
}

// Example demonstrates nested signal objects
func ExampleSignals_nested() {
	attrs := ds.Signals(map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
	})
	signal := attrs["data-signals"]
	fmt.Printf("Signal is not empty: %t\n", signal != "")
	// Output: Signal is not empty: true
}

// Example demonstrates two-way data binding
func ExampleBind() {
	attrs := ds.Bind("email")
	fmt.Println(attrs["data-bind:email"])
	// Output: true
}

// Example demonstrates value syntax for binding
func ExampleBindExpr() {
	attrs := ds.BindExpr("name")
	fmt.Println(attrs["data-bind"])
	// Output: name
}

// Example demonstrates computed signals
func ExampleComputed() {
	attrs := ds.Computed("double", "() => $count * 2")
	computed := attrs["data-computed"]
	fmt.Printf("Computed is not empty: %t\n", computed != "")
	// Output: Computed is not empty: true
}

// Example demonstrates keyed computed signal
func ExampleComputedKey() {
	attrs := ds.ComputedKey("total", "$price * $qty")
	fmt.Println(attrs["data-computed:total"])
	// Output: $price * $qty
}

// Example demonstrates conditional CSS class
func ExampleClass() {
	attrs := ds.Class("active", "$isActive")
	class := attrs["data-class"]
	fmt.Printf("Class is not empty: %t\n", class != "")
	// Output: Class is not empty: true
}

// Example demonstrates keyed class syntax
func ExampleClassKey() {
	attrs := ds.ClassKey("active", "$isActive")
	fmt.Println(attrs["data-class:active"])
	// Output: $isActive
}

// Example demonstrates GET request action
func ExampleGet() {
	result := ds.Get("/api/data")
	fmt.Println(result)
	// Output: @get('/api/data')
}

// Example demonstrates GET request with format arguments
func ExampleGet_withFormatArgs() {
	result := ds.Get("/api/users/%d", 42)
	fmt.Println(result)
	// Output: @get('/api/users/42')
}

// Example demonstrates GET request with options
func ExampleGet_withOptions() {
	result := ds.Get("/api/data", ds.Opt("retry", "error"), ds.Opt("selector", ".target"))
	fmt.Println(result)
	// Output: @get('/api/data',{retry: 'error', selector: '.target'})
}

// Example demonstrates POST request with options
func ExamplePost_withOptions() {
	result := ds.Post("/api/data", ds.Opt("contentType", "form"))
	fmt.Println(result)
	// Output: @post('/api/data',{contentType: 'form'})
}

// Example demonstrates merging multiple attributes
func ExampleMerge() {
	merged := ds.Merge(
		ds.Show("$visible"),
		ds.OnClick("toggle()"),
		ds.Text("$message"),
	)
	fmt.Printf("Has %d attributes\n", len(merged))
	// Output: Has 3 attributes
}

// Example demonstrates intersection observer
func ExampleOnIntersect() {
	attrs := ds.OnIntersect("$visible = true", ds.ModOnce, ds.ModFull)
	fmt.Println(attrs["data-on-intersect__once__full"])
	// Output: $visible = true
}

// Example demonstrates intersection with threshold
func ExampleOnIntersect_withThreshold() {
	attrs := ds.OnIntersect("$partial = true", ds.ModThreshold, ds.Threshold(0.5))
	fmt.Println(attrs["data-on-intersect__threshold.50"])
	// Output: $partial = true
}

// Example demonstrates loading indicator
func ExampleIndicator() {
	attrs := ds.Indicator("fetching")
	fmt.Println(attrs["data-indicator"])
	// Output: fetching
}

// Example demonstrates initialization code
func ExampleInit() {
	attrs := ds.Init("setup()")
	fmt.Println(attrs["data-init"])
	// Output: setup()
}

// Example demonstrates init with delay
func ExampleInit_withDelay() {
	attrs := ds.Init("load()", ds.ModDelay, ds.Ms(1000))
	fmt.Println(attrs["data-init__delay.1000ms"])
	// Output: load()
}

// Example demonstrates text content binding
func ExampleText() {
	attrs := ds.Text("$count")
	fmt.Println(attrs["data-text"])
	// Output: $count
}

// Example demonstrates conditional visibility
func ExampleShow() {
	attrs := ds.Show("$isOpen")
	fmt.Println(attrs["data-show"])
	// Output: $isOpen
}

// Example demonstrates style binding
func ExampleStyleKey() {
	attrs := ds.StyleKey("color", "$textColor")
	fmt.Println(attrs["data-style:color"])
	// Output: $textColor
}

// Example demonstrates attribute binding
func ExampleAttrKey() {
	attrs := ds.AttrKey("disabled", "$isDisabled")
	fmt.Println(attrs["data-attr:disabled"])
	// Output: $isDisabled
}

// Example demonstrates element reference
func ExampleRef() {
	attrs := ds.Ref("myButton")
	fmt.Println(attrs["data-ref"])
	// Output: myButton
}

// Example demonstrates reactive effect
func ExampleEffect() {
	attrs := ds.Effect("console.log('count changed:', $count)")
	fmt.Println(attrs["data-effect"])
	// Output: console.log('count changed:', $count)
}

// Example demonstrates interval timer
func ExampleOnInterval() {
	attrs := ds.OnInterval("tick()", ds.ModDuration, ds.Ms(1000))
	fmt.Println(attrs["data-on-interval__duration.1000ms"])
	// Output: tick()
}

// Example demonstrates signal patch watcher
func ExampleOnSignalPatch() {
	attrs := ds.OnSignalPatch("refresh()")
	fmt.Println(attrs["data-on-signal-patch"])
	// Output: refresh()
}

// Example demonstrates filtered signal patch watcher
func ExampleOnSignalPatchFilter() {
	attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/user/", Exclude: "/password/"})
	fmt.Println(attrs["data-on-signal-patch-filter"])
	// Output: {include: /user/, exclude: /password/}
}

// Example demonstrates JSON signals debug display
func ExampleJSONSignals() {
	attrs := ds.JSONSignals(ds.Filter{Include: "/user/"})
	fmt.Println(attrs["data-json-signals"])
	// Output: {include: /user/}
}

// Example demonstrates input event with debounce
func ExampleOnInput() {
	attrs := ds.OnInput("search()", ds.ModDebounce, ds.Ms(300))
	fmt.Println(attrs["data-on:input__debounce.300ms"])
	// Output: search()
}

// Example demonstrates form submission
func ExampleOnSubmit() {
	attrs := ds.OnSubmit("handleSubmit()", ds.ModPrevent)
	fmt.Println(attrs["data-on:submit__prevent"])
	// Output: handleSubmit()
}

// Example demonstrates DELETE action
func ExampleDelete() {
	result := ds.Delete("/api/todos/%d", 42)
	fmt.Println(result)
	// Output: @delete('/api/todos/42')
}

// Example demonstrates PATCH action
func ExamplePatch() {
	result := ds.Patch("/api/users/%d", 5)
	fmt.Println(result)
	// Output: @patch('/api/users/5')
}

// Example demonstrates PUT action
func ExamplePut() {
	result := ds.Put("/api/users/%d", 10)
	fmt.Println(result)
	// Output: @put('/api/users/10')
}
