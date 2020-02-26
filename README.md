# godooj

Golang library to connect to Odoo via JSON-RPC

## Install

```sh
go get gitlab.com/ausuwardi/godooj
```

## Use

Create $HOME/.odoo.toml

```toml
[servers.testing]
server = "http://localhost:8069"
database = "demo"
user = "admin"
password = "admin"

[servers.production]
server = "https://demo.odoo.com"
database = "demo"
user = "admin"
password = "secret"
```

Then in your go code, use the library to connect:

```go
package main

import (
    "log"

    odoo "github.com/ausuwardi/godooj"
)

func main() {
    client, err := odoo.ClientConnect("testing")
    if err != nil {
        log.Error("Error connecting to Odoo")
        return
    }

    records, err := client.SearchRead(
        "product.product",
        odoo.List{
            odoo.List{"active", "=", true},
            odoo.List{"name", "like", "apple"},
        },
        []string{"id", "name", "price"},
    )
    if err != nil {
        log.Error("Error retrieving data")
        return
    }

    for _, rec := range records {
        id, _ := odoo.IntField(rec, "id")
        name, _ := odoo.StringField(rec, "name")
        price, _ := odoo.FloatField(rec, "price")

        log.Infof("Product %d - %s - %f", id, name, price)
    }
}
```

