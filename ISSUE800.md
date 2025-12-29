### [\#800] SMTP Error code received doesn't match LMTP response
**URL**: https://github.com/foxcpp/maddy/issues/800
**State**: OPEN
**Labels**: bug, 

**Body:**
# Description of the bug

When using LMTP to offload IMAP delivery to dovecot for instance, the return SMTP error code is not properly propagated

# Steps to reproduce

1. Send a mail to a maddy server that is using LTMP with dovecot as an IMAP backend
2. Let's say the recipient is exceeding its quota with this email
3. Maddy receives back a 552 error code from Dovecot
4. The error code is rewritten to 452 as per: https://github.com/foxcpp/maddy/blob/master/internal/smtpconn/smtpconn.go#L121 (for legacy support of too many RCPT If I got it correctly ?)
5. The sender mail server receives back a 554 (Internal server error) from maddy

# Log files

Here is a log of the different interaction (sender/receiver/dovecot). Both sender and receiver are running a maddy server:

https://hastebin.milkywan.fr/aqupejalus.pl

# Configuration file

Simple LMTP taget:
```
target.lmtp local_mailboxes {
  targets unix:///var/run/dovecot2/lmtp-maddy
}
```
Called by a simple message pipeline:
```
msgpipeline local_routing {
  destination $(local_domains) {
    modify {
      replace_rcpt &local_rewrites
    }
    deliver_to &local_mailboxes
  }
  default_destination {
    reject 550 5.1.1 "User doesn't exist"
  }
}
```
For dovecot:
```
service lmtp {
  unix_listener lmtp-maddy {
    mode = 0600
    user = maddy
  }
}
```

# Environment information

* maddy version: 0.8.1 linux/amd64 go1.24.6