### [\#731] How can I pull local_domains from table.sql_query
**URL**: https://github.com/foxcpp/maddy/issues/731
**State**: OPEN
**Labels**: new feature, 

**Body:**
# Use case

I am building a systme where I add domains to maddy via api. I have already built the api part (if anyone finds this useful I'm happy to share). Now I'd like to, inatead of passing local_domains via ENV, read them from a postgres table, in a way that I do not need to restart maddy to make these domains visible to the modules.

Note alternatives you considered and why they are not useful.

# Your idea for a solution
Tried this
```
table.sql_query allowed_domains {
    driver postgres
    named_args no
    dsn "postgres dsn"
    lookup "SELECT domain FROM domains where domain = $1"
    list "SELECT domain FROM domains"
}

whatever {
		modify {
                dkim {
                    domains $(primary_domain) $(local_domains) &allowed_domains
                    selector default
                    key_path dkim-keys/{domain}-{selector}.key
                    sig_expiry 120h # 5 days
                    hash sha256
                    newkey_algo rsa2048
                }
            }
}
```
How your solution would work in general?

Did not, domains are not in the list

- [x ] I'm willing to help with the implementation