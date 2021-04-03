SHELL := /usr/bin/env bash

.PHONY: app_build
app_build:
	yarn build

.PHONY: server_build
server_build:
	mage buildAll

# https://devcenter.heroku.com/articles/heroku-cli
.PHONY: heroku_setup
heroku_setup:
	heroku login
	./bin/heroku_setup.sh $(APP) $(ORG)

.PHONY: heroku_push
# heroku_push: app_build
# heroku_push: server_build
heroku_push:
	heroku container:login
	docker build -t registry.heroku.com/$(APP)/web .
	docker push registry.heroku.com/$(APP)/web
	heroku container:release web -a $(APP)
