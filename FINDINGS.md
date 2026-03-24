# Security Vulnerabilities

**date:** 2026-03-20
**tool:** Trivy v0.69.3
**image:** rest-api-app:latest

## Summary
Image contains one HIGH severity vulnerability in zlib library (CVE-2026-22184).

## CVE-2026-22184 — Buffer Overflow in zlib (untgz utility)

### Description
zlib versions prior to `1.3.2` contain a global buffer overflow in the `untgz`
utility located in `contrib/untgz`. The vulnerability allows an attacker to execute
arbitrary code when calling `untgz` with an excessively long archive name
in the command line.

### Impact
Low - the vulnerability affects `untgz` utility which is not used by this application.
However, it must be fixed to comply with security best practices.

### Fix:
Update packages in the container by adding  following lines to Dockerfile:
```dockerfile
RUN apk update && apk upgrade
```
