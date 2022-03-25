# GoSQL ![](https://github.com/twharmon/gosql/workflows/Test/badge.svg)
Query builder with some handy utility functions.

## Documentation
For full documentation see the [pkg.go.dev](https://pkg.go.dev/github.com/twharmon/gosql?tab=doc).

## Examples
```go
// Open database and create connection
sqliteDB, _ := sql.Open("sqlite3", "my-db.sql")
db := gosql.New(sqliteDB)

// Define a struct that includes a primary key
type User struct {
    ID       int `gosql:"primary"`
    Email    string
    IsActive bool
}

// Insert a row into the table
db.Insert(&User{
    ID: 1,
    Email: "gopher@example.com",
    IsActive: true,
})

// Select a row from the table
var user User
db.Select("*").Where("id = ?", 1).Get(&user)

// Update the row in the table
user.Email = "gosql@example.com"
db.Update(&user)

// Delete the row from the table
db.Delete(&user)
```

## Benchmarks
```
BenchmarkInsert-10            	    5637	    209484 ns/op	     448 B/op	      23 allocs/op
BenchmarkUpdate-10            	   90866	     12887 ns/op	     576 B/op	      27 allocs/op
BenchmarkSelect-10            	   90318	     13125 ns/op	     768 B/op	      41 allocs/op
BenchmarkSelectMany-10        	   12435	     96761 ns/op	   10640 B/op	     838 allocs/op
BenchmarkSelectManyPtrs-10    	   10000	    100512 ns/op	   11248 B/op	     938 allocs/op
```

## Contribute
Make a pull request
