### [\#772] Max_message_size not taking effect
**URL**: https://github.com/foxcpp/maddy/issues/772
**State**: OPEN
**Labels**: 

**Body:**
When I try to send larger attachments (30.33MB), it always fails, but attachments of 13.66MB can be sent successfully. My settings allow for much more than 30.33MB, so why am I still encountering failures when sending the larger file?

```
imap: listening on tls://0.0.0.0:993	

imap: listening on tcp://0.0.0.0:143	

imapsql: loaded a large mailbox with 10110 messages, beware of performance issues	

submission: rDNS error	{"reason":"DNS response contained records which contain invalid names","src_ip":"192.168.1.3:48138"}

submission: incoming message	{"msg_id":"552598d7","sender":"xxxx","src_host":"[127.0.0.1]","src_ip":"192.168.1.3:48138","username":"xxxx"}

submission: RCPT ok	{"msg_id":"552598d7","rcpt":"xxxx"}

submission: DATA error	{"msg_id":"552598d7","reason":"I/O error while writing buffer: buffer: failed to write file: SMTP error 552: Maximum message size exceeded"}

submission: aborted	{"msg_id":"552598d7"}

```

![Image](https://github.com/user-attachments/assets/41a82680-c9b3-45ce-8b94-724ffb7e6103)