package main

import (
	"fmt"
	"log"

	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/mac"
)

func main() {

	kh, err := keyset.NewHandle(mac.HMACSHA256Tag256KeyTemplate())
	if err != nil {
		log.Fatal(err)
	}

	m, err := mac.New(kh)
	if err != nil {
		log.Fatal(err)
	}

	mac, err := m.ComputeMAC([]byte("this data needs to be MACed"))
	if err != nil {
		log.Fatal(err)
	}

	if m.VerifyMAC(mac, []byte("this data needs to be MACed")); err != nil {
		log.Fatal("MAC verification failed")
	}

	fmt.Println("MAC verification succeeded.")

}
