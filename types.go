package main

const (
	TagValueStyleCamel = "camel"
	TagValueStyleSnake = "snake"
	TagValueStyleGo    = "go"
	TagValueStyleUpper = "upper"
	TagValueStyleLower = "lower"
)

type TagDescribe struct {
	TagKey        string
	TagValueStyle string
}
