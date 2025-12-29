### [\#783] proxy_protocol error
**URL**: https://github.com/foxcpp/maddy/issues/783
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

When adding proxy_protocol directive in smtp, submission and imap, I get the following error :
panic: runtime error: hash of unhashable type proxyprotocol.Listener

# Configuration file
```bash
imap tls://0.0.0.0:993 tcp://0.0.0.0:143 {
    proxy_protocol
    auth &lldap_auth
    storage &local_mailboxes
}
```

- maddy version: 0.8.1