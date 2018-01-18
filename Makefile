SRC := $(filter-out backend/%_test.go,$(wildcard backend/*.go))
ALL := $(wildcard backend/*.go)
DEVELOPMENT_CONFIG = "config/development.yml"
TEST_CONFIG = "../config/test.yml"
HUGO_CONFIG = "./config/hugo.development.yml"
PRIVATE_KEY_PATH = "./keys/passwords/app.rsa"
PUBLIC_KEY_PATH = "./keys/passwords/app.rsa.pub"
SSL_CERT_PATH = "./keys/ssl/server.crt"
SSL_KEY_PATH = "./keys/ssl/server.key"

build:
	rm -rf dist
	mkdir dist
	cd frontend && NODE_ENV=production brunch build --production
	cp -R frontend/public/cms dist/cms
	go build -o graphia-cms ${ALL}
	hugo --config ${HUGO_CONFIG}

# Faster for working with cucumber
build-dev:
	rm -rf dist
	mkdir dist
	cd frontend && brunch build
	cp -R frontend/public/cms dist/cms
	go build -o graphia-cms ${ALL}

build-backend-dev:
	go build -o graphia-cms ${ALL}

test:
	go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

cucumber: build-dev
	cd tests/frontend && cucumber

keep-testing:
	ls backend/*.go | entr -r go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

keep-building:
	#ls backend/*.go frontend/src/**/*.* | entr -r make build-dev
	find backend frontend/src -name "*.go" -or -name "*.js" -or -name "*.vue" | entr -r make build-dev

keep-building-backend:
	ls backend/*.go | entr -r make build-backend-dev

run-backend:
	ls backend/*.go | entr -r go run ${SRC} -log-to-file=true -config ${DEVELOPMENT_CONFIG}

run-frontend:
	#brunch watch --server frontend
	cd frontend && brunch watch --env development

cleanup:
	rm -rf tests/tmp/**/*

generate-password-keys:
	openssl genrsa -out ${PRIVATE_KEY_PATH} 1024
	openssl rsa -in ${PRIVATE_KEY_PATH} -pubout > ${PUBLIC_KEY_PATH}

generate-ssl-keys:
	openssl req \
		-new \
		-newkey rsa:4096 \
		-days 365 \
		-nodes \
		-x509 \
		-subj "/C=GB/ST=England/L=Manchester/O=Graphia Ltd./OU=dev/CN=localhost" \
		-keyout ${SSL_KEY_PATH} \
		-out ${SSL_CERT_PATH}
