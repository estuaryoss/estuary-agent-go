package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func DoesFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func IsFolder(folderName string) bool {
	info, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func AppendFile(fileName string, content string) {
	CreateFileIfNotExist(fileName)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed opening file: %s", fileName)
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(fmt.Sprintf("Failed appending to file: %s", fileName))
	}
	file.Close()
}

func ReadFile(fileName string) []byte {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print(fmt.Sprintf("Failed reading content from file: %s", fileName))
	}
	return content
}

func WriteFile(fileName string, content []byte) {
	CreateFileIfNotExist(fileName)
	err := ioutil.WriteFile(fileName, content, 0644)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed writing content to file: %s", fileName))
	}
}

func WriteFileJson(fileName string, content interface{}) {
	CreateFileIfNotExist(fileName)
	jsonContent, err := json.Marshal(content)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed writing JSON content to file: %s", fileName))
	}
	err = ioutil.WriteFile(fileName, jsonContent, 0644)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed writing content to file: %s", fileName))
	}
}

func CreateFileIfNotExist(fileName string) {
	if DoesFileExists(fileName) {
		return
	}
	emptyFile, err := os.Create(fileName)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create empty file: %s", fileName))
	}
	emptyFile.Close()
}

func CreateDir(dirName string) {
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		return
	}
	err := os.Mkdir(dirName, 0644)
	if err != nil {
		log.Print(fmt.Sprintf("Failed creating dir: %s", dirName))
	}
}

func RecreateFile(fileName string) {
	DeleteFile(fileName)
	CreateFileIfNotExist(fileName)
}

func DeleteFile(fileName string) {
	if !DoesFileExists(fileName) {
		return
	}
	err := os.Remove(fileName)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to delete file: %s", fileName))
	}
}

func DeleteFiles(fileNames []string) {
	for _, fileName := range fileNames {
		err := os.Remove(fileName)
		if err != nil {
			log.Printf(fmt.Sprintf("Failed to delete file: %s", fileName))
		}
	}
}

func RecreateFiles(fileNames []string) {
	for _, fileName := range fileNames {
		RecreateFile(fileName)
	}
}

func OpenFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf(err.Error())
	}
	return file
}

func CloseFile(file *os.File) error {
	return file.Close()
}
