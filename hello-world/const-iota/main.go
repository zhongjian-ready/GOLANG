package main

func main() {
	const (
		a = iota // 0
		b        // 1
		c        // 2
		d = "hello"
		e        // "hello"
		f = iota // 5
		g        // 6
	)

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)
	println(g)
}	