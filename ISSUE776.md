### [\#776] io_debug: Extra option for only logging on failure
**URL**: https://github.com/foxcpp/maddy/issues/776
**State**: OPEN
**Labels**: new feature, 

**Body:**
# Use case

io_debug is currently [a boolean flag](https://maddy.email/reference/endpoints/smtp/#io_debug-boolean).

Turning it on to log everything when debugging either Maddy's sources or on bringing up a new server makes sense. However, it would also be useful for longer-term debugging of more complex issues with server setup if it could log the entire communication only when an error has occurred. Thus successful connections do not need the extra info, whilst problem communication can capture this for days or weeks without needing to generate overly verbose logs.

There is also the possibility that extended logging on failure/rejections/etc. could be valuable for feeding into systems like fail2ban or CrowdSec for blacklisting spam and attackers, enabling such systems to differentiate better between malicious connections and innnocent failures.


# Your idea for a solution

How your solution would work in general?
I haven;t yet looked to see how difficult this would be in terms of implementation. But for the interface, simply adding a third option to the config directive - "on", "off", and e.g. "on_failure" seems suitable to me.

- [x] I'm willing to help with the implementation
  * I think I can implement but would appreciate guidance and a review of changes to ensure it's implemented in the required places.