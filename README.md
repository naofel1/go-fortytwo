# go-fortytwo

[![GitHub tag (latest
SemVer)](https://img.shields.io/github/v/tag/naofel1/go-fortytwo?label=go%20module)](https://github.com/naofel1/go-fortytwo/tags)
[![Go
Reference](https://pkg.go.dev/badge/github.com/naofel1/go-fortytwo.svg)](https://pkg.go.dev/github.com/naofel1/go-fortytwo)
[![GitHub](https://img.shields.io/github/license/naofel1/go-fortytwo)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/naofel1/go-fortytwo)](https://goreportcard.com/report/github.com/naofel1/go-fortytwo)

**go-fortytwo** is a client for the [42
API](https://api.intra.42.fr/apidoc), written in
[Go](https://golang.org/).

## Installation

```sh
go get github.com/naofel1/go-fortytwo
```

## Usage

To obtain an API key, follow 42 documentation [getting started
guide](https://api.intra.42.fr/apidoc/guides/getting_started).

```go
import "github.com/naofel1/go-fortytwo"

(...)

client := fortytwo.NewClient(ctx, "client_id", "client_secret", "redirect_url", []string{"public"})

achivements, err := client.List(context.Background(), &fortytwo.AchievementQueryRequest{
      Pagination: &fortytwo.Pagination{
         Cursor:   1,
         PageSize: 10,
      },
   })
if err != nil {
    // Handle error...
}
```

👉 Check out the docs on
[pkg.go.dev](https://pkg.go.dev/github.com/naofel1/go-fortytwo) for a complete
reference and the [examples](/examples) directory for more example code.

## Status

The FortyTwo API golang package is in alpha. If you want to contribute you are
welcome to do so. Please read the [contributing guidelines](CONTRIBUTING.md) before you start.

## License

[MIT License](LICENSE)
