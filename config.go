package main

const PrefixCmd = '.'

type Query byte

const (
	QueryAdd     Query = '+'
	QueryRandom  Query = '.' // no arg
	QueryInspect Query = '?' // no arg
)
