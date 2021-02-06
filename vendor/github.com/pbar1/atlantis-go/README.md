# Atlantis SDK for Go

[![GoDoc](https://godoc.org/github.com/pbar1/atlantis-go?status.svg)](https://godoc.org/github.com/pbar1/atlantis-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pbar1/atlantis-go)](https://goreportcard.com/report/github.com/pbar1/atlantis-go)

Go library for creating [custom `run` commands][1] for [Atlantis][2].

```sh
go get github.com/pbar1/atlantis-go
```

## Usage

```go
package main

import (
	"log"

	atlantis "github.com/pbar1/atlantis-go"
)

func main() {
	step, err := atlantis.NewRunStep()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Pull request number: %d\n", step.PullNum)
	log.Printf("Terraform plan file: %s\n", step.Planfile)
}
```

[1]: https://www.runatlantis.io/docs/custom-workflows.html#reference
[2]: https://www.runatlantis.io
