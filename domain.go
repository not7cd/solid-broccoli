package main

type Keyword struct {
	ID   int64
	Name string
}

type Sentence struct {
	ID   int64
	KeywordID   int64
	Value string
}