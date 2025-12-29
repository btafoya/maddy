### [\#748] Delivery target wrapper to ignore errors
**URL**: https://github.com/foxcpp/maddy/issues/748
**State**: OPEN
**Labels**: mta-in, good first issue, 

**Body:**
Might be useful if some `deliver_to` calls are intended to produce non-critical copies (e.g. push notifications, like #735).

```
deliver_to ignore_error http ...
```