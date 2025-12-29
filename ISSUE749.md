### [\#749] Dynamic configuration expression support
**URL**: https://github.com/foxcpp/maddy/issues/749
**State**: OPEN
**Labels**: 

**Body:**
1. Define evaluation context that is passed across the codebase via context.Context.
2. Replace some (most?) configuration values with generic dynamic like:
```go
type DynamicString struct {
	...
}

func (ds DynamicString) Evaluate(ctx context.Context) (string, error)
```
3. Add support for single-quote escaping in maddy configuration syntax.
```
modify {
  replace_rcpt expr 'address_username(key)'
}
```
4. Define configuration syntax (like above) for dynamic expressions.