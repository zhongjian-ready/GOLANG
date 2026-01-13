package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var s string
	s = "hello world"

	var allType interface{}
	allType = s

	str, ok := allType.(string)
	if ok {
		fmt.Println("转换成功", str)
	} else {
		fmt.Println("转换失败")
	}

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Println("open /dev/tty failed", err)
		return
	}
	
	var r io.Reader
	r = tty

	var w io.Writer
	w = r.(io.Writer)

	w.Write([]byte("请输入内容: \n"))

	var buf [128]byte
	n, err := r.Read(buf[:])
	if err != nil {
		fmt.Println("read from tty failed", err)
		return
	}
	fmt.Printf("read %d bytes: %s\n", n, string(buf[:n]))
}
