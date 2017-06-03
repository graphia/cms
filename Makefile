SRC := $(filter-out backend/%_test.go,$(wildcard backend/*.go))
ALL := $(wildcard backend/*.go)
DEVELOPMENT_CONFIG = "config/development.yml"
TEST_CONFIG = "../config/test.yml"
HUGO_CONFIG = "./config/hugo.development.yml"
PRIVATE_KEY_PATH = "./keys/app.rsa"
PUBLIC_KEY_PATH = "./keys/app.rsa.pub"

build:
	rm -rf dist
	mkdir dist
	cd frontend && NODE_ENV=production brunch build --production
	cp -R frontend/public/cms dist/cms
	go build -o graphia-cms ${ALL}
	hugo --config ${HUGO_CONFIG}

test:
	go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

keep-testing:
	ls backend/*.go | entr -r go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

run-backend:
	ls backend/*.go | entr -r go run ${SRC} -log-to-file=true -config ${DEVELOPMENT_CONFIG}

run-frontend:
	brunch watch --server frontend

cleanup:
	rm -rf tests/tmp/**/*

generate-keys:
	openssl genrsa -out ${PRIVATE_KEY_PATH} 1024
	openssl rsa -in ${PRIVATE_KEY_PATH} -pubout > ${PUBLIC_KEY_PATH}
