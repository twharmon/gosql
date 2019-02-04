package gosql

import "reflect"

// Post -> Comment
type oneToMany struct {
	field reflect.StructField // Comments []*Comment
}

// Comment -> Post
type manyToOne struct {
	field  reflect.StructField // Post *Post
	column string              // "post_id"
}

// Post -> Tag
type manyToMany struct {
	field  reflect.StructField // Tags []*Tag
	table  string              // "post_tag"
	column string              // "tag_id"
}
