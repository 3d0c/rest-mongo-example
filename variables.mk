### project name
BINARY		:= lyre-be

### version
MAJOR            := 0
MINOR            := 0
PATCH            := 1
PRODUCT_VERSION  := ${MAJOR}.${MINOR}.${PATCH}
COMMIT			 := $(shell git rev-parse --short HEAD)
FULL_VERSION     := ${PRODUCT_VERSION}-${COMMIT}
DATE			 := $(shell date +%FT%T%z)

### go vars
GO_ENV_VARS			= GO111MODULE=on CGO_ENABLED=0 GOPROXY=https://proxy.golang.org,direct
BUILD_FLAGS	= -ldflags "-w -s -X github.com/teal-seagull/lyre-be-v4/cmd/lyre-be.version=${FULL_VERSION} -X github.com/teal-seagull/lyre-be-v4/cmd/lyre-be.date=${DATE}"

### variables for testing
COVERAGE_OUT		:= coverage.out
COVERAGE_HTML		:= coverage.html
KEYS_PATH			:= tmp/keys
TMP					:= /tmp
UTS_KEYS			:= tmp/data/keys

### Platforms
DARWIN	:= darwin
LINUX	:= linux

### Path for install
DESTDIR			:= ""
BIN_DIR			:= ${DESTDIR}/usr/bin
SYSTEMD_DIR		:= ${DESTDIR}/etc/systemd/system
PROD_DIR		:= ${DESTDIR}/opt/lyre-be
