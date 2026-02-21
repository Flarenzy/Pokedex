# Pokedex CLI

[![Go Tests](https://github.com/Flarenzy/Pokedex/actions/workflows/tests.yml/badge.svg)](https://github.com/Flarenzy/Pokedex/actions/workflows/tests.yml)
![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/Flarenzy/Pokedex/badges/coverage-badge.json)

A command-line Pokedex application written in Go.

This project lets you:
- browse location areas (`map`, `mapb`)
- explore encounter data for an area (`explore`)
- attempt to catch Pokemon (`catch`)
- inspect caught Pokemon (`inspect`)
- list your caught collection (`pokedex`)

## Origin

This repository was originally created by following the Boot.dev course project on building a Pokedex in Go, then expanded with additional refactors, testing, and CI.

## Testing

Run all tests locally:

```bash
go test ./...
```

Run coverage locally:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## CI

GitHub Actions runs the full test suite on:
- pushes to `main`
- pull requests targeting `main`

Coverage is also generated in CI on every run:
- step summary includes total and per-package/function coverage
- `coverage.out` and `coverage.txt` are uploaded as workflow artifacts
- coverage badge JSON is published to the `badges` branch (`coverage-badge.json`)

Workflow file: `.github/workflows/tests.yml`
