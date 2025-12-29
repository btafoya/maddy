### [\#777] io_debug for target.remote
**URL**: https://github.com/foxcpp/maddy/issues/777
**State**: OPEN
**Labels**: new feature, 

**Body:**
# Use case

Hi, getting some repeated failures to send outgoing messages; in particular TLS failures are difficult to resolve after the fact and without detail on what happened. It'd be really helpful to be able to see the raw outgoing SMTP transport without MITMing myself; it seems I can do this for incoming IMAP and SMTP, but for outgoing connections that configuration directive does not seem to be supported.

See also #776