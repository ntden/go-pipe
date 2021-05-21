## go-pipe

Go module allowing to concisely execute functions one after another, where a function's outputs are the next one's inputs.  
This **should not be used in situations where performance is important** because it uses reflection and thus will slow down your program compared to calling functions manually.

### Install

```sh
$ go get -u github.com/ntden/go-pipe
```

### Usage

```go
package main

import "github.com/ntden/go-pipe"

func main() {
    p, err := pipe.New(
        func(int a, int b) float32 { ... },
        func(float32 a) (string, error) { ... },
        func(string a) { ... }
    )

    _, err = pipe.Execute(10, 11)
}

```
