# Fields Hook

## Usage

It can be simply used by adding a hook with fields as parameters:

```go
package main

import (
    "github.com/krostar/logrushooks"
    "github.com/sirupsen/logrus"
)

func init() {
    logrus.AddHook(logrushooks.NewFieldHook(&logrushooks.FieldOptions{
        "version": "4.2.1",
        "weather": 42,
    }))
}

func main() {
    logrus.Info("Hello World")
    // INFO[0000] Hello World            version="4.2.1" weather=42
}
```

## Impact on logging time

```sh
# Benchmark with a benchtime of 1mn
# on a Mac Book Pro 2017
BenchmarkLogWithoutCallerHook-4   50000000   1495 ns/op
BenchmarkLogWithFieldsHook-4      50000000   1966 ns/op
```
