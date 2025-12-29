### [\#750] Configuration reload
**URL**: https://github.com/foxcpp/maddy/issues/750
**State**: OPEN
**Labels**: 

**Body:**
This is a major refactoring since maddy has a lot of state and resource handles to deal with.

## Design overview

Current register - initialize pipeline is replaced with register - configure - start.

```go 
type Module interface {
  Configure(ctx context.Context, inlineArgs string, cfg *config.Map) error
  Start(ctx context.Context) error
  Close(ctx context.Context) error
}
```

All module instances are first registered, then lazily configured. After that, all modules are started. If module is unused (and is therefore not configured) - startup fails (like now).

All instances are stored in a Registry object that replaces current global `instances` variable. 

If module needs to access some exclusive resources - like BoltDB handler or TCP listener, the corresponding object is created and owner by a global resource holder. 

When maddy is restarted, all resource holders are preserved, but Registry is created from scratch. All configuration loading is reexecuted, module instances are recreated. Module instances re-obtain existing resource handles from corresponding holders. Then modules are started. If any resource holders remain unused - they are closed.

Network listeners are `dup`'ed on return to modules. This is done because the only way to interrupt running Accept is to close to listener. 

## Plan

- [x] Register - Configure - Start module initialization
- [x] Non-global module registry
- [x] Resource holders
- [x] Network listener resource holders
- [x] Reinitialization with resource hand-over on SIGUSR1