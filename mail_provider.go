package main

import (
  "github.com/emersion/go-imap"
)

type MailProvider interface {
  GetMailboxes() (map[string]bool, error)
  CreateMailbox(string) error
  GetMessages(string) ([]*imap.Message, error)
}
