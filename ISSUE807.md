### [\#807] hetzner dns api outdated
**URL**: https://github.com/foxcpp/maddy/issues/807
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

Hetzner DNS Console API is deprecated and it appears AMCE challenges are failing (at least I suspect that might be the cause for https://github.com/foxcpp/maddy/discussions/806)

libdns/hetzner has implemented support for the new Hetzner DNS Cloud API: https://github.com/libdns/hetzner/commit/43c0630a7ad716f9ebcccf7fc159749232254ef0

maddy needs to use `"github.com/libdns/hetzner/v2"` now instead. But it appears that causes a terrible dependency update chain. I have not had success building maddy after such an update.

# Steps to reproduce

See https://github.com/foxcpp/maddy/discussions/806

# Log files

Use a service like hastebin.com or attach a file if it is big

# Configuration file

See https://github.com/foxcpp/maddy/discussions/806

# Environment information

* maddy version: 0.8.1