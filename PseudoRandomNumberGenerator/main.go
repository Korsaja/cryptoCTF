package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)
const KEYSIZE = 32
var charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
var holdrand uint32

func Rand()uint32{
	holdrand = holdrand * 214013 + 2531011
	return ( holdrand >> 16) & 0x7fff
}

func generateKey()string{
	var str strings.Builder
	for i:=0;i < KEYSIZE;i++{
		char := charSet[Rand() % 61]
		str.WriteByte(char)
	}
	return str.String()
}

func getSign(cp []byte)string{
	magicPlusVersion := []byte{0x42,0x5a,0x68}
	var buff bytes.Buffer
	for i := range magicPlusVersion{
		 buff.WriteByte(magicPlusVersion[i]^cp[i])
	}
	return buff.String()
}

func checkBlock(b byte)[]byte{
	out := make([]byte,9)
	var j int
	for i:=0x31;i<=0x39;i++{
		out[j] = b ^ uint8(i)
		j++
	}
	return out
}

func contains(b []byte,x byte)bool{
	for i:=0;i<len(b);i++{
		if b[i] == x{
			return true
		}
	}
	return false
}


func xor(src,key []byte)[]byte{
	out := make([]byte,len(src))
	kL := len(key)
	for i:=0;i < len(src);i++{
		out[i] = src[i] ^ key[i%kL]
	}
	return out
}


func main() {
	// current cipher text sign
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	content , err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	cipherText := []byte{0x23, 0x17, 0x5D, 0x70, 0x5A, 0x11, 0x6D, 0x67, 0x37, 0x08}
	// see https://ru.wikipedia.org/wiki/Bzip2
	sign := getSign(cipherText)
	block := checkBlock(cipherText[3])
	bytePi := cipherText[4] ^ 0x31
	start := 1354320000 // 2012 1 dec
	end := 1356998400   // 2013 1 jan
	var key string

	for i := start; i < end; i++ {
		holdrand = uint32(i)
		key = generateKey()
		if strings.Contains(key[:3], sign) &&
			contains(block,key[3]) && key[4] == bytePi{
			fmt.Printf("[+] Have candidate -> %s\n",key)
			break
		}
	}
	fmt.Printf("******************************************************\n")
	plain := xor(content,[]byte(key))
	cr := bzip2.NewReader(bytes.NewReader(plain))
	scan := bufio.NewScanner(cr)
	scan.Split(bufio.ScanLines)
	for scan.Scan(){
		ll := scan.Text()
		fmt.Printf("%s\n", ll)
	}
	if err := scan.Err(); err != nil{
		log.Fatal(err)
	}

	fmt.Printf("******************************************************\n[+] Done.\n")

}









