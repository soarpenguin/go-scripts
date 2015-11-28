#!/usr/bin/env bash

go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'

