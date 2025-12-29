### [\#756] Body encoding cannot be NIL
**URL**: https://github.com/foxcpp/maddy/issues/756
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

According to [this comment](https://github.com/pimalaya/himalaya/issues/525#issuecomment-2613920135), maddy seems to send NIL body encoding to FETCH response (within BODYSTRUCTURE) whereas it [cannot be null](https://github.com/pimalaya/himalaya/issues/525#issuecomment-2616931726).

# Steps to reproduce

NC

# Log files

```log
maddy  | 2025-01-25T10:30:52.572441139Z [debug] tls.loader.acme/handshake: choosing certificate {"identifier":"imap.services.domain.tld","num_choices":1}                   
maddy  | 2025-01-25T10:30:52.572507024Z [debug] tls.loader.acme/handshake: default certificate selection results        {"hash":"ad2fea10f731303445f1be751303e1e2dcfb39760b252
43baa899121e9b6aa02","identifier":"imap.services.domain.tld","issuer_key":"acme-v02.api.letsencrypt.org-directory","managed":true,"subjects":["imap.services.domain.tld"]}
maddy  | 2025-01-25T10:30:52.572521371Z [debug] tls.loader.acme/handshake: matched certificate in cache {"expiration":"2025-04-19T14:50:02.000","hash":"ad2fea10f731303445f1be
751303e1e2dcfb39760b25243baa899121e9b6aa02","managed":true,"remote_ip":"192.168.1.2","remote_port":"38888","subjects":["imap.services.domain.tld"]}                         
maddy  | 2025-01-25T10:30:52.575758785Z imap: * OK [CAPABILITY IMAP4rev1 LITERAL+ SASL-IR CHILDREN UNSELECT MOVE IDLE APPENDLIMIT AUTH=PLAIN COMPRESS] IMAP4rev1 Service Ready
maddy  | 2025-01-25T10:30:52.575972430Z imap: 0.0.wnJF090B AUTHENTICATE PLAIN AHByaW50ZXIAZ2MzJTNpRkZGSkR3UmRMNGg0OGI=                                                        
maddy  | 2025-01-25T10:30:52.803061876Z imap: 0.0.wnJF090B OK [CAPABILITY IMAP4rev1 LITERAL+ SASL-IR CHILDREN UNSELECT MOVE IDLE APPENDLIMIT I18NLEVEL=1 SORT THREAD=ORDEREDSU
BJECT COMPRESS NAMESPACE] AUTHENTICATE completed                                                                                                                              
maddy  | 2025-01-25T10:30:52.803515145Z imap: 0.1.lgySQXmV SELECT INBOX                                                                                                       
maddy  | 2025-01-25T10:30:52.803906957Z [debug] imapsql: initialized uidMap for selected mailbox: 2 [8 9]                                                                     
maddy  | 2025-01-25T10:30:52.804090726Z imap: * FLAGS (\Seen \Answered \Flagged \Deleted \Draft old)                                                                          
maddy  | 2025-01-25T10:30:52.804132715Z imap: * OK [PERMANENTFLAGS (\Seen \Answered \Flagged \Deleted \Draft \* old)] Flags permitted.                                        
maddy  | 2025-01-25T10:30:52.804191286Z imap: * OK [UNSEEN 2] Message 2 is first unseen                                                                                       
maddy  | 2025-01-25T10:30:52.804287068Z imap: * OK [UIDVALIDITY 1437467386] UIDs valid                                                                                        
maddy  | 2025-01-25T10:30:52.804301274Z imap: * 2 EXISTS                                                                                                                      
maddy  | 2025-01-25T10:30:52.804324398Z imap: * 0 RECENT                                                                                                                      
maddy  | 2025-01-25T10:30:52.804340288Z imap: * OK [UIDNEXT 10] Predicted next UID                                                                                            
maddy  | 2025-01-25T10:30:52.804395293Z imap: 0.1.lgySQXmV OK [READ-WRITE] SELECT completed                                                                                   
maddy  | 2025-01-25T10:30:52.804660335Z imap: 0.2.kxqbTv0z FETCH 2:1 (UID FLAGS ENVELOPE BODYSTRUCTURE)                                                                       
maddy  | 2025-01-25T10:30:52.804985040Z [debug] imapsql: resolved false 8:9 to 8:9                                                                                            
maddy  | 2025-01-25T10:30:52.805530003Z [debug] imapsql: scanMessages: scanned msgId=8 (seq 1) [UID FLAGS ENVELOPE BODYSTRUCTURE]                                             
maddy  | 2025-01-25T10:30:52.805637647Z [debug] imapsql: scanMessages: scanned msgId=9 (seq 2) [UID FLAGS ENVELOPE BODYSTRUCTURE]                                             
maddy  | 2025-01-25T10:30:52.805692511Z imap: * 1 FETCH (UID 8 FLAGS (old \Seen) ENVELOPE ("Fri, 24 Jan 2025 18:10:41 +0100" "this is a test" ((NIL NIL "printer" "services.ca
sona.rocks")) ((NIL NIL "printer" "services.domain.tld")) ((NIL NIL "printer" "services.domain.tld")) ((NIL NIL "printer" "services.domain.tld")) NIL NIL NIL "<ggzxrvlr
b43i4xfot7bws2qntl6j5lxwe4a6eqzdseplqlylgo@ct6wl36ufvlu>") BODYSTRUCTURE ("text" "plain" ("charset" "us-ascii") NIL NIL NIL 7 1 NIL ("inline" ()) NIL NIL))                   
maddy  | 2025-01-25T10:30:52.805801507Z imap: * 2 FETCH (UID 9 FLAGS (old) ENVELOPE ("Fri, 24 Jan 2025 18:15:31 +0100" "bla" ((NIL NIL "printer" "services.domain.tld")) ((N
IL NIL "printer" "services.domain.tld")) ((NIL NIL "printer" "services.domain.tld")) ((NIL NIL "printer" "services.domain.tld")) NIL NIL NIL "<du46gib7rpdxuv4d5mb3sa25c
anezopr6cxdqe2wdqpjeuqkhr@rzdxpuppmuws>") BODYSTRUCTURE ("text" "plain" ("charset" "us-ascii") NIL NIL NIL 10 1 NIL ("inline" ()) NIL NIL))                                   
maddy  | 2025-01-25T10:30:52.805922988Z imap: 0.2.kxqbTv0z OK FETCH completed
```

# Configuration file

NC

# Environment information

NC

---
@petrm feel free to add incomplete or missing information information