SRC := $(filter-out backend/%_test.go,$(wildcard backend/*.go))
ALL := $(wildcard backend/*.go)
DEFAULT_CONFIG = "config/default.yml"
TEST_CONFIG = "../config/test.yml"

build:
	go build -o graphia-cms ${ALL}

test:
	go test -v ${ALL} -log-to-file=false -config=${TEST_CONFIG}

keep-testing:
	ls backend/*.go | entr -r go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

run-backend:
	ls backend/*.go | entr -r go run ${SRC} -log-to-file=true -config ${DEFAULT_CONFIG}

run-frontend:
	brunch watch --server

cleanup:
	rm -rf tests/tmp/**/*
