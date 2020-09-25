package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

var(
	hexONE = "1c0111001f010100061a024b53535009181c"
	hexTWO = "686974207468652062756c6c277320657965"
)




func main(){
	var thePlainMessage string
	_HEXONE_, _ := hex.DecodeString(hexONE)
	_HEXTWO_, _ := hex.DecodeString(hexTWO)

	for i := range _HEXONE_{
		thePlainMessage  += string(_HEXONE_[i] ^ _HEXTWO_[i])
		time.Sleep(500*time.Millisecond)
		fmt.Println(thePlainMessage)
	}

}
