### [\#786] database is locked error when using reroute
**URL**: https://github.com/foxcpp/maddy/issues/786
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

I am trying to implement a simple mailing-list-like functionality using maddy only. I have added a corresponding configuration block so that recipients of emails addressed to a specified address get rewritten according to an aliases file and delivered either locally or injected to the remote queue. However when I send an email to both the mailing list address and some other address on the server, the IMAP database seems to get deadlocked with the email failing to deliver and breaking the delivery of any further emails until the restart of maddy.

# Steps to reproduce

Modify the default configuration file as below and send an email addressed to both mailing.list@example.org and postmaster@example.org.

# Log files

```
imapsql: readUids (listMsgUidsRecent) []: database is locked   {"mbox":"INBOX","mboxId":1,"username":"postmaster@example.org","uid":1}
imapsql: GetMailbox [INBOX]: imapsql: serialization failure, try again later   {"username":"postmaster@example.org","uid":1}
smtp: DATA error        {"msg_id":"dc377837","reason":" GetMailbox INBOX: imapsql: serialization failure, try again later"}
```

# Configuration file

I use the default configuration file with the following block added to the top of the default `local_routing` msgpipeline:
```
    destination mailing.list@example.org {
		modify {
			replace_rcpt file /etc/maddy/mailing-list
		}

		reroute {
			destination mailing.list@example.org {
				reject 550 5.1.1 "User doesn't exist"
			}

			destination postmaster $(local_domains) {
				deliver_to &local_routing
			}

			default_destination {
				deliver_to &remote_queue
			}
		}
	}
```

I also tried swapping the `deliver_to &local_routing` line above for `deliver_to &local_mailboxes` with the same outcome.

Note that the default `destination postmaster $(local_domains)` section is still present below the above to facilitate delivery to addresses different than the virtual mailing list.

# Environment information

* maddy version: 0.8.1 on glibc-based Void Linux