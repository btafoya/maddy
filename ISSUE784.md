### [\#784] Bug: ghcr images are lacking arm64
**URL**: https://github.com/foxcpp/maddy/issues/784
**State**: OPEN
**Labels**: bug, 

**Body:**
# Describe the bug

The images provided in ghcr.io/foxcpp/maddy seem to lack arm64.

# Steps to reproduce

1. Try to pull Maddy on arm64 host from ghcr.io.
2. See it fail due to manifest not including anything other than amd64

# Log files

```
docker pull ghcr.io/foxcpp/maddy:0.8.1
0.8.1: Pulling from foxcpp/maddy
no matching manifest for linux/arm64/v8 in the manifest list entries
```