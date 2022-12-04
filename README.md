# kaelly-portals 

![CI status](https://github.com/kaellybot/kaelly-portals/actions/workflows/build.yml/badge.svg?branch=main)

Application to retrieve dimension portals from different sources, written in Go

## Current supported sources

- [dofus-portals](https://dofus-portals.fr)

## Generate client boilerplate

```Bash
# CLI Installation
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Examples
oapi-codegen -package payloads -generate types,client,spec openapi.yaml > openapi.gen.go
oapi-codegen -package dofusportals -generate types,client,spec payloads/dofusportals/openapi.yaml > payloads/dofusportals/openapi.gen.go
```