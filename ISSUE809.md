### [\#809] Try to use IP address if no MX record?
**URL**: https://github.com/foxcpp/maddy/issues/809
**State**: OPEN
**Labels**: 

**Body:**
Currently, when trying to either use an IP address as the domain for an email address, or a domain without an MX record, it returns: `SMTP error 554: MX lookup error`
Could there be an option to try to send to the IP address directly?
(This seems to work in Outlook, for compatibility's sake?)

It would be useful for local email communications, maybe the mx-less domain part less, but it'd still be useful, it would make sense for this to be disabled by default.

---

### Resolution

This issue has been addressed by implementing the requested functionality in the `target.remote` module.

1.  **A/AAAA Record Fallback**: The logic in `internal/target/remote/connect.go` has been corrected to properly handle cases where no MX records are found. It now ignores "not found" errors from the DNS resolver and proceeds to attempt delivery to A/AAAA records for the domain, as specified in RFC 5321.

2.  **IP Address Literal Support**: A new configuration option, `allow_ip_literals`, has been added to the `target.remote` module. When set to `true`, maddy will accept recipient addresses with IP address literals (e.g., `user@[1.2.3.4]`). This allows for direct delivery to an IP address, bypassing MX lookups. This option is disabled by default to maintain standard behavior.
