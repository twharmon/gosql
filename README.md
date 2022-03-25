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

// Register all structs corresponding to a table in the database
db.Register(User{})

// Select a row from the table
var user User
db.Select("*").Get(&user)

// Update the row in the table
user.Email = "somenewemail@example.com"
db.Update(&user)
```

## Benchmarks
```
BenchmarkInsert-4        	    836143 ns/op	     400 B/op	      17 allocs/op
BenchmarkUpdate-4        	     22923 ns/op	     488 B/op	      18 allocs/op
BenchmarkSelect-4        	     24934 ns/op	     648 B/op	      26 allocs/op
BenchmarkSelectMany-4    	    127559 ns/op	    6568 B/op	     328 allocs/op
BenchmarkSelectManyPtrs-4	    130752 ns/op	    7976 B/op	     428 allocs/op
```

## Contribute
Make a pull request
