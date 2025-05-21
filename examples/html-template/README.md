# Example: Generating HTML output

The template file needs to follow go's
[html/template](https://pkg.go.dev/html/template) rules.
[Here's](../../internal/ui/assets/template.html) `ecsv`'s built in template.

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

The built in template generates an HTML file that looks like the following:

![ecsv-html](https://github.com/user-attachments/assets/dbde169a-3253-42cd-b5ff-0f2f99cecf58)
