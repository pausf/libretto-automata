BIN := bin/libretto
PKG := ./cmd/libretto

# The version the binary reports, taken from git rather than kept in a source file.
#
# A hardcoded constant drifts from the tag the moment you forget to bump it, and it
# drifts silently — the binary keeps claiming a version nobody released. Asking git
# means the answer cannot be wrong.
#
#   v0.2.0            exactly on a tag, clean tree
#   v0.2.0-3-gabc123  three commits past it
#   v0.2.0-3-gabc123-dirty   with uncommitted changes, which is worth seeing
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

# Where `make link` puts the command, and under which names.
#
# `libretto` is the one to type. The long form is linked too because it is what the
# repository is called, and somebody who remembers the project rather than the
# command should still find it.
#
# The binary prints whatever name it was invoked as, so every link describes itself
# correctly with no extra work. Override either:  make link NAMES=lb
PREFIX ?= $(HOME)/.local/bin
NAMES  ?= libretto libretto-automata

.PHONY: build test fmt vet preview clean link unlink

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)

test:
	go test ./...

# Excludes the slow paths: real-git integration and teatest flows.
test-short:
	go test -short ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

# The panel, once, with colour forced. Also the way to eye the font fallback:
#   LIBRETTO_ASCII=safe make preview
preview: build
	@$(BIN) preview

# Put `lib` on PATH. A symlink, not a copy — so `make build` updates the installed
# command and there is never a stale binary pretending to be the current one.
#
# Refuses to overwrite anything that is not already our own link. A tool whose whole
# promise is "never clobbers what it did not create" does not get to make an
# exception for itself.
link: build
	@mkdir -p $(PREFIX)
	@for n in $(NAMES); do \
		dst="$(PREFIX)/$$n"; \
		if [ -e "$$dst" ] && [ "$$(readlink "$$dst")" != "$(CURDIR)/$(BIN)" ]; then \
			echo "refusing  $$dst — exists and is not ours"; \
			echo "          it is: $$(readlink "$$dst" 2>/dev/null || echo 'a real file')"; \
			continue; \
		fi; \
		ln -sfn "$(CURDIR)/$(BIN)" "$$dst"; \
		echo "linked    $$dst"; \
	done
	@command -v libretto >/dev/null 2>&1 \
		&& echo "          \`libretto\` is on PATH" \
		|| echo "          NOTE: $(PREFIX) is not on your PATH"

unlink:
	@for n in $(NAMES); do \
		dst="$(PREFIX)/$$n"; \
		if [ "$$(readlink "$$dst" 2>/dev/null)" = "$(CURDIR)/$(BIN)" ]; then \
			rm "$$dst"; echo "removed   $$dst"; \
		else \
			echo "skipped   $$dst — not our link"; \
		fi; \
	done

clean:
	rm -rf bin
