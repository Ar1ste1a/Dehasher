package export

import (
	"Dehash/internal/files"
	"Dehash/internal/sqlite"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func WriteCredsToFile(creds []sqlite.Creds, outputFile string, fileType files.FileType) error {
	var data []byte
	var err error

	switch fileType {
	case files.JSON:
		data, err = json.MarshalIndent(creds, "", "  ")
	case files.XML:
		data, err = xml.MarshalIndent(creds, "", "  ")
	case files.YAML:
		data, err = yaml.Marshal(creds)
	case files.TEXT:
		var outStrings []string
		for _, c := range creds {
			outStrings = append(outStrings, c.ToString()+"\n")
		}
		data = []byte(strings.Join(outStrings, ""))
	default:
		return errors.New("unsupported file type")
	}

	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s.%s", outputFile, fileType.String())
	return os.WriteFile(filePath, data, 0644)
}

func WriteToFile(results sqlite.DehashedResults, outputFile string, fileType files.FileType) error {
	var data []byte
	var err error

	result := results.Results

	switch fileType {
	case files.JSON:
		data, err = json.MarshalIndent(result, "", "  ")
	case files.XML:
		data, err = xml.MarshalIndent(result, "", "  ")
	case files.YAML:
		data, err = yaml.Marshal(result)
	case files.TEXT:
		var outStrings []string
		for _, r := range result {
			out := fmt.Sprintf(
				"Id: %s\nEmail: %s\nIpAddress: %s\nUsername: %s\nPassword: %s\nHashedPassword: %s\nHashType: %s\nName: %s\nVin: %s\nAddress: %s\nPhone: %s\nDatabaseName: %s\n\n",
				r.Id, r.Email, r.IpAddress, r.Username, r.Password, r.HashedPassword, r.HashType, r.Name, r.Vin, r.Address, r.Phone, r.DatabaseName)
			outStrings = append(outStrings, out)
		}
		data = []byte(strings.Join(outStrings, ""))
	default:
		return errors.New("unsupported file type")
	}

	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s.%s", outputFile, fileType)
	return ioutil.WriteFile(filePath, data, 0644)
}
