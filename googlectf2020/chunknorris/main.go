package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)



type factordb struct {
	ID      string          `json:"id"`
	Status  string          `json:"status"`
	Factors [][]interface{} `json:"factors"`
}

type RSA struct {P,Q,E,N,C,D big.Int}
func NewRSA(r io.Reader)RSA{
	var (
		n,e,c big.Int
	)
	begin := "0x"
	hexNumber := make([][]byte,0)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
		for s.Scan() {
		text := s.Text()
		index := strings.Index(text,begin)
		b, err := hex.DecodeString(strings.Trim(text[index+2:],"\n"))
		if err != nil{
			fmt.Println(len(text[index+2:]))
			log.Fatal(err)
		}
		hexNumber = append(hexNumber,b)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	e.SetInt64(65537)
	n.SetBytes(hexNumber[0])
	c.SetBytes(hexNumber[1])


	req, err := http.Get(fmt.Sprintf("http://factordb.com/api?query=%s", n.String()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err create request %s\n", err)
	}
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var fdb factordb
	err = decoder.Decode(&fdb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err unmarshal response %s\n", err)
	}
	p, _ := new(big.Int).SetString(fdb.Factors[0][0].(string),10)
	q, _ := new(big.Int).SetString(fdb.Factors[1][0].(string),10)
	fi := (p.Sub(p,big.NewInt(1))).Mul(p,(q.Sub(q,big.NewInt(1))))
	d := new(big.Int).ModInverse(&e,fi)
	return RSA{
		P: *p,Q: *q,E: e,N: n,C: c,D: *d,
	}
}
func (r *RSA)decrypt(){
	const size = 64
	var bb big.Int
	plain := new(big.Int).Exp(&r.C,&r.D,&r.N) // decrypt message
		var db = make([]byte,size)
	dx := size
	buff := big.NewInt(0xff)

	for plain.BitLen() > 0{
		dx--
		b := bb.And(plain,buff).Int64()
		db[dx] = byte(b)
		plain.Rsh(plain,8)
	}
	fmt.Println("Decoded number as text:", string(db[dx:]))
}





func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rsa := NewRSA(f)
	rsa.decrypt()
}
