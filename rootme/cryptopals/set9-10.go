package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func pkcs7(b []byte,size int) []byte {
	if size > 255 {
		panic("invalid lenght size, max 255")
	}
	padding := make([]byte, size, size)
	leght := len(b)
	if leght < size {
		copy(padding, b)
		countBytes := size - leght%size
		for i := 0; i < countBytes; i++ {
			padding[leght + i] = byte(countBytes)
		}
	}
	return padding
}

func EncryptCBC(src []byte,b cipher.Block,IV []byte) []byte{
	out := make([]byte,len(src))
	var xor = func(a,b []byte)(xorBytes []byte){
		if len(a) > len(b){
			a = a[:len(b)]
		}
		xorBytes = make([]byte,len(a))
		for i := range a{
			xorBytes[i] = a[i] ^ b[i]
		}
		return
	}
	iv := IV
	for i:=0; i < len(src)/b.BlockSize(); i++{
		copy(out[i*b.BlockSize():],
			xor(src[i*b.BlockSize():(i+1)*b.BlockSize()],iv))
		b.Encrypt(out[i*b.BlockSize():],out[i*b.BlockSize():])
		iv = out[i*b.BlockSize():(i+1)*b.BlockSize()]
	}
	return out
}

func DecryptCBC(src []byte,b cipher.Block,IV []byte) []byte{
	out := make([]byte,len(src))
	temp := make([]byte,b.BlockSize())
	iv := IV
	var xor = func(a,b []byte)(xorBytes []byte){
		if len(a) > len(b){
			a = a[:len(b)]
		}
		xorBytes = make([]byte,len(a))
		for i := range a{
			xorBytes[i] = a[i] ^ b[i]
		}
		return
	}
	for i:=0; i < len(src)/b.BlockSize();i++{
		b.Decrypt(temp,src[i*b.BlockSize():])
		copy(out[i*b.BlockSize():],xor(temp,iv))
		iv = src[i*b.BlockSize():(i+1)*b.BlockSize()]
	}
	return out
}




func readFile(r io.Reader)string{
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var text strings.Builder
	for scanner.Scan(){
		text.WriteString(scanner.Text())
	}
	return text.String()
}



func main() {

	text := []byte("TELEPORT")
	key := []byte("YELLOW SUBMARINE")
	text = pkcs7(text,32)
	iv := make([]byte,16)
	b, _ := aes.NewCipher(key)

	encryptedText := EncryptCBC(text,b,iv)
	fmt.Printf("Encrypted text: %s\n",encryptedText)
	decryptedText := DecryptCBC(encryptedText,b,iv)
	fmt.Printf("Dencrypted text: %s\n",decryptedText)
	f, _ := os.Open("10.txt")
	fileText := readFile(f)

	candidate, _ := base64.StdEncoding.DecodeString(fileText)

	result := DecryptCBC(candidate,b,iv)

	fmt.Println(string(result))


}

