# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

While configo is at `0.x`, the public API and the `configo` tag format may change
in any release. Breaking changes are called out under their own heading so they
are easy to find.

## v0.1.0 - 2026-10-02

### Added

- `LICENSE` (MIT), `CONTRIBUTING.md` and package documentation in `doc.go`.
- GitHub Actions CI running `gofmt`, `go mod tidy`, `go vet` and `go test -race`.
- A release workflow triggered by pushing a `v*` tag: it validates the tag
  against the module path, re-runs the suite, publishes the GitHub release from
  this file, and warms the module proxy for pkg.go.dev.
- Issue and pull request templates, and Dependabot updates for Go modules and
  GitHub Actions.

[Unreleased]: https://github.com/Marfien/configo/commits/main
