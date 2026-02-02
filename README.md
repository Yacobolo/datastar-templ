# datastar-templ

Type-safe [Datastar](https://data-star.dev) attribute helpers for [templ](https://github.com/a-h/templ) templates.

## Overview

`datastar-templ` is a Go library that provides compile-time type safety for Datastar attributes in templ templates. It bridges the gap between Go's templ templating system and Datastar's hypermedia framework, enabling you to build reactive web applications with full IDE autocomplete and type checking.

## Features

- **Type-Safe Attributes**: Compile-time checks for all Datastar attributes
- **Comprehensive Event Coverage**: 60+ typed DOM event handlers
- **SSE Action Builders**: Type-safe HTTP action expressions
- **Modifier System**: Full support for Datastar modifiers (debounce, throttle, etc.)
- **Signal Management**: Handle reactive state with signals and computed values
- **Merge Capability**: Combine multiple Datastar attributes on a single element
- **Flexible API**: Multiple syntax options for different use cases

## Installation

```bash
go get github.com/yourusername/datastar-templ
```

## Usage

Import the package (commonly aliased as `ds`):

```go
import ds "github.com/yourusername/datastar-templ"
```

### Basic Attributes

Use Datastar attributes in your templ templates:

```go
templ Button() {
    <button { ds.OnClick("$open = true")... }>
        Click me
    </button>
}
```

### Signals

Define reactive state:

```go
templ Component() {
    <div { ds.Signals(map[string]any{
        "count": 0,
        "message": "Hello",
    })... }>
        <span { ds.Text("$count")... }></span>
    </div>
}
```

### Event Handlers

Handle DOM events with modifiers:

```go
templ SearchInput() {
    <input 
        type="text"
        { ds.OnInput("search()", ds.ModDebounce, ds.Ms(300))... }
    />
}
```

### SSE Actions

Make server requests:

```go
templ TodoItem(id int) {
    <div>
        <button { ds.OnClick(ds.Delete("/api/todos/%d", id))... }>
            Delete
        </button>
    </div>
}
```

### Conditional Rendering

Show/hide elements reactively:

```go
templ Modal() {
    <div { ds.Show("$isOpen")... }>
        <button { ds.OnClick("$isOpen = false")... }>
            Close
        </button>
    </div>
}
```

### Data Binding

Bind form inputs to signals:

```go
templ Form() {
    <input 
        type="text"
        { ds.Bind("name")... }
    />
}
```

### Merging Attributes

Combine multiple Datastar attributes:

```go
templ ComplexElement() {
    <div { ds.Merge(
        ds.Show("$visible"),
        ds.OnClick("toggle()"),
        ds.Class("active", "$isActive"),
    )... }>
        Content
    </div>
}
```

## API Reference

### Core Functions

- `Merge(attrs ...templ.Attributes) templ.Attributes` - Combine multiple Datastar attributes

### Signal Management

- `Signals(m map[string]any) templ.Attributes` - Define reactive signals
- `SignalKey(key string, value any) templ.Attributes` - Define a single signal
- `Computed(key, expr string) templ.Attributes` - Define computed values
- `ComputedKey(key, expr string) templ.Attributes` - Define a single computed value

### Data Binding

- `Bind(key string, mods ...Modifier) templ.Attributes` - Bind to a signal
- `BindExpr(expr string, mods ...Modifier) templ.Attributes` - Bind with custom expression

### DOM Manipulation

- `Text(expr string) templ.Attributes` - Set element text content
- `Show(expr string) templ.Attributes` - Toggle element visibility
- `Class(key, expr string) templ.Attributes` - Toggle CSS classes
- `Style(key, expr string) templ.Attributes` - Set inline styles
- `Attr(key, expr string) templ.Attributes` - Set element attributes

### Event Handlers

Mouse events: `OnClick`, `OnDblClick`, `OnMouseOver`, `OnMouseEnter`, etc.

Keyboard events: `OnKeyDown`, `OnKeyUp`, `OnKeyPress`

Form events: `OnInput`, `OnChange`, `OnSubmit`, `OnFocus`, `OnBlur`

And 50+ more event handlers...

### SSE Actions

- `Get(url string, args ...any) string` - HTTP GET request
- `Post(url string, args ...any) string` - HTTP POST request
- `Put(url string, args ...any) string` - HTTP PUT request
- `Patch(url string, args ...any) string` - HTTP PATCH request
- `Delete(url string, args ...any) string` - HTTP DELETE request

With options:
```go
ds.Get("/api/data", ds.Opt("indicator", ".spinner"))
```

### Modifiers

Duration helpers:
- `Duration(amount int, unit string) Modifier`
- `Ms(amount int) Modifier`
- `Seconds(amount int) Modifier`
- `Threshold(percent int) Modifier`

Event modifiers:
- `ModDebounce`, `ModThrottle`, `ModOnce`, `ModPassive`, `ModCapture`, etc.

Casing modifiers:
- `Camel`, `Kebab`, `Snake`, `Pascal`

Timing modifiers:
- `Leading`, `NoLeading`, `Trailing`, `NoTrailing`

### Watchers

- `OnIntersect(expr string, mods ...Modifier) templ.Attributes` - Intersection observer
- `OnInterval(expr string, mods ...Modifier) templ.Attributes` - Timed intervals
- `OnSignalPatch(expr string) templ.Attributes` - Watch signal changes

### Utilities

- `Ref(key string) templ.Attributes` - Reference an element
- `Indicator(selector string, mods ...Modifier) templ.Attributes` - Loading indicator
- `Init(expr string) templ.Attributes` - Run code on initialization
- `Effect(expr string) templ.Attributes` - Run reactive effects
- `Ignore() templ.Attributes` - Ignore Datastar processing

## Examples

### Todo List

```go
templ TodoList(todos []Todo) {
    <div { ds.Signals(map[string]any{"todos": todos})... }>
        <form { ds.OnSubmit(ds.Post("/todos"), ds.ModPrevent)... }>
            <input { ds.Bind("newTodo")... } />
            <button type="submit">Add</button>
        </form>
        
        <ul>
            for _, todo := range todos {
                <li>
                    <span { ds.Text(fmt.Sprintf("$todos[%d].title", todo.ID))... }></span>
                    <button { ds.OnClick(ds.Delete("/todos/%d", todo.ID))... }>
                        Delete
                    </button>
                </li>
            }
        </ul>
    </div>
}
```

### Search with Debounce

```go
templ Search() {
    <div { ds.Signals(map[string]any{"query": "", "results": []string{}})... }>
        <input 
            type="search"
            { ds.Bind("query")... }
            { ds.OnInput(ds.Get("/search?q=$query"), ds.ModDebounce, ds.Ms(300))... }
        />
        
        <ul>
            <li { ds.Text("$results.length + ' results'")... }></li>
        </ul>
    </div>
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
datastar-templ/
├── ds.go          # Core types and helpers
├── consts.go      # Datastar constants
├── attrs.go       # Attribute builders
├── events.go      # Event handlers
├── actions.go     # SSE actions
└── *_test.go      # Tests
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Resources

- [Datastar Documentation](https://data-star.dev)
- [templ Documentation](https://templ.guide)
- [Go Documentation](https://golang.org/doc)
