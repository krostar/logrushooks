# Caller Hook

## Usage

The quickest and easiest way to use it is to simply add a hook with default parameters:

```go
package main

import (
    "github.com/krostar/logrushooks"
    "github.com/sirupsen/logrus"
)

func init() {
    logrus.AddHook(logrushooks.NewCallerHook(logrushooks.DefaultCallerOptions))
}

func main() {
    logrus.Info("Hello World")
    // INFO[0000] Hello World            caller="/Users/krostar/go/src/github.com/krostar/mypackage/main.go:13"
}
```

You can also change the name of the key and remove the app package path from the value:

```go
package main

import (
    "github.com/krostar/logrushooks"
    "github.com/sirupsen/logrus"
)

func init() {
    logrus.AddHook(logrushooks.NewCallerHook(&logrushooks.CallerOptions{
        AppPackage:      "github.com/krostar/mypackage",
        CallerKey:       "custom-caller-name",
    }))
}

func main() {
    logrus.Info("Hello World")
    // INFO[0000] Hello World            custom-caller-name="main.go:16"
}
```

For more detailed example you can also check the `caller_test.go` file.

## Impact on logging time

```sh
# Benchmark with a benchtime of 1mn
# on a Mac Book Pro 2017
BenchmarkLogWithoutCallerHook-4   50000000   1523 ns/op
BenchmarkLogWithCallerHook-4      20000000   4523 ns/op
```
