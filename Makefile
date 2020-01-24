
GO_SRCS = $(wildcard *.go)

get-my-imap: $(GO_SRCS)
	go build -o $@

clean: get-my-imap
	rm get-my-imap
.PHONY: clean

run: get-my-imap
	./get-my-imap
.PHONY: run
