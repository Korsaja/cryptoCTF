package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

type RSA struct {
	e, c, n *big.Int
	p, q, d    *big.Int
	request func()
	factordb
}
func Constructor(reader io.Reader) *RSA {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	r := &RSA{}

	for scanner.Scan() {
		str := strings.Split(scanner.Text(), ":")
		switch str[0] {
		case "e":
			bigInt, ok := new(big.Int).SetString(str[1], 10)
			if ok {
				r.e = bigInt
			}
		case "c":
			bigInt, ok := new(big.Int).SetString(str[1], 10)
			if ok {
				r.c = bigInt
			}
		case "n":
			bigInt, ok := new(big.Int).SetString(str[1], 10)
			if ok {
				r.n = bigInt
			}
		}
	}
	r.request = func() {
		req, err := http.Get(fmt.Sprintf("http://factordb.com/api?query=%s", r.n.String()))
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
		r.factordb = fdb
	}
	r.request()
	p, _ := new(big.Int).SetString(r.Factors[0][0].(string),10)
	q, _ := new(big.Int).SetString(r.Factors[1][0].(string),10)
	r.p = p;r.q = q
	fi := (r.p.Sub(r.p,big.NewInt(1))).Mul(r.p,(r.q.Sub(r.q,big.NewInt(1))))
	r.d = new(big.Int).ModInverse(r.e,fi)
	return r
}

func (r *RSA)decrypt(){
	const size = 32
	var bb big.Int
	plain := new(big.Int).Exp(r.c,r.d,r.n) // decrypt message

	var db = make([]byte,size)
	dx := size
	buff := big.NewInt(0xff)

	for plain.BitLen() > 0{
		dx--
		db[dx] = byte(bb.And(plain,buff).Int64())
		plain.Rsh(plain,8)
	}
	fmt.Println("Decoded number as text:", string(db[dx:]))
}
//c ^d mod n
func main() {
	filename := "rsanoob (1).txt"
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-]File unread with error %s.\n", err)
	}
	// get data for task
	rsainfo := Constructor(fd)
	rsainfo.decrypt()

}
