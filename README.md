# connectivity-gateway

Demo service. It returns the connectivity status of a vehicle to the mobile
app and the owner portal, with a status message in the caller's language.

Everything here is invented for a demo. The VINs are not real and this is not
Lucid code.

## Run

```
go run .
curl -H "Accept-Language: fr-CA" localhost:8080/v1/vehicles/5LDEMAA0000000001/connectivity
```

## Test

```
go test ./...
```

## Pipeline

`.github/workflows/ci.yml` runs `go vet`, the tests and `govulncheck`. The
scan job fails while the open finding in `security/` is unresolved.

## Conventions

- Never push to `main`. Open a pull request.
- A pull request description lists what changed, the test results and the
  re-scan result.
