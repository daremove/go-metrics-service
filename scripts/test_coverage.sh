#!/bin/sh

go test -v -coverpkg="../..." -coverprofile=coverage.out -covermode=count ../internal/...
go tool cover -func coverage.out | grep total | awk '{print $3}'
