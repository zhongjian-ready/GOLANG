package main

import (
	"fmt"
	"unsafe"
)

func main() {
    fmt.Println("--- Basic Types ---")
    fmt.Printf("bool:    %d bytes\n", unsafe.Sizeof(bool(true)))
    fmt.Printf("int:     %d bytes\n", unsafe.Sizeof(int(0)))
    fmt.Printf("int8:    %d bytes\n", unsafe.Sizeof(int8(0)))
    fmt.Printf("int16:   %d bytes\n", unsafe.Sizeof(int16(0)))
    fmt.Printf("int32:   %d bytes\n", unsafe.Sizeof(int32(0)))
    fmt.Printf("int64:   %d bytes\n", unsafe.Sizeof(int64(0)))
    fmt.Printf("float64: %d bytes\n", unsafe.Sizeof(float64(0)))

    fmt.Println("\n--- Reference Types (Header size) ---")
    var p *int
    var s string = "Hello World This text is long but the header is small"
    var sl []int = []int{1, 2, 3, 4, 5}
    
    fmt.Printf("pointer: %d bytes\n", unsafe.Sizeof(p))
    fmt.Printf("string:  %d bytes\n", unsafe.Sizeof(s))
    fmt.Printf("slice:   %d bytes\n", unsafe.Sizeof(sl))
}