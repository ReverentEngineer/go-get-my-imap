package main

import (
  "flag"
  "gopkg.in/yaml.v2"
  "io/ioutil"
)

type IMAPConfig struct {
  Hostname string
  Tls bool
  Username string
  Password string
  Maildir string
}

func Configure(config *IMAPConfig) (error) {
  configPathPtr := flag.String("config", "config.yml", "The configuration file to perform authentication.")
  flag.Parse()

  data, err := ioutil.ReadFile(*configPathPtr)

  if (err != nil) {
    return err
  }
  err = yaml.Unmarshal(data, config);

  if err != nil {
    return err
  }

  return nil
}
