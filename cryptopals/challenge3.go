package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"sort"
	"sync"
	"unicode/utf8"
)

const _TARGET_ = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
const _FILE_ = "alice.txt"

type Message struct {
	key rune
	words []byte
	score float64
}

type messages []Message
func (m messages) Len() int           { return len(m) }
func (m messages) Less(i, j int) bool {
	switch fd := m[i].score - m[j].score; {
	case fd < 0:
		return false
	case fd > 0:
		return true
	}
	return m[i].key < m[j].key
}
func (m messages) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }



func XOR(message []byte,key rune)Message{
	res := make([]byte,len(message))
	for i,c := range message{
		res[i] = c ^ byte(key)
	}
	return Message{
		key:key,
		words:res,
		score: EnglishScore(res),
	}
}
func EnglishScore(b []byte)(result float64){
	freq := FreqRune()
	for _, j := range b{
		if value, ok := freq[rune(j)];ok{
			result += value
		}
	}
	return
}
func FreqRune()map[rune]float64{
	text, err := ioutil.ReadFile(_FILE_)
	if err != nil{
		fmt.Println(err)
	}
	totalCount := utf8.RuneCountInString(string(text))
	m := make(map[rune]float64)
	for _, r := range string(text){
		m[r]++
	}
	for r, _ := range m{
		m[r] = (m[r] * 100)/float64(totalCount)
	}
	return m
}



func main(){
	_HEXTARGET_, _ := hex.DecodeString(_TARGET_)

	var (
		ch = make(chan Message,len(FreqRune()))
		wg sync.WaitGroup
		done = make(chan struct{})
		scores messages
	)
	for k ,_ := range FreqRune(){
		wg.Add(1)
		go func(key rune) {
			ch <- XOR(_HEXTARGET_,key)
			wg.Done()
		}(k)
	}

	go func() {
		for {
			select {
			case m := <-ch:
				scores = append(scores,m)
			case <-done:
				return
			}
		}
	}()


	wg.Wait()
	close(ch)
	sort.Sort(scores)

	fmt.Printf("[%c] %.2f: %s\n",scores[0].key,scores[0].score,string(scores[0].words))
}