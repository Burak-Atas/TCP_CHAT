package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Server başlatılırken hata oluştu: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Sunucu başlatılıyor :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Bağlantı kabul edilemedi: %s", err.Error())
			continue
		}

		c := s.newClient(conn)
		go c.readInput()

	}
}
