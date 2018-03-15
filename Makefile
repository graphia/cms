SRC := $(filter-out backend/%_test.go,$(wildcard backend/*.go))
ALL := $(wildcard backend/*.go)
DEVELOPMENT_CONFIG = "config/development.yml"
TEST_CONFIG = "../config/test.yml"
PRIVATE_KEY_PATH = "./keys/passwords/app.rsa"
PUBLIC_KEY_PATH = "./keys/passwords/app.rsa.pub"
SSL_CERT_PATH = "./keys/ssl/server.crt"
SSL_KEY_PATH = "./keys/ssl/server.key"

refresh-dist:
	rm -rf dist
	mkdir dist

build: refresh-dist
	cd frontend && NODE_ENV=production brunch build --production
	cp -R frontend/public/cms dist/cms
	go build -o graphia-cms ${ALL}

# Faster for working with cucumber
build-dev: refresh-dist
	cd frontend && brunch build
	cp -R frontend/public/cms dist/cms
	go build -o graphia-cms ${ALL}

build-backend-dev:
	go build -o graphia-cms ${ALL}

test:
	go test -v ${ALL} -log-to-file=true -config=${TEST_CONFIG}

cucumber: build-dev
	cd tests/frontend && bundle exec cucumber -p minimal

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
	rm -rf tests/tmp/**/* frontend/public/**/*

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
