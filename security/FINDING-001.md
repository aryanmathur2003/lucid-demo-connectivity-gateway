# FINDING-001: vulnerable dependency golang.org/x/text

| Field | Value |
| --- | --- |
| Severity | High |
| Advisory | CVE-2022-32149 (Go advisory GO-2022-1059) |
| Package | golang.org/x/text |
| Installed version | v0.3.7 |
| Fixed version | v0.3.8 or later |
| Affected symbol | `language.ParseAcceptLanguage` |
| Source | Dependency scan of `go.mod` |
| Due | 30 days from detection under the policy for high findings |

## Summary

`language.ParseAcceptLanguage` can take a long time to parse a crafted
Accept-Language header. An attacker can use this to cause a denial of
service.

## What the security team needs to close this

- Whether the affected symbol is reachable from a request in this service.
- The upgrade, with the build and tests passing.
- A re-scan showing the advisory no longer applies.

This finding was written by hand for a demo. Replace it with real scanner
output before presenting.
