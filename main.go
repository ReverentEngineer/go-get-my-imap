package main

import (
  "log"
)

func main() {

  config := IMAPConfig{}


  if err := Configure(&config); err != nil {
    log.Fatal(err)
  }

  remote, err := IMAPConnect(config)

  if err != nil {
    log.Fatal(err)
  }
  local, err := MaildirInit(config.Maildir)

  if err != nil {
    log.Fatal(err)
  }

  mailbox_sync(remote, local)
}
