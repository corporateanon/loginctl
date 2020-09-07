# loginctl

Retrieves information about current users activity in Linux desktop systems

```go
package main

import (
	"log"

	"github.com/corporateanon/loginctl"
)

func main() {
	lctl, err := loginctl.NewFromRegularUsers()
	if err != nil {
		log.Panicln(err)
	}

	log.Println(lctl.GetSessionInfo())

}
```
