package main

import "fmt"

func main() {
	f()
	fmt.Println("Returned normally from f.")
}

func f() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	fmt.Println("Calling g.")
	g(0)
	fmt.Println("Returned normally from g.")
}

func g(i int) {
	if i > 3 {
		fmt.Println("Panicking!")
		panic(fmt.Sprintf("%v", i))
	}
	defer fmt.Println("Defer in g", i)
	fmt.Println("Printing in g", i)
	g(i + 1)
}

//-------------------------------------------------
//
//func explode() {
//	// Cause a panic.
//	panic("WRONG")
//}
//
//func main() {
//	// Handle errors in defer func with recover.
//	defer func() {
//		if err := recover(); err != nil {
//			// Handle our error.
//			fmt.Println("FIX")
//			fmt.Println("ERR", err)
//		}
//	}()
//	// This causes an error.
//	explode()
//	fmt.Println("reach here ") // 达不到
//}
//

//func explode() {
//	// Cause a panic.
//	panic("WRONG")
//
//}
//
//func throwPanic(f func()) {
//	defer func() {
//		if err := recover(); err != nil {
//			// Handle our error.
//			fmt.Println("FIX")
//			fmt.Println("ERR", err)
//		}
//	}()
//	f()
//	fmt.Println(" finish")
//}
//
//func main() {
//	throwPanic(explode)
//	fmt.Println("reach here ") //可达
//}
