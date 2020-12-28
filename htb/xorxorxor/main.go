package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

var target = "134af6e1297bc4a96f6a87fe046684e8047084ee046d84c5282dd7ef292dc9"

func xor(b []byte,key []byte)[]byte{
	out :=  make([]byte,len(b))
	for i := range b{
		out[i] = b[i] ^ key[i % len(key)]
	}
	return out
}


func main(){

	flag_start := "HTB{"
	h, err := hex.DecodeString(target)
	if err != nil{
		log.Fatal(err)
	}

	key := xor(h[:4],[]byte(flag_start))

	fmt.Println(string(xor(h,key)))

}
