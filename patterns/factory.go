package pattern

import (
	"fmt"
)

type Creator struct {
	factory factory
}

func (c *Creator) Operation() {
	product := c.factory.factoryMethod()
	product.method()
}

type factory interface {
	factoryMethod() Product
}

type ConcreteCreator struct{}

func (c *ConcreteCreator) factoryMethod() Product {
	return new(ConcreteProduct)
}

type Product interface {
	method()
}

type ConcreteProduct struct{}

func (p *ConcreteProduct) method() {
	fmt.Println("ConcreteProduct.method()")
}
