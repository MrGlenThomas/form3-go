# form3-go

Form3 accounts API take home test

go-form3 is a Go client library for accessing the [Form3 API](https://api-docs.form3.tech/)

By [Glen Thomas](glen-thomas.com) [glen.thomas@outlook.com](mailto:glen.thomas@outlook.com)

I have never used Go before. I currently mainly code with C# and TypeScript.

When I needed to know how to do something (e.g. initialize a variable, create a function, use the go cli, etc.) I searched online. Stackoverflow had many answers I needed.

I am not familiar with commonly used go code styling (indentation, alignment, comments, naming conventions etc.) or project structure. I have looked at some example repositories to get an idea of common practice. I may have laid things out in a frustrating way (do people generally order files as package name, imports, types, functions?).

The test criteria specified using the fake API for integration testing. I also wanted to have some basic unit tests for fast feedback so am using httptest to intercept HTTP requests and provide static responses.

I saw in the [API documentation](https://api-docs.form3.tech/api.html#introduction-and-api-conventions) that some headers are required in all requests so have a shared function for generating a new request object and adding those headers in.

I have one large integration test that covers all operations. I should break this into smaller test functions that cover specific scenarios, but need to consider how the underlying database will affect the test design (i.e. the same data store is shared across parallel-executed tests).

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
