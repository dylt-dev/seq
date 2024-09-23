# `seq`

### Making Go Iterators Go
```
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dylt-dev/seq"
)

/*
Demo app that demonstrates the following

- Create a seq.LineSeq to get data from an io.Reader line by line
- Use seq.Where() to filter the results
- Use seq.Skip() to skip N matches
- Use seq.Limit() to limit the number of results
- Use seq.Iter() to create an iterator and for..range on the results
- Use seq.Seq[T].Next() manually to confirm that after for..range the iterator is as expected

First 10 names from petnames.txt
--------------
0. AJ
1. Abbey
2. Abbie
3. Abel
4. Abigail
5. Ace
6. Adam
7. Admiral
8. Aires
9. Ajax
*/

func main () {

	// Open file of cute pet names
	var rd io.Reader
	var err error
	rd, err = os.Open("./petnames.txt")
	if err != nil { panic("Unable to open petnames file")}
	// Create a LineSeq from the file -- this is a Sequence that returns the contents of a reader \n by \n
	var sq seq.Seq[string] = seq.NewLineSeq(rd)
	var filter seq.FilterFunc[string] = func (name string) bool { return strings.HasPrefix(name, "Ab")}
	sq = seq.Where(sq, filter)		// Only names starting with 'Ab'
	sq = seq.Skip(sq, 3)			// Skip the first 3 matches
	sq = seq.Limit(sq, 1)			// Only get 1 match
	// Print the one match
	var name string
	for name = range seq.Iter(sq) {
		fmt.Printf("name=%s\n", name)
	}
	// Confirm that the end state of the sequence is as expected
	name, err = sq.Next()
	fmt.Printf("After loop: name=%s err=%s\n", name, err.Error())
}
```