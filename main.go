package main

import (
	"fmt"
	"github.com/google/gops/signal"
	"io/ioutil"
	"net"
)

func main() {
	sigs := map[string]byte{
		"Stats":    signal.Stats,
		"Version":  signal.Version,
		"MemStats": signal.MemStats,
	}
	for label, s := range sigs {
		c, err := net.Dial("tcp", `127.0.0.1:57223`)
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := []byte{s}
		_, err = c.Write(buf)
		if err == nil {
			m, err := ioutil.ReadAll(c)
			if err == nil {
				if len(m) > 0 {
					fmt.Printf("%s :: %d bytes: %s\n", label, len(m), m)
				}
			}
		}
		c.Close()
	}
}
