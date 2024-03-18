# Example: Generating HTML output

The template file needs to follow go's
[html/template](https://pkg.go.dev/html/template) rules. A sample file is
located [here](./example.html).

ecsv provides the output data represented via the struct `HTMLData`:

```go
type HTMLDataRow struct {
	Data   []string
	InSync bool
}
type HTMLData struct {
	Title     string
	Columns   []string
	Rows      []HTMLDataRow
	Errors    *[]error
	Timestamp string
}
```

You will primarily be interested in iterating over the field `Rows`.
`HTMLDataRow.InSync` signifies whether the versions for a system are in sync
or not, and you can leverage that to render a row in a particular style.
