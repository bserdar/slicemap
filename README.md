[![GoDoc](https://godoc.org/github.com/bserdar/slicemap?status.svg)](https://godoc.org/github.com/bserdar/slicemap)
[![Go Report Card](https://goreportcard.com/badge/github.com/bserdar/slicemap)](https://goreportcard.com/report/github.com/bserdar/slicemap)

# Slicemap

This package implements a map that uses a slice of values as key.

Go does not allow using slices as map keys. This especially becomes
cumbersome if you have a composite key with non-constant number of
elements. This package uses a nested map structure to associate values
with keys where the keys are a slice of a comparable object.

```
sm:=SliceMap[string,int]{} // Create a map[[]string]int
sm.Put([]string{"a","b","c"},1)
sm.Put([]string{"d","e"},2)
fmt.Println(sm.Get([]string{"a","b","c"})) // Prints 1, true
fmt.Println(sm.Get([]string{"d","e"}))  // Prints 2, true
fmt.Println(sm.Get([]string{"f","g"}))  // Prints 0, false
```

