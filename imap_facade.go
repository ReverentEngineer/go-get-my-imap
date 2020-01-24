package main

import (
  "github.com/emersion/go-imap/client"
  "github.com/emersion/go-imap"
)

type IMAPFacade struct {
  client *client.Client
}

func IMAPConnect(config IMAPConfig) (*IMAPFacade,error) {

  var err error
  var c *client.Client

  if config.Tls {
    c, err = client.DialTLS(config.Hostname, nil)
  } else {
    c, err = client.Dial(config.Hostname)
  }

  if err != nil {
    return nil, err
  }

  err = c.Login(config.Username, config.Password)

  if err != nil {
    return nil, err
  }

  return &IMAPFacade{client: c}, nil
}

func (r *IMAPFacade) GetMailboxes() ([]string, error) {
  mailboxes := make(chan *imap.MailboxInfo, 10)
  err := make(chan error, 1)
  go func () {
    err <- (*r.client).List("", "*", mailboxes)
  }()
  mailbox_list := make([]string, 0)
  for mailbox := range mailboxes {
    mailbox_list = append(mailbox_list, mailbox.Name)
  }
  return mailbox_list, <-err
}

func (r *IMAPFacade) CreateMailbox(name string) error {
  if err := (*r.client).Create(name); err != nil {
    return err
  }
  return nil
}

func (r *IMAPFacade) GetMessages(mailbox string) (<-chan *imap.Message, <-chan error) {
  mbox, err := (*r.client).Select(mailbox, false)
  messages := make(chan *imap.Message, 10)
  done := make(chan error, 1)
  if err != nil {
    done <- err
    return nil, done
  }

  from := uint32(1)
  to := mbox.Messages
  seqset := new(imap.SeqSet)
  seqset.AddRange(from, to)

  var section imap.BodySectionName
  items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchUid}
  go func() {
    done <- (*r.client).Fetch(seqset, items, messages)
  }()
  return messages, done
}
