GOTFLATS ?=
SHELL = /bin/bash

.PHONY: gen validate

gen:
	find gen/suites -maxdepth 1 -mindepth 1 -type d -print0 | xargs -I '{}' -n1 -0 bash -c 'dir="$$(basename {})" && echo "=== $${dir} ===" && cd {} && go run . -o "../../../corpus/$${dir}"'

validate:
	cd ./validate && go run ./main.go
