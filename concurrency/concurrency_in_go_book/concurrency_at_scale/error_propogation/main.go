package main

import "fmt"

/*
Error Wrapping:

* Wrap errors only at module boundaries.
	- ex: func1 -> func2 -> func3 -> func4 -> func5 -> func6, if func6 returns an error, then wrap the error at func1. if all are in same boundary, becuase no need to log each internal function error, it depends on context. Sometimes we log that info
* Add meaningful context to errors (e.g., input data, high-level messages).


TODO- how log-store is working.

Define custom error types to ensure all errors are well-formed and structured.
*/

func main() {
	fmt.Println((1 << 9) | (1 << 8))
}
