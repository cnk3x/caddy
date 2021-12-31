.PHONY: test ctest covdir coverage docs linter qtest clean dep release logo license envvar
APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
BINARY:="identity"
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all: envvar build
	@echo "$@: complete"

envvar:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"

build:
	@mkdir -p bin/
	@rm -rf ./bin/*
	@CGO_ENABLED=0 go build -o ./bin/authdbctl $(VERBOSE) \
		-ldflags="-w -s \
		-X main.appVersion=$(APP_VERSION) \
		-X main.gitBranch=$(GIT_BRANCH) \
		-X main.gitCommit=$(GIT_COMMIT) \
		-X main.buildUser=$(BUILD_USER) \
		-X main.buildDate=$(BUILD_DATE)" \
		-gcflags="all=-trimpath=$(GOPATH)/src" \
		-asmflags="all=-trimpath $(GOPATH)/src" \
		cmd/authdbctl/*.go
	@./bin/authdbctl --version
	@./bin/authdbctl --help
	@echo "$@: complete"

linter:
	@echo "Running lint checks"
	@golint -set_exit_status ./...
	@echo "$@: complete"

gtest:
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./*.go
	@echo "$@: complete"

test: envvar covdir linter gtest coverage
	@echo "$@: complete"

ctest: covdir linter
	@richgo version || go get -u github.com/kyoh86/richgo
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./...
	@echo "$@: complete"

covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage
	@echo "$@: complete"

coverage:
	@#go tool cover -help
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./...
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"
	@echo "$@: complete"

docs:
	@mkdir -p .doc
	@go doc -all > .doc/index.txt
	@cat .doc/index.txt
	@echo "$@: complete"

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/
	@echo "$@: complete"

qtest:
	@echo "Perform quick tests ..."
	@#go test -v -run TestVersioned *.go
	@#go test -v -run TestNewID *.go
	@#time richgo test -v -run TestNewPublicKey *.go
	@#time richgo test -v -run TestNewUser *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewPublicKey *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewCode pkg/qr/*.go
	@time richgo test -v -coverprofile=.coverage/coverage.out -run TestDatabaseAuthentication *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewMfaToken *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out internal/tag/*.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run "Test.*Database.*" *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run "TestDatabaseGetUsers" *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out internal/tag/*.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewEmailAddress *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewRole *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewPassword *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewName *.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestDatabasePolicy *.go
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"
	@#go tool cover -func=.coverage/coverage.out | grep database
	@echo "$@: complete"


dep:
	@echo "Making dependencies check ..."
	@golint || go get -u golang.org/x/lint/golint
	@go get -u github.com/kyoh86/richgo
	@versioned || go get -u github.com/greenpau/versioned/cmd/versioned
	@echo "$@: complete"

license:
	@versioned || go get -u github.com/greenpau/versioned/cmd/versioned
	@for f in `find ./ -type f -name '*.go'`; do versioned -addlicense -copyright="Paul Greenberg greenpau@outlook.com" -year=2020 -filepath=$$f; done
	@echo "$@: complete"

release:
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m 'updated VERSION file'
	@versioned -sync database.go
	@git add database.go
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@#echo "git push --delete origin v$(APP_VERSION)"
	@#echo "git tag --delete v$(APP_VERSION)"

logo:
	@mkdir -p assets/docs/images/
	@gm convert -background black -font Bookman-Demi \
		-size 640x320 "xc:black" \
		-draw "fill white gravity center text 0,0 'Go\nidentity'" \
		-pointsize 96 \
		assets/docs/images/logo.png
	@echo "$@: complete"
