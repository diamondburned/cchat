package repository

type Struct struct {
	Comment
	Name   string
	Fields []StructField
}

type StructField struct {
	Comment
	NamedType
}
