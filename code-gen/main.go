package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type YamlToStruct struct {
	Trigger       string `yaml:"trigger,omitempty" json:"trigger,omitempty"`
	Service       string `yaml:"service,omitempty" json:"service,omitempty"`
	TriggerValues struct {
		Method string `yaml:"method,omitempty" json:"method,omitempty"`
		Path   string `yaml:"path,omitempty" json:"path,omitempty"`
	} `yaml:"trigger_values,omitempty" json:"trigger_values,omitempty"`
}

type CapabilitiesStruct struct {
	ContractID string `yaml:"contract_id,omitempty" json:"contract_id,omitempty"`
}

func createFilePath(data map[string]string) string {
	path_ := filepath.Dir(data["packageName"])
	filePath := path_ + "/" + data["packageName"] + "/"

	return filePath
}

func updateToken(data map[string]string, fileList []string, nameList []string, idx int, conventionList []string) bool {
	for _, v := range fileList {
		readFile := ""
		if idx == 0 {
			readFile = "./files/" + v + ".txt"
		} else {
			filePath_ := createFilePath(data)
			readFile = filePath_ + v + ".go"
			if v == "package" {
				readFile = filePath_ + data["packageName"] + ".go"
			}

		}

		input, errReadFile := ioutil.ReadFile(readFile)
		if errReadFile != nil {
			fmt.Println(errReadFile)
			log.Fatal(errReadFile)
		}

		filePath := createFilePath(data)

		fileName := filePath + v + ".go"
		if v == "package" {
			fileName = filePath + data["packageName"] + ".go"
		}

		replaceText_ := bytes.Replace(input, []byte(conventionList[idx]), []byte(nameList[idx]), -1)

		errFile_ := ioutil.WriteFile(fileName, replaceText_, 0666)
		if errFile_ != nil {
			fmt.Println(errFile_)
			os.Exit(1)
		}
	}

	log.Println("File Creation Successful!!")
	return true

}

func YamlToStructModifier(
	trigger string,
	service string,
	method string,
	path string) (*yaml.Node, error) {

	app := YamlToStruct{
		Trigger: trigger,
		Service: service,
		TriggerValues: struct {
			Method string `yaml:"method,omitempty" json:"method,omitempty"`
			Path   string `yaml:"path,omitempty" json:"path,omitempty"`
		}{method, path},
	}
	marshalledApp, err := yaml.Marshal(&app)
	if err != nil {
		return nil, err
	}

	node := yaml.Node{}
	if err := yaml.Unmarshal(marshalledApp, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

func YamlToStructCapabilitiesModifier(contract_id string) (*yaml.Node, error) {
	app := CapabilitiesStruct{
		ContractID: contract_id,
	}
	marshalledApp, err := yaml.Marshal(&app)
	if err != nil {
		return nil, err
	}

	node := yaml.Node{}
	if err := yaml.Unmarshal(marshalledApp, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func createTriggersAndCapabilities(data map[string]string) bool {
	yamlNode := yaml.Node{}
	///// yaml file and file path should be written here(Where triggers and contract_id will be generated)
	sourceYaml, errReadFile := ioutil.ReadFile("manifest.yaml")
	if errReadFile != nil {
		fmt.Println(errReadFile)
		log.Fatal(errReadFile)
	}
	err := yaml.Unmarshal([]byte(sourceYaml), &yamlNode)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	triggerText_ := data["triggerText"]
	serviceText_ := data["contractText"]
	method_ := data["apiMethod"]
	path_ := data["path"]

	result, err := YamlToStructModifier(triggerText_, serviceText_, method_,
		path_)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	capabilitiesResult, capabilitiesErr := YamlToStructCapabilitiesModifier(data["contractText"])
	if capabilitiesErr != nil {
		log.Fatalf("error: %v", capabilitiesErr)
	}

	appIdx := -1
	for i, k := range yamlNode.Content[0].Content {
		if k.Value == "triggers" {
			appIdx = i + 1
			break
		}
	}
	appIdxCapabilities := -1
	for i, k := range yamlNode.Content[0].Content {
		if k.Value == "capabilities" {
			appIdxCapabilities = i + 1
			break
		}
	}
	yamlNode.Content[0].Content[appIdx].Content = append(
		yamlNode.Content[0].Content[appIdx].Content, result.Content[0])
	yamlNode.Content[0].Content[appIdxCapabilities].Content = append(
		yamlNode.Content[0].Content[appIdxCapabilities].Content, capabilitiesResult.Content[0])

	out, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(out))

	///// yaml file and file path should be written here(Where triggers and contract_id will be generated)
	errFile_ := ioutil.WriteFile("manifest.yaml", out, 0666)
	if errFile_ != nil {
		fmt.Println(errFile_)
		os.Exit(1)
	}

	log.Println("YAML File written successful!!")

	return true
}

func main() {
	///// yaml file and file path should be written here (Where all key data will be taken)
	input, errYaml := ioutil.ReadFile("sample.yaml")
	if errYaml != nil {
		log.Fatal(errYaml)
	}
	data := make(map[string]string)

	errUnmarshal := yaml.Unmarshal(input, &data)

	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}
	fileList := []string{"package", "category", "contractid", "model", "name"}

	conventionList := []string{"<package>", "<contractid>", "<modifiedContractId>", "<capabilityStruct>"}
	errDir := os.Mkdir(data["packageName"], 0755)
	if errDir != nil {
		log.Fatal(errDir)
	}

	packageName := data["packageName"]
	contractName := data["contractText"]
	capabilityModel := data["capabilityStruct"]
	modifiedContractText_ := strings.Replace(data["contractText"], ":", "_", 5)
	var nameList []string
	nameList = append(nameList, packageName)
	nameList = append(nameList, contractName)
	nameList = append(nameList, modifiedContractText_)
	nameList = append(nameList, capabilityModel)

	for i := 0; i < len(nameList); i++ {
		updateToken(data, fileList, nameList, i, conventionList)
	}

	createTriggersAndCapabilities(data)

}
