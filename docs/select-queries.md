# Select Queries

## Setup

Define a type.
```
type User struct {
	ID       int
	Name     string
	IsActive bool
}

DB.Register(User{})
```


## Usage

Select all fields.
```
var user User
DB.Select("*").Get(&user)
```

Select some fields.
```
var user User
DB.Select("id", "name").Get(&user)
```

Select with a where clause.
```
var user User
DB.Select("*").Where("is_active = ?", true).Get(&user)
```

Select many rows.
```
var users []User
DB.Select("*").Where("is_active = ?", true).Get(&users)
```

