# GoSQL [![Build Status](https://travis-ci.com/twharmon/gosql.svg?branch=master)](https://travis-ci.com/twharmon/gosql) [![Coverage Status](https://coveralls.io/repos/github/twharmon/gosql/badge.svg?branch=master)](https://coveralls.io/github/twharmon/gosql?branch=master)
Query builder with some handy utility functions.

## Install
`go get github.com/twharmon/gosql`

## Usage
```
// connect to a database
sqlDB, _ := sql.Open(/* your database connection */)
db := gosql.Conn(sqlDB)

// define a struct to be associated with a table in your database
// table names MUST be lower_snake_case singular of struct name
type User struct {
    // ID must be first field in struct and must be int64
    ID int64
    Name string
    Active bool
}

// you must register all structs
gosql.Register(User{})


// now you are ready to go

insert a new user
newUser := User{
    Name: "New Gopher",
    Active: true,
}
db.Insert(&newUser)
// newUser.ID is set after inserted into database

// select a user
var user User
db.Select("*").Where("id = ?", 1).To(&user)

// update the user
user.Name = "Gopher"
db.Update("user").Set("name = ?", user.Name).Where("id = ?", user.ID).Exec()

// delete the user
db.Delete("user").Where("id = ?", user.ID).Exec()

// count all users
count, _ := db.Count("user", "*").Exec()

// count all active users
count, _ := db.Count("user", "*").Where("active = ?", true).Exec()

```

For full documentation see [godoc](https://godoc.org/github.com/twharmon/gosql).

## Contribute
Make a pull request
