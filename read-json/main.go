package main

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"io/ioutil"
	"log"
	"os"
)

type StructJson struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	// Read the JSON file
	fileName := "file.json"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data StructJson
	err = json.Unmarshal(content, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("READ DATA OF FILE")
	log.Println(fmt.Sprintf("Name : %s", data.Name))
	log.Println(fmt.Sprintf("Email : %s", data.Email))
	log.Println(fmt.Sprintf("Age : %d", data.Age))

	newAge := data.Age + 1 // this should be for age
	newData := StructJson{
		Name:  faker.Name(),
		Age:   newAge,
		Email: faker.Name(),
	}

	log.Println("REWRITE DATA OF FILE")
	log.Println(fmt.Sprintf("Name from %s to %s", data.Name, newData.Name))
	log.Println(fmt.Sprintf("Email from %s to %s", data.Email, newData.Email))
	log.Println(fmt.Sprintf("Age from %d to %d", data.Age, newData.Age))

	log.Println("MODIFIED FILE JSON")

	modified, err := json.Marshal(newData)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(fileName, modified, os.ModePerm)
	if err != nil {
		fmt.Println("Failed for writing JSON file:", err)
	}

	log.Println("SUCCESS MODIFIED JSON FILE")
}
