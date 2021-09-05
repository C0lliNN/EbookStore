package internal

type Bar struct {
	Name string
}

func NewBar() Bar {
	return Bar{Name: "Some-name"}
}

type Foo struct {
	Bar Bar
}

func NewFoo(bar Bar) Foo {
	return Foo{Bar: bar}
}