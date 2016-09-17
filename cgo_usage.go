package main

// typedef int (*intFunc) ();
//
// int
// bridge_int_func(intFunc f)
// {
//		return f();
// }
//
// int fortytwo()
// {
//	    return 42;
// }
import "C"
import "fmt"

func main() {
	f := C.intFunc(C.fortytwo)
	fmt.Println(int(C.bridge_int_func(f)))
	// Output: 42
}

// // #cgo CFLAGS: -DPNG_DEBUG=1
// // #cgo amd64 386 CFLAGS: -DX86=1
// // #cgo LDFLAGS: -lpng
// // #include <png.h>
// import "C"
