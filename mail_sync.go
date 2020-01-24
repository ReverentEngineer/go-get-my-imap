package main

import (
  "sync"
  "github.com/emersion/go-imap"
)

func DownloadMailboxes(mailboxes []string, remote *IMAPFacade, local *Maildir)  {
  for _,mailbox := range mailboxes {
    messages,_ := (*remote).GetMessages(mailbox)
    var wg sync.WaitGroup
    for message := range messages {
      wg.Add(1)
      go func(mailbox string, message *imap.Message, wg *sync.WaitGroup) {
        (*local).WriteMessage(mailbox, message)
        wg.Done()
      }(mailbox, message, &wg)
    }
  }
}

func CreateMailboxes(mailboxes []string, local *Maildir) { 
  for _,mailbox := range mailboxes {
    (*local).CreateMailbox(mailbox)
  }
}


func mailbox_sync(remote *IMAPFacade, local *Maildir) {
  remote_mailboxes,_ := (*remote).GetMailboxes()
  CreateMailboxes(remote_mailboxes, local)
  DownloadMailboxes(remote_mailboxes, remote, local)
}
