package main

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
	"unicode"
)

func main() {
	n := `456378902858290907415273676326459758501863587455889046415299414290812776158851091008643992243505529957417209835882169153356466939122622249355759661863573516345589069208441886191855002128064647429111920432377907516007825359999`
	e := `65537`
	c := `41662410494900335978865720133929900027297481493143223026704112339997247425350599249812554512606167456298217619549359408254657263874918458518753744624966096201608819511858664268685529336163181156329400702800322067190861310616`
	N, E, C, _ := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	N, _ = N.SetString(n, 10)
	E, _ = E.SetString(e, 10)
	C, _ = C.SetString(c, 10)
	S := new(big.Int)
	S, _ = S.SetString("2",10)
	C2 := S.Exp(S,E,N)
	sendMessage := C.Mul(C,C2)

	conn, err := net.Dial("tcp","challenge01.root-me.org:51031")
	if err != nil{
		log.Fatal(err)
	}
	var bf = make([]byte,1024)
	size, err := conn.Read(bf)
	if err != nil{
		log.Fatal(err)
	}
	conn.Write([]byte(fmt.Sprintf("%s\r\n",sendMessage.String())))
	size, err = conn.Read(bf)
	if err != nil{
		log.Fatal(err)

	}
	text := strings.Split(string(bf[:size]),":")
	number := text[1]
	num := make([]rune,0)
	for _, r := range []rune(number){
		if unicode.IsDigit(r){
			num = append(num,r)
		}
	}
	mp := new(big.Int)
	mp, _ = mp.SetString(string(num),10)
	res := mp.Div(mp,big.NewInt(2))
	fmt.Println(string(res.Bytes()))
}


