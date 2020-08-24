# sh
A library for easily running shell commands

## Install

```shell script
go get -u github.com/bcatubig/sh
```

## Usage

```go
package main
import (
    "fmt"
    "log"
    "os"
    
    "github.com/bcatubig/sh"
)

func main() {
    tfDir := "./infra/terraform"
    c := sh.NewCommand("terraform", 
        sh.Args(
            "init", 
            "-no-color", 
            "-detailed-exit-code", 
            tfDir,
        ),
        sh.ExpectedReturnCode(2),
        sh.Writers(os.Stdout),
    )
    
    output, err := c.Run()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(output.String())
}
```
