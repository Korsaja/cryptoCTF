package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"io"
	"math"
	"math/bits"
	"os"
	"strings"
	"unicode/utf8"
)

// set 1 block 5
func repeatingKeyXOR(key,plaintext []byte)[]byte{
	ciphertext := make([]byte,len(plaintext))
	for i := 0; i < len(plaintext); i++{
		shift := key[i%len(key)]
		c := plaintext[i] ^ shift
		ciphertext[i] = c
	}
	return ciphertext
}

var freqSymbolMap = make(map[rune]float64)




func hammingDistance(a,b []byte)int{
	if len(a) != len(b){
		return -1
	}
	var distance int
	for i := range a{
		distance += bits.OnesCount8(a[i] ^ b[i])
	}
	return distance
}

func loadCipherText(r io.Reader)string{
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var result strings.Builder
	for scanner.Scan(){
		result.WriteString(scanner.Text())
	}
	return result.String()
}


func searchXORKeySize(in []byte, f func(a,b []byte)(int))int{
	var res int
	bestScore := math.MaxFloat64
	for keySize := 2 ; keySize <= 40; keySize++{
		a,b := in[:keySize*4],in[keySize*4:keySize*4*2]
		score := float64(f(a,b)) / float64(keySize) //normalization
		//fmt.Fprintf(os.Stderr,"\rscore: %0.2f key size %d",score,keySize)
		//time.Sleep(200*time.Millisecond)
		if score < bestScore{
			res = keySize
			bestScore = score
		}
	}
	fmt.Println()
	return res
}


func base64ToBinary(text string)[]byte{
	binaryText, err := base64.StdEncoding.DecodeString(text)
	if err != nil{
		panic("err decode text to binary")
	}
	return binaryText
}

func freqSymbol(r io.Reader)map[rune]float64{
	freq := make(map[rune]float64)
	var sb strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan(){sb.WriteString(scanner.Text())}
	err := scanner.Err()
	if err != nil{
		log.Fatal(err)
	}
	totalCount := utf8.RuneCountInString(sb.String())

	for _, r := range sb.String(){
		freq[r]++
	}
	for r, _ := range freq{
		freq[r] = (freq[r] * 100)/float64(totalCount)
	}
	return freq
}

func singleXOR(message []byte,key byte)[]byte{
	res := make([]byte,len(message))
	for i,c := range message{
		res[i] = c ^ byte(key)
	}
	return res
}

func englishScore(b []byte)(result float64){
	for _, j := range b{
		if value, ok := freqSymbolMap[rune(j)];ok{
			result += value
		}
	}
	return
}


func findSingleXOR(in []byte)(res []byte,key byte,score float64){
	for i:=0;i < 256;i++{
		out := singleXOR(in,byte(i))
		s := englishScore(out)
		if s > score{
			res = out
			score = s
			key = byte(i)
		}
	}
	return
}



func exploit(ciphertext []byte,
	hammingFunc func(a,b []byte)int)string{
	keySize := searchXORKeySize(ciphertext,hammingFunc)
	columns := make([]byte,(len(ciphertext)+keySize-1)/keySize)
	key := make([]byte,keySize)
	for i:=0;i < keySize;i++{
		for row := range columns{
			if row*keySize+i >= len(ciphertext){
				continue
			}
			columns[row] = ciphertext[row*keySize+i]
		}
		_ ,k ,_ := findSingleXOR(columns)
		key[i] = k
	}
	return string(key)
}


func main(){
	file , err := os.Open("set1_challenge6")
	if err != nil{
		log.Fatal(err)
	}
	ffile, err := os.Open("alice.txt")
	if err != nil{
		log.Fatal(err)
	}

	freqSymbolMap = freqSymbol(ffile)
	base64text := loadCipherText(file)
	ciphertext := base64ToBinary(base64text)
	key := exploit(ciphertext,hammingDistance)
	fmt.Println(key)

}
