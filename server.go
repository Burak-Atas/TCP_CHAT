package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	files    File
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case NICK_NAME:
			s.Create_Nick(cmd.client, cmd.args)
		case JOIN_ROOM:
			s.Join_Room(cmd.client, cmd.args)
		case ROOMS:
			s.listRooms(cmd.client)
		case MESAGE:
			s.msg(cmd.client, cmd.args)
		case DOWNLOAD_FİLE:
			s.files.Download_file(cmd.client, cmd.args)
		case UPLOAD_FİLE:
			s.Upload_File(cmd.client, cmd.args)
		case QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {

	log.Printf("Yeni kullanıcı katıldı: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) Create_Nick(c *client, args []string) {

	if len(args) < 2 {
		c.msg("nick kullanımı doğrulanamadı. kullanım: /nick 'NAME' ")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("Yeni isim kabul edildi ;%s", c.nick))

}

func (s *server) Join_Room(c *client, args []string) {

	if len(args) < 2 {
		c.msg("oda ismi  kullanımı doğrulanamadı. kullanım: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]

	if !ok {

		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s odaya katıldı.", c.nick))

	c.msg(fmt.Sprintf("hoşgeldiniz %s", roomName))
}

func (s *server) listRooms(c *client) {

	var rooms []string

	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Aktif odalar: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {

	if len(args) < 2 {
		c.msg("mesaj gönderimi doğrulanamadı , kullanım: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) Upload_File(c *client, args []string) {

}

func (s *server) quit(c *client) {

	log.Printf("Client sohbetten ayrıldı : %s", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("Görüşmek üzeri kendine iyi günler dileriz =(")
	c.conn.Close()

}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {

		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		broadcastmsg := fmt.Sprintf("%s odadan ayrıldı", c.nick)
		oldRoom.broadcast(c, broadcastmsg)

	}
}
