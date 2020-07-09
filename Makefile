LDFLAGS ?= -s -w
BINARY=bin/s3fetch

SOURCES := s3fetch.go

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags="$(LDFLAGS) -X 'main.Revision=$$(git rev-parse HEAD)' -X 'main.Version=dev-build'" -o $(BINARY) $(SOURCES)

clean:
	$(RM) -rf bin/*

fmt: $(SOURCES)
	go fmt $(SOURCES)

.PHONY: clean fmt
