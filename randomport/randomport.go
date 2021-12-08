package randomport

import (
	"fmt"
	"log"
	"net"
)

// Get random port to listen to
func Get() (port string) {
	for {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				log.Printf("%v", err)
			}
		}()
		port = fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
		break
	}
	return
}
