package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const service = `docker.hackthebox.eu:32437`

type Stage int       // 1 - login 2 - pass 3 - ciphertext 4 - flag
func (s *Stage) Inc() { *s++ }

const (
	username Stage = iota + 1
	password
	cipher
	flag
)

func bitFlipping(b []byte)string{
	token := getToken(b)
	h, _ := hex.DecodeString(string(token))
	h[0] = h[0] ^ byte('x') ^ byte('a')
	return hex.EncodeToString(h)
}

func getToken(b []byte)[]byte{
	 i := bytes.Index(b,[]byte(": "))
	 temp := b[i+2:]
	 return bytes.Split(temp,[]byte("\n"))[0]
}

func main() {
	user := []byte("xdmin")
	pass := []byte("g0ld3n_b0y")

	conn, err := net.Dial("tcp", service)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buff := make([]byte, 1<<(10*1)) // one mb
	var stage Stage = 1
	loop:
	for {
		n, _ := conn.Read(buff)
		switch stage {
		case username:
			fmt.Fprintf(os.Stdout, "%s%s\n",
				string(buff[:n]), string(user))
			conn.Write(user)
			time.Sleep(time.Second)
			stage.Inc()
		case password:
			fmt.Fprintf(os.Stdout, "%s%s\n",
				string(buff[:n]), string(pass))
			conn.Write(pass)
			time.Sleep(time.Second)
			stage.Inc()
		case cipher:
			b := bitFlipping(buff[:n])
			conn.Write([]byte(b))
			time.Sleep(time.Second)
			stage.Inc()
		case flag:
			fmt.Fprintf(os.Stdout, "%s\n",
				string(buff[:n]))
			stage.Inc()
		default:
			break loop
		}
	}
	fmt.Println("done.")
}
