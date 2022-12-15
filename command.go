package main

type commandID int

//
const (
	NICK_NAME commandID = iota
	JOIN_ROOM
	ROOMS
	MESAGE
	DOWNLOAD_FİLE
	UPLOAD_FİLE
	QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
