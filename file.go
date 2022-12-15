package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

type File struct {
}

func (f File) Download_File_s(c *client, chanel chan string, bytechan chan []byte, args []string) {

	file, err := os.Open(args[1])

	if err != nil {

		openfilemsg := fmt.Sprintf("Server'e kaydedilecek dosya açılırken hata oluştu;\nfilename:%v", args[1])
		log.Println(openfilemsg)
		c.msg(openfilemsg)

	}

	filesize := make([]byte, 1024)
	file.Read(filesize)

	chanel <- file.Name()
	bytechan <- filesize
}

func (f File) Download_file(c *client, args []string) {

	if len(args) < 2 {
		c.msg("dosya gönderimi doğrulanamadı. kullanım: /download file_name")
		return
	}

	filename := make(chan string)
	file_exp := make(chan []byte, 64)

	go f.Download_File_s(c, filename, file_exp, args)
	os.Mkdir(c.nick, fs.FileMode(os.O_WRONLY))

	file, _ := os.OpenFile(c.nick+"/"+<-filename, os.O_RDWR|os.O_CREATE, 0755)
	read := strings.NewReader(string(<-file_exp))

	io.Copy(file, read)
}

func Download_file_u(c *client, chanel chan string, bytechan chan []byte) {

}
