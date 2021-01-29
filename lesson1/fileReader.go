package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	arguments := os.Args[1:]

	argsError := ValidArgs(arguments)

	if argsError != nil {
		writeToFile("error.txt", argsError.Error())
		return
	}

	for _, fileName := range arguments {
		fileError := FileExist(fileName)
		if fileError != nil {
			writeToFile("error.txt", fileError.Error())
			return
		}
		jsonError := ValidJson(fileName)
		if jsonError != nil {
			writeToFile("error.txt", jsonError.Error())
			return
		}

	}
	firstByte, _  := readJSON(arguments[0])
	secondByte, _ := readJSON(arguments[1])
	
	var json1 map[string]interface{}
	var json2 map[string]interface{}

	json.Unmarshal(firstByte, &json1)
	json.Unmarshal(secondByte, &json2)

	var add []string
	var same []string
	var change []string
	var remove []string

	for key := range json2 {
		if _, ok := json1[key]; !ok {
			add = append(add, key)
			continue
		}
		
		if json1[key] == json2[key] {
			same = append(same, key)
		}else {
			change = append(change, key)
		}
	}

	for key := range json1 {
		if _, ok := json2[key]; !ok {
			remove = append(remove, key)
		}
	}

	sort.Strings(add)
	sort.Strings(same)
	sort.Strings(change)
	sort.Strings(remove)

	writeSliceToFile("add.txt", add)
	writeSliceToFile("same.txt", same)
	writeSliceToFile("change.txt", change)
	writeSliceToFile("remove.txt", remove)
}

func ValidArgs(args []string) error {
	if len(args) != 2 {
		return errors.New("args: enter 2 filenames")
	}

	for _, str := range args {
		if len(str) < 1 {
			return errors.New("args: empty arguments")
		}
	}

	return nil
}

func FileExist(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func ValidJson(fileName string) error {
	byteValue, openError := readJSON(fileName)
	if openError != nil {
		return openError
	}
	var result map[string]interface{}
	unmarshalError := json.Unmarshal(byteValue, &result)

	if unmarshalError != nil {
		return unmarshalError
	}

	return nil
}

func writeToFile(fileName string, textToWrite string) {
	file, fileError := os.OpenFile("./"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if fileError != nil {
		fmt.Println("Something goes wrong")
	}

	defer file.Close()                         

	if _, err := file.Write([]byte(textToWrite + "\n")); err != nil {
		fmt.Println("Cannot write to file")
	}
}

func readJSON(fileName string) ([]byte, error) {
	jsonFile, openError := os.Open(fileName)

	if openError != nil {
		return nil, openError
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue, nil
}

func writeSliceToFile(fileName string, sliceToWrite []string){
	for _, item := range sliceToWrite {
		writeToFile(fileName, item)
	}
}
