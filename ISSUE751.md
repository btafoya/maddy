### [\#751] maddy 0.8 something broke email_localpart_optional
**URL**: https://github.com/foxcpp/maddy/issues/751
**State**: OPEN
**Labels**: bug, release-blocker, 

**Body:**
# Describe the bug

I have maddy set up as described here: https://maddy.email/multiple-domains/.

```
    check {                   
        authorize_sender {
            user_to_email chain {  
            step email_localpart_optional           # remove domain from username if present
            step email_with_domain $(local_domains) # expand username with all allowed domains
            }                                        
        }                                                         
    } 
```
Something has changed in 0.8, because when I authenticate with `user@email.tld`, I get the error below.

```
maddy  | 2025-01-25T19:31:05.563762356Z submission/sasl: authentication failed  {"reason":"no auth. provider accepted creds, last err: unknown credentials","src_ip":"192.168.1.2:59842","username":"printer@domain.tld"}
```

Authenticating with just `user` allows the client to log in, but then I end up with:
```
maddy  | 2025-01-25T19:38:41.643632364Z submission/sasl: authentication failed  {"reason":"no auth. provider accepted creds, last err: unknown credentials","src_ip":"192.168.1.2:39680","username":"printer@domain.tld"}
```

For now I rolled back to 0.7.1.

# Configuration file

Nothing was changed in the config.

* maddy version: 0.8