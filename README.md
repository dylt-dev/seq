# `seq`

### Making Go Iterators Fun
```
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dylt-dev/seq"
)

func main () {
	rd, err := os.Open("./petnames.txt")
	if err != nil { panic("Unable to open petnames file")}
	petNamesSeq := seq.NewLineSeq(rd)
	// load all the names into an array
	names := []string{}
	for name := range seq.Iter1(petNamesSeq) {
		// Sequences let you check if the previous read actually resulted in an error, like EOF
		if petNamesSeq.Err() == nil {
			fmt.Printf("name=%s\n", name)
			names = append(names, name)
		} else {
			// EOFs are ok. Other errors are terrifying.
			if errors.Is(petNamesSeq.Err(), io.EOF) {
				fmt.Println("Normal EOF reached. All is well.")
			} else {
				fmt.Printf("%v\n", err)
				panic("Unepxected error!")
			}
		}
	}
	fmt.Printf("%d name(s) in the array\n", len(names))
}
```