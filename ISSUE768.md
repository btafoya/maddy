### [\#768] Don't DKIM-oversign Sender header field
**URL**: https://github.com/foxcpp/maddy/issues/768
**State**: OPEN
**Labels**: security, mta-out, good first issue, 

**Body:**
Sender header field is oversigned even though it is not documented. To make maddy more compatible with the mailing list, it should be signed only if it is included and NOT oversigned.