package main

import "fmt"

// Layer 1 interface
type Layer1 interface {
	Method1() string
}

// Layer 2 interface
type Layer2 interface {
	Layer1
	Method1() string // Same method name as Layer1
}

// Concrete type implementing both layers
type MyType struct{}
type MyType2 struct{}

func (t MyType) Method1() string {
	return "Layer1 implementation"
}

// Corrected method name to match Layer1's Method1
func (t MyType2) Method1() string { // Same name as Layer1's Method1
	return "Layer2 implementation"
}

func callMe(myObj Layer2) {
	myObj.Method1()
}

func main() {
	var obj Layer2 = MyType{}

	// Method shadowing: Method1 from Layer2 takes precedence
	fmt.Println(obj.Method1()) // Output: Layer1 implementation


	// Accessing Layer1's Method2 is not possible:

	callMe(obj)

	if layer1Obj, ok := obj.(Layer2); ok {
		fmt.Println(layer1Obj.Method1()) // Output: Layer1 implementation
	}
}