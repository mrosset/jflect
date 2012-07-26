package main

// Field data type
type Field struct {
	name  string
	gtype string
	tag   string
}

// Simplifies Field construction
func NewField(name, gtype string) Field {
	return Field{goField(name), gtype, goTag(name)}
}

// Provides Sorter interface so we can keep field order
type FieldSort []Field

func (s FieldSort) Len() int { return len(s) }

func (s FieldSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s FieldSort) Less(i, j int) bool {
	return s[i].name < s[j].name
}
