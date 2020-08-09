# retry-go [![GoDoc](https://godoc.org/github.com/megapctr/retry-go?status.svg)](http://godoc.org/github.com/megapctr/retry-go)

Minimal retry library

## Usage

Basic example:
```go
var res *http.Response
var err error
for r := retry.Exp(5, time.Second); r.Next(err); {
  res, err = http.Get("https://example.com")
}
if err != nil {
  // ...
}
```

Reusable retryer with custom DurationFunc
```go
var retryF = retry.Factory(3, func (i time.Duration) time.Duration {
  return 15 * i * time.Second
})

var res *http.Response
var err error
for r := retryF(); r.Next(err); {
  res, err = http.Get("https://example.com")
}
if err != nil {
  // ...
}
```
