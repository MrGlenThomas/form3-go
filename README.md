# form3-go

Form3 accounts API take home test

go-form3 is a Go client library for accessing the [Form3 API](https://api-docs.form3.tech/)

By [Glen Thomas](glen-thomas.com) [glen.thomas@outlook.com](mailto:glen.thomas@outlook.com)

I have never used Go before. I currently mainly code with C# and TypeScript.

I am not familiar with commonly used go code styling (indentation, alignment, comments, etc.) or project structure. I have looked at some example repositories to get an idea of common practice.

The test criteria specified using the fake API for integration testing. I also wanted to have some basic unit tests for fast feedback so am using httptest to intercept HTTP requests and provide static responses.

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
