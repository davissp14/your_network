package network

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SharedFileList struct {
	Files []string `json:"shared_files"`
}

func ListFiles(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f := strings.Split(path, "/")[1:]
		p := strings.Join(f, "/")
		files = append(files, p)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

const BUFFERSIZE = 1024

func Download(conn Connection, req Request, shared_dir string) {
	conn.Send(req)
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	conn.Conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	conn.Conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	newFile, err := os.Create(fmt.Sprintf("%s/%s", shared_dir, fileName))
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	var percentage float32
	percentage = 0
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn.Conn, (fileSize - receivedBytes))
			conn.Conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, conn.Conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
		fmt.Println(receivedBytes)
		if (float32(receivedBytes) / float32(fileSize) * 100) > float32(percentage+5) {
			// fmt.Printf("\nReceived %d bytes of %d", receivedBytes, fileSize)
			percentage = (float32(receivedBytes) / float32(fileSize)) * float32(100)
			fmt.Printf("\nPercentage complete: %d", int32(percentage))
		}
	}
	fmt.Printf("\nFile `%s` stored at `%s`!", fileName, shared_dir)
}

func SendFileToClient(conn Connection, pathToFile string) {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		fmt.Printf("\nFilename doesn't exist at path: %s", pathToFile)
		return
	}
	file, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	conn.Conn.Write([]byte(fileSize))
	conn.Conn.Write([]byte(fileName))
	sendBuffer := make([]byte, 4096)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Conn.Write(sendBuffer)
	}
	return
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
