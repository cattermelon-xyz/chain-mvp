package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func getDataFolder() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	dataFolder := filepath.Join(homedir, "htg")
	err = os.MkdirAll(dataFolder, 0755)
	if err != nil {
		log.Println("err in ", dataFolder)
		log.Println(err)
	}
	filepath.Join(homedir+"/htg", "tmp")
	err = os.MkdirAll(dataFolder, 0755)
	filepath.Join(homedir+"/htg", "initative")
	err = os.MkdirAll(dataFolder, 0755)
	return dataFolder, err
}

func ReadFile(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return nil
	}
	return content
}

func Save(content []byte, id string, datatype string) bool {
	dataFolder, _ := getDataFolder()
	var err error
	if datatype == "node" {
		err = ioutil.WriteFile(dataFolder+"/tmp/"+id+".json", content, 0755)
	} else if datatype == "initiative" {
		err = ioutil.WriteFile(dataFolder+"/initative/"+id+".json", content, 0755)
	}
	if err != nil {
		log.Println("err in Write file")
		log.Println(err)
		return false
	}
	return true
}
