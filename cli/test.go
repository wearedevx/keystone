package main

import 	(
	"fmt"
	 "github.com/cossacklabs/themis/gothemis/keys"
	)

func main() {
	pair, err := keys.New(keys.TypeEC)
	fmt.Printf("err: %+v\n", err)
	fmt.Printf("pair: %+v\n", pair)

}
