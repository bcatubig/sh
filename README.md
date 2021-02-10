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
        sh.WithArgs(
            "plan", 
            "-no-color",
            "-out",
            "plan.tfplan",
            "-detailed-exit-code", 
            tfDir,
        ),
        // We should expect to see code 2 returned
        sh.WithExpectedReturnCode(2),
        // Stream output to stdout as well
        sh.WithWriters(os.Stdout),
    )
    
    // result will always be populated with stdout/stderr
    result, err := c.Run()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result.Output.String())
    fmt.Println(result.ReturnCode)
}
```
