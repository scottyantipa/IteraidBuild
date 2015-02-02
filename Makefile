export GOPATH := $(CURDIR)

EXECUTABLES = cleanDB git_objects server init_repo mark_repo_open mongo_utils new_iteraid_static_server

.PHONY: all
all: build

.PHONY: fast
fast:
	make -C ./ -j

build: $(EXECUTABLES)
	make -C ./ui

.PHONY: clean
clean:
	rm -rf repo
	rm -rf bin
	./bash/kill_static_servers.sh || true # do the || true because this may error, but we don't care

.PHONY: clean_db
clean_db:
	./bin/cleanDB

.PHONY: repo
repo:
	./bin/init_repo

.PHONY: $(EXECUTABLES)
$(EXECUTABLES):
	go install $@

.PHONY: serve
serve:
	./bin/server


# Full reset -- clean, init the repo, start server
.PHONY: reset
reset:
	make -C ./ clean
	make -C ./
	make -C ./ repo
	make -C ./ serve

.PHONY: install
install:
	make -C ./ui install
