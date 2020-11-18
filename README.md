# form3-go

Form3 accounts API take home test

go-form3 is a Go client library for accessing the [Form3 API](https://api-docs.form3.tech/)

By [Glen Thomas](glen-thomas.com) [glen.thomas@outlook.com](mailto:glen.thomas@outlook.com)

I have never used Go before. I currently mainly code with C# and TypeScript.

When I needed to know how to do something (e.g. initialize a variable, create a function, use the go cli, etc.) I searched online. Stackoverflow had many answers I needed. I also checked out some open source Go HTTP client repositories to learn by example.

I am not familiar with commonly used Go code styling (indentation, alignment, comments, naming conventions etc.) or project structure. I have looked at some example repositories to get an idea of common practice. I may have laid things out in an unusual way (do people generally order files as; package, imports, types, functions?). I haven't looked at using any linter, but given more time would add something in. My IDE does some auto-formatting for me.

The test criteria specified using the fake API for integration testing. I also wanted to have some basic unit tests for fast feedback so am using [httptest](https://pkg.go.dev/net/http/httptest) to intercept HTTP requests and provide static responses to verify the outgoing HTTP request and handling of the response works as expected.

I saw in the [API documentation](https://api-docs.form3.tech/api.html#introduction-and-api-conventions) that some headers are required in all requests so have a shared function for generating a new request object and adding those headers in.

The Docker compose fails on the first up because the account api is not ready to receive requests. I am handling this with a loop in the shell script that probes the accounts endpoint of the API for a response before beginning the tests.

I initially used the golang:alpine image for fastest download/run speed but hit [this issue](https://github.com/golang/go/issues/28065), so used golang:stretch instead.

## Usage

A basic example to list accounts:

```go
import "form3.tech/go-form3/form3"

func main() {
	client := form3.NewClient(nil)
	accounts, _, err := client.Accounts.List(context.Background(), &form3.ListOptions{PageNumber: 1, PageSize: 50})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
```

See [Examples](/examples).

## Testing

To run unit tests `go test -run 'Unit'`

To run integration tests `docker-compose up`
