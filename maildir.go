package main

import (
  "io/ioutil"
  "io"
  "os"
  "errors"
  "fmt"
  "crypto/md5"
  "github.com/emersion/go-imap"
)

type Maildir struct  {
  RootPath string
}

func MaildirInit(path string) (*Maildir, error) {

  fileInfo, err := os.Stat(path)
  if os.IsNotExist(err) {
    if err := os.Mkdir(path, 0700); err != nil {
      return nil, err
    }
  } else if err != nil {
    return nil, err
  } else if fileInfo.IsDir() == false {
    return nil, errors.New("Maildir path is not a directory.")
  }

  return &Maildir{RootPath: path}, nil
}

func GenerateUniqueName(folder string, message *imap.Message) (string, error) {
  uid := (*message).Uid
  secs := (*message).Envelope.Date.Unix()
  nsecs := (*message).Envelope.Date.UnixNano()
  flags := ""
  hostname, err := os.Hostname()
  folder_md5 := fmt.Sprintf("%x",md5.Sum([]byte(folder)))
  if err != nil {
    return "", err
  }
  return fmt.Sprintf("%d_%d.%s,U=%d,FMD5=%s,%s", secs, nsecs, hostname, uid, folder_md5, flags), nil
}

func (m *Maildir) GetMailboxes() (map[string]bool, error) {
  var err error
  mailboxes := make(map[string]bool)

  fileInfo, err := ioutil.ReadDir(m.RootPath)
  if err != nil {
    goto Error
  }

  for _, file := range fileInfo {
    mailboxes[file.Name()] = true
  }

  Error:
  return mailboxes, err
}

func (m *Maildir) CreateMailbox(name string) error {

  if err := os.Mkdir(fmt.Sprintf("%s/%s", m.RootPath, name), 0700); err != nil && !os.IsExist(err) {
    return err
  }

  if err := os.Mkdir(fmt.Sprintf("%s/%s/tmp", m.RootPath, name), 0700); err != nil && !os.IsExist(err) {
    return err
  }


  if err := os.Mkdir(fmt.Sprintf("%s/%s/new", m.RootPath, name), 0700); err != nil && !os.IsExist(err) {
    return err
  }


  if err := os.Mkdir(fmt.Sprintf("%s/%s/cur", m.RootPath, name), 0700); err != nil && !os.IsExist(err) {
    return err
  }

  return nil
}

func (m *Maildir) GetMessages(mailbox string) ([]*imap.Message, error) {
  _, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", m.RootPath, mailbox))
  if err != nil {
    return nil, err
  }

  return nil, nil
}

func (m *Maildir) WriteMessage(mailbox string, message *imap.Message) {
  name,_ := GenerateUniqueName(mailbox, message)
  filepath := fmt.Sprintf("%s/%s/new/%s", m.RootPath, mailbox, name)
  if _, err := os.Stat(filepath); os.IsNotExist(err) {
    f,_ := os.Create(filepath)
    defer f.Close()
    var section imap.BodySectionName
    r := (*message).GetBody(&section)
    io.Copy(f, r)
  }
}

