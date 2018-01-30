package network

import (
	"fmt"
	"io"
	"net"
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

func Download(node Node, req Request, shared_dir string) {
	connection := node.Conn

	req.SendOnExisting(connection)
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	connection.Read(bufferFileName)
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
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
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

func SendFileToClient(connection net.Conn, req Request, pathToFile string) {
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
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	fS := strings.Replace(fileSize, ":", "", -1)
	fN := strings.Replace(fileName, ":", "", -1)
	fmt.Printf("\nSending file `%s` (%s Bytes)", fN, fS)

	gofuncs := int64(16)

	for num := int64(0); num < gofuncs; num++ {
		go func(conn net.Conn) {
			sendBuffer := make([]byte, 4096)
			for {
				_, err = file.Read(sendBuffer)
				if err == io.EOF {
					break
				}
				connection.Write(sendBuffer)
			}
		}(connection)

	}
	// fmt.Printf("\nSuccessfully sent `%s`!", fN)
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
