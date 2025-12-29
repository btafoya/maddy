### [\#788] Bug report: IMAP APPENDLIMIT behave chaos
**URL**: https://github.com/foxcpp/maddy/issues/788
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

IMAP APPENDLIMIT settings behave in effected, conflict with maddy settings.

# Steps to reproduce

1. Run maddy in Docker as instructed in https://maddy.email/docker/
2. Create a `creds` and `imap-acct` for an account
3. Now we can expect this account has no LIMIT (NIL) with APPENDLIMIT because account and mboxes are both has no limit.

```
# docker exec -it maddy maddy imap-acct appendlimit foxcpp@maddy.test
No limit

# sqlite3 imapsql.db 'select * from mboxes where msgsizelimit is null'
1|1|INBOX|1|0||7|1439825317||5
2|1|Sent|1|0||1|3335251900|\Sent|0
3|1|Trash|1|0||1|3434545230|\Trash|0
4|1|Junk|1|0||1|418662061|\Junk|0
5|1|Drafts|1|0||1|1872697851|\Drafts|0
6|1|Archive|1|0||2|556772642|\Archive|0
```

4. But in effect, this account has an APPENDLIMIT to ZERO

```
# python
>>> import imaplib
>>> s=imaplib.IMAP4_SSL('127.0.0.1')
>>> s.login('foxcpp@maddy.test', 'foxcpp')
('OK', [b'[CAPABILITY IMAP4rev1 LITERAL+ SASL-IR CHILDREN UNSELECT MOVE IDLE APPENDLIMIT I18NLEVEL=1 SORT THREAD=ORDEREDSUBJECT COMPRESS NAMESPACE] LOGIN completed'])
>>> s.status('INBOX', '(APPENDLIMIT)')
('OK', [b'INBOX (APPENDLIMIT 0)'])
```

According to [RFC7889 The IMAP APPENDLIMIT Extension](https://www.rfc-editor.org/rfc/rfc7889.html), the `APPENDLIMIT` flag in `s.login()` indicates server support this extension and client will need to discover upload limits for each mailbox. And the STATUS Response indicates server will not accept any APPEND commands at all for the affected mailboxes.

I am expecting to see the following STATUS response, which indicates the mailbox has no limit.

```
>>> s.status('INBOX', '(APPENDLIMIT)')
('OK', [b'INBOX (APPENDLIMIT NIL)'])
```

# Log files

No need

# Configuration file

Exactly same as https://github.com/foxcpp/maddy/blob/fa47d40f6d510a431d2bbc238c7d36a58774ae2f/maddy.conf.docker

# Environment information

* maddy version: 0.8.1