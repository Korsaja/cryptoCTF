package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/bits"
	"os"
	"unicode/utf8"
)

func main() {

	img, err := ioutil.ReadFile(`ch3.bmp`)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(`alice.txt`)
	freq := freqSymbol(f)
	keySize := searchKeySize(img) // lenght key

	var singleXor = func(m []byte, k byte) []byte {
		res := make([]byte, len(m))
		for i := 0; i < len(m); i++ {
			res[i] = m[i] ^ k
		}
		return res
	}
	var score = func(m []byte) float64 {
		var res float64
		for _, r := range m {
			if val, ok := freq[rune(r)]; ok {
				res += val
			}
		}
		return res
	}
	var findSingleXor = func(m []byte)(key byte){
		var scores float64
		for i := 0; i < 256; i++ {
			out := singleXor(m,byte(i))
			s := score(out)
			if s > scores{
				scores = s
				key = byte(i)
			}
		}

		return
	}

	columns := make([]byte, (len(img)+keySize-1)/keySize)
	key := make([]byte,keySize)
	for i := 0; i < keySize; i++ {
		for idx := range columns {
			if idx*keySize+i >= len(img) {
				continue
			}
			columns[idx] = img[idx*keySize+i]
		}
		k := findSingleXor(columns)
		key[i] = k

	}
	i, err := decrypt(img,key)
	if err != nil{
		fmt.Println(err)
	}

	f, err = os.Create("res.bmp")
	if err != nil{
		log.Fatal(err)
	}
	bmp.Encode(f,i)

}


func decrypt(in,key []byte)(image.Image,error){
	key = bytes.ToLower(key)
	res := make([]byte,len(in))
	for i:=0;i < len(in);i++{
		shift := key[i%len(key)]
		res[i] = in[i] ^ shift
	}
	img, err  := bmp.Decode((bytes.NewReader(res)))
	if err != nil{
		return nil,err
	}
	return img,nil

}

func freqSymbol(r io.Reader) map[rune]float64 {
	var text string
	fs := make(map[rune]float64)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		text += s.Text()
	}
	for _, r := range text {
		fs[r]++
	}
	for r, _ := range fs {
		fs[r] = (fs[r] * 100) / float64(utf8.RuneCountInString(text))
	}
	return fs
}

func searchKeySize(in []byte) int {
	var score float64
	var bestScore = math.MaxFloat64
	var keySize int
	for i := 2; i <= 40; i++ {
		a, b := in[:i*4], in[i*4:i*4*2]
		score = float64(hammingDistance(a, b)) / float64(i)
		if score < bestScore {
			keySize = i
			bestScore = score
		}
	}
	return keySize
}

func hammingDistance(a, b []byte) int {
	var distance int
	for i := range a {
		distance += bits.OnesCount8(a[i] ^ b[i])
	}
	return distance
}
