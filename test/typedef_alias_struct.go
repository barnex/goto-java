package main

type A struct{ v int }
type B A

func (a A) AV()  { println(a.v) }
func (a *A) AP() { a.v++ }
func (b B) BV()  { println(b.v) }
func (b *B) BP() { b.v++ }

func main() {
	var a A
	a = A{}
	a = A{0}
	a = A{v: 0}
	a = A(struct{ v int }{0})

	var b B
	b = B(a) //b.set(a.data()) // or b.set(a.v1, a.v2, ...)?
	b.BP()   // inc b
	a.AV()   // unchanged
	b.BV()   // incremented
}
