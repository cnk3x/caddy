# libdns-westcn
west.cn (西部数码) provider for libdns

## Example

```go
import (
	westcn "github.com/cnk3x/libdns-westcn"
	"github.com/libdns/libdns"
)

ctx := context.TODO()

zone := "example.com."

// configure the DNS provider (choose any from github.com/libdns)
provider := &westcn.Provider{Username: "user", Password: "api_password"}

// list records
recs, err := provider.GetRecords(ctx, zone)

// create records (AppendRecords is similar)
newRecs, err := provider.SetRecords(ctx, zone, []libdns.Record{
	Type:  "A",
	Name:  "sub",
	Value: "1.2.3.4",
})

// delete records (this example uses provider-assigned ID)
deletedRecs, err := provider.DeleteRecords(ctx, zone, []libdns.Record{
	ID: "foobar",
})

// no matter which provider you use, the code stays the same!
// (some providers have caveats; see their package documentation)
```