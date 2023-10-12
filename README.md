# example-go-slug-mtime

Used to test [this change](https://github.com/hashicorp/go-slug/pull/47) to go-slug.

## Instructions

1. Clone this repo in a peer directory to `go-slug`
2. `go run main.go` will slug and unpack the `archive-dir` and the hash comparison should fail
3. Check out the branch associated with link above in `go-tfe`
4. `go mod edit -replace github.com/hashicorp/go-slug=../go-slug` will reference the fixed package
5. `go run main.go` should now succeed
6. You can run this tool against other directories by specifying an argument `go run main.go "./archive-dir-absolute"`
