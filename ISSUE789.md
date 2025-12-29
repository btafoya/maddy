### [\#789] TLS Connection Failure with Go 1.22+ Due to Disabled RSA Cipher Suites （Version 0.7.1 is good, but problems occur in later versions )
**URL**: https://github.com/foxcpp/maddy/issues/789
**State**: OPEN
**Labels**: bug, 

**Body:**
Problem Description
After upgrading maddy to compile with Go 1.23 (or Go 1.22+), sending emails to certain mail servers fails due to TLS handshake errors. This occurs because Go 1.22+ removed support for the insecure RSA key exchange cipher suites by default. Remote servers that only support these legacy cipher suites (e.g., TLS_RSA_WITH_*) reject the connection, preventing email delivery.

Root Cause
As documented in Go's release notes ([Go 1.22](https://go.dev/doc/go1.22#tls) and [Go 1.23](https://go.dev/doc/go1.23#tls)):

Go 1.22 removed RSA key exchange-based cipher suites from the default TLS configuration.

Go 1.23 further removed 3DES-based suites.
This causes compatibility issues with older mail servers still relying on RSA key exchange.

Proposed Solutions

Compile-Time Workaround
Re-enable RSA support via the GODEBUG flag during compilation:

bash
GODEBUG=tlsrsakex=1 go build ./cmd/maddy  
Configuration File Workaround
Explicitly specify legacy-compatible cipher suites in maddy.conf:

go
tls {  
  ciphers "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 TLS_RSA_WITH_AES_128_CBC_SHA"  
  # Add other required suites (prioritize ECDHE where possible)  
}  
Note: Include TLS_RSA_WITH_* suites only if strictly necessary, as they are less secure.

Recommendations

Documentation: Update build/configuration docs to warn users about Go 1.22+ compatibility issues and provide the above workarounds.

Fallback Logic: Consider adding runtime warnings when TLS handshakes fail due to cipher mismatches, suggesting recompilation with tlsrsakex=1 or config adjustments.

Reproduction Steps

Compile maddy with Go 1.23 (no custom GODEBUG).

Attempt delivery to a server enforcing RSA key exchange (e.g., older Exchange/SMTP services).

Observe TLS handshake failure in logs:

text
tls: no cipher suite supported by both client and server  
References

Go TLS defaults: [crypto/tls/common.go](https://github.com/golang/go/blob/master/src/crypto/tls/common.go#L1562-L1581)

Go 1.22 release notes: [tlsrsakex=1](https://go.dev/doc/go1.22#tls)

Thank you for maintaining maddy! Let me know if further details are needed.