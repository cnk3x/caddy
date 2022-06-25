# Linode for `libdns`

[![Go Reference](https://pkg.go.dev/badge/test.svg)](https://pkg.go.dev/github.com/tosie/libdns-linode)

This package implements the [libdns interfaces](https://github.com/libdns/libdns) for the [Linode Domains API](https://www.linode.com/docs/api/domains/).

## Authenticating

To authenticate you need to supply a Linode [Personal Access Token](https://cloud.linode.com/profile/tokens).

## Example

Here's a minimal example of how to get all DNS records for a zone. See also: [provider_test.go](https://github.com/tosie/libdns-linode/blob/master/provider_test.go)

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tosie/libdns-hetzner"
)

func main() {
	token := os.Getenv("LIBDNS_LINODE_TOKEN")
	if token == "" {
		fmt.Printf("LIBDNS_LINODE_TOKEN not set\n")
		return
	}

	zone := os.Getenv("LIBDNS_LINODE_ZONE")
	if token == "" {
		fmt.Printf("LIBDNS_LINODE_ZONE not set\n")
		return
	}

	p := &linode.Provider{
		APIToken: token,
	}

	records, err := p.GetRecords(context.WithTimeout(context.Background(), time.Duration(15*time.Second)), zone)
	if err != nil {
        fmt.Printf("Error: %s", err.Error())
        return
	}

	fmt.Println(records)
}

```
