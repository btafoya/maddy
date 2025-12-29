### [\#808] Enabling Dovecot SASL in submission config hangs maddy on startup, systemd times out and kills it
**URL**: https://github.com/foxcpp/maddy/issues/808
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

When following the documentation for using Dovecot SASL with maddy by enabling it in the submission config, starting maddy with systemd appears to cause it to hang and eventually systemd kills it due to timing out.

The problem does not happen when disabling SASL.

# Steps to reproduce

Enable Dovecot SASL as documentation says in https://maddy.email/reference/auth/dovecot_sasl/

Try starting maddy after.

# Log files

systemd journal:
```
systemd[1]: Starting maddy.service - maddy mail server...
maddy[2950036]: smtp: listening on tcp://0.0.0.0:25
systemd[1]: maddy.service: start operation timed out. Terminating.
systemd[1]: maddy.service: Failed with result 'timeout'.
systemd[1]: Failed to start maddy.service - maddy mail server.
```

# Configuration file

`10-master.conf` Dovecot:
```
service auth {
  unix_listener auth-maddy-client {
    mode = 0600
    user = maddy
  }
}
```

maddy config:
```
$(hostname) = mail.vern.cc
$(primary_domain) = vern.cc
$(local_domains) = $(primary_domain)

tls file /etc/letsencrypt/live/$(primary_domain)/fullchain.pem /etc/letsencrypt/live/$(primary_domain)/privkey.pem

target.lmtp local_mailboxes {
    targets unix:///var/run/dovecot/lmtp-maddy
}

hostname $(hostname)

table.chain local_rewrites {
    optional_step regexp "([A-Za-z0-9_]+)[^A-Za-z0-9_](.+)@(.+)" "$1@$3"
    optional_step static {
        entry postmaster postmaster@$(primary_domain)
    }
    optional_step file /etc/maddy/aliases
    optional_step file /etc/maddy/aliases
    optional_step file /etc/maddy/aliases
    optional_step file /etc/maddy/aliases
}

msgpipeline local_routing {
    destination postmaster $(local_domains) {
        modify {
            replace_rcpt &local_rewrites
#            replace_rcpt file /etc/maddy/aliases
#            replace_rcpt file /etc/maddy/aliases
#            replace_rcpt file /etc/maddy/aliases
        }

        deliver_to &local_mailboxes
    }

    default_destination {
        reject 550 5.1.1 "User doesn\'t exist"
    }
}

smtp tcp://0.0.0.0:25 {
    limits {
        # Up to 20 msgs/sec across max. 10 SMTP connections.
        all rate 20 1s
        all concurrency 10
    }

    dmarc yes
    check {
        require_mx_record
        dkim
        spf
    }

    source $(local_domains) {
        reject 501 5.1.8 "Use Submission for outgoing SMTP"
    }
    default_source {
        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }
        default_destination {
            reject 550 5.1.1 "User doesn\'t exist"
        }
    }
}

submission tls://192.168.122.1:465 tcp://192.168.122.1:587 \
           #tls://[fe80::5054:ff:fe07:328f]:465 tcp://[fe80::5054:ff:fe07:328f]:587 
           tls://10.7.0.2:465 tcp://10.7.0.2:587 \
           tls://[2a01:2a01:4ff:f0::]:465 tcp://[2a01:2a01:4ff:f0::]:587 {
    limits {
        # Up to 50 msgs/sec across any amount of SMTP connections.
        all rate 50 1s
    }

    auth dovecot_sasl unix:///var/run/dovecot/auth-maddy-client
    #auth pass_table file /etc/dovecot/passwd
    #auth pass_table file /etc/maddy/passwd

    source $(local_domains) {
        check {
            authorize_sender {
                prepare_email &local_rewrites
                user_to_email identity
                from_normalize auto
                auth_normalize auto
            }
        }

        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }
        default_destination {
            modify {
                dkim $(primary_domain) $(local_domains) default
            }
            deliver_to &remote_queue
        }
    }
    default_source {
        reject 501 5.1.8 "Non-local sender domain"
    }
}

submission tls://127.0.0.1:465 tcp://127.0.0.1:587 {
    auth dummy
    source $(local_domains) {
        check {
            authorize_sender {
                prepare_email &local_rewrites
                user_to_email static {
                    entry "root" "*"
                }
            }
        }

        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }
        default_destination {
            modify {
                dkim $(primary_domain) $(local_domains) default
            }
            deliver_to &remote_queue
        }
    }
    default_source {
        reject 501 5.1.8 "Non-local sender domain"
    }
}

target.remote outbound_delivery {
    limits {
        # Up to 20 msgs/sec across max. 10 SMTP connections
        # for each recipient domain.
        destination rate 20 1s
        destination concurrency 10
    }
    mx_auth {
        dane
        mtasts {
            cache fs
            fs_dir mtasts_cache/
        }
        local_policy {
            min_tls_level encrypted
            min_mx_level none
        }
    }
}

target.queue remote_queue {
    target &outbound_delivery

    autogenerated_msg_domain $(primary_domain)
    bounce {
        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }
        default_destination {
            reject 550 5.0.0 "Refusing to send DSNs to non-local addresses"
        }
    }
}
```

# Environment information

* maddy version: `v0.8.2-0.20250309124430-fa47d40f6d51 linux/amd64 go1.24.4`

_Note that I am posting this on behalf of the vern.cc admin._

---

### Resolution

This issue has been resolved by modifying the `auth.dovecot_sasl` module to prevent blocking during maddy's startup sequence.

The root cause of the hang was that the module would immediately attempt to connect to the Dovecot SASL socket during its initialization phase. If the Dovecot service was not yet ready to accept connections, this would cause maddy to block indefinitely, leading to a startup timeout from systemd.

The fix involved the following changes to `internal/auth/dovecot_sasl/dovecot_sasl.go`:

1.  **Lazy Connections**: The blocking network dial has been removed from the `Init` function.
2.  **On-Demand Initialization**: The module now connects to the Dovecot socket and fetches the list of supported SASL mechanisms on the first actual authentication attempt, rather than at startup.

This ensures that maddy can start up successfully even if Dovecot is not immediately available, making the startup process more robust and avoiding the race condition.