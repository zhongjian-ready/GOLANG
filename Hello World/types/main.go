package main

import (
	"fmt"
	"unsafe"
)

// --- Struct for demonstrating aggregate types ---
type Person struct {
	Name string
	Age  int
}

// --- Interface definition for demonstrating interface types ---
type Speaker interface {
	Speak()
}

// Person implements the Speaker interface
func (p Person) Speak() {
	fmt.Println(p.Name, "says: Hello Interface!")
}

func main() {
	fmt.Println("================== 1. Basic Types ==================")
	// --- Boolean (bool) ---
	var b bool = true
	fmt.Printf("[Bool] Value: %v, Size: %d byte\n", b, unsafe.Sizeof(b))

	// --- Numeric Types ---
	// int size depends on the OS (usually 8 bytes on 64-bit systems)
	var i int = 100
	var f64 float64 = 3.14159
	// Complex number: Real + Imaginary parts
	var c128 complex128 = 1 + 2i
	
	fmt.Printf("[Int] Value: %d, Size: %d bytes\n", i, unsafe.Sizeof(i))
	fmt.Printf("[Float64] Value: %f\n", f64)
	fmt.Printf("[Complex128] Value: %v\n", c128)

	// --- Byte vs Rune ---
	// byte is an alias for uint8, usually used for ASCII or raw binary data
	var myByte byte = 'A' 
	// rune is an alias for int32, represents a Unicode Code Point, used for Chinese chars etc.
	var myRune rune = '国' 
	
	fmt.Printf("[Byte/uint8] Value: %v, Char: %c\n", myByte, myByte)
	fmt.Printf("[Rune/int32] Value: %v, Char: %c (Size: %d bytes)\n", myRune, myRune, unsafe.Sizeof(myRune))

	// --- String ---
	str := "Hello"
	fmt.Printf("[String] Value: %s, Header Struct Size: %d bytes (not content length)\n", str, unsafe.Sizeof(str))

	
	fmt.Println("\n================== 2. Aggregate Types ==================")
	// --- Array ---
	// Features: Fixed length, "Value Type". Assignment copies the entire array.
	var arr [3]int = [3]int{1, 2, 3}
	arrCopy := arr // Here a [Full Copy] occurs
	arrCopy[0] = 999 
	
	fmt.Printf("[Array] Original: %v (Unchanged)\n", arr)
	fmt.Printf("[Array] Copy:     %v (Only copy modified)\n", arrCopy)

	// --- Struct ---
	// Features: Also "Value Type". Assignment copies all fields.
	p1 := Person{Name: "Alice", Age: 30}
	p2 := p1 // Here a [Full Copy] occurs
	p2.Name = "Bob"

	fmt.Printf("[Struct] Original: %+v (Unchanged)\n", p1)
	fmt.Printf("[Struct] Copy:     %+v (Only copy modified)\n", p2)


	fmt.Println("\n================== 3. Reference Types ==================")
	// Features: Variable stores a "descriptor" (contains pointer). Assignment copies the descriptor, but points to the same underlying data.
	
	// --- Slice ---
	slice := []int{1, 2, 3}
	sliceCopy := slice // Copied slice header (ptr+len+cap), ptr points to same array // Here a [Shallow Copy] occurs
	sliceCopy[0] = 888 

	fmt.Printf("[Slice] Original: %v (Modified! Because underlying array is shared)\n", slice)

	// --- Map ---
	m := map[string]int{"key": 10}
	mCopy := m // Copied map pointer
	mCopy["key"] = 20

	fmt.Printf("[Map]   Original:  %v (Modified!)\n", m)

	// --- Pointer ---
	// Specifically used to store memory addresses
	num := 50
	ptr := &num // Get address
	*ptr = 100  // Modify target value via pointer

	fmt.Printf("[Pointer] Num modified via pointer: %d\n", num)

	// --- Channel ---
	ch := make(chan int)
	fmt.Printf("[Channel] This is a reference type, addr/descriptor: %v\n", ch)

	// --- Function ---
	myFunc := func() { fmt.Println("Hello Func") }
	fmt.Printf("[Function] Function can be assigned to var, type: %T\n", myFunc)


	fmt.Println("\n================== 4. Interface Types ==================")
	// --- Empty Interface (interface{}) ---
	// Can hold any type (similar to Java Object or C void*)
	var anyVal interface{}
	
	anyVal = 123
	fmt.Printf("[Interface{}] Hold Int: %v (Underlying Type: %T)\n", anyVal, anyVal)
	
	anyVal = "String Value"
	fmt.Printf("[Interface{}] Hold String: %v (Underlying Type: %T)\n", anyVal, anyVal)

	// --- Polymorphism ---
	var s Speaker = Person{Name: "SpeakerObj", Age: 18}
	fmt.Print("[Interface] Polymorphic Call: ")
	s.Speak() // Calls Person's Speak method
}
