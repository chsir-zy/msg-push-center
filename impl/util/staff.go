package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type staff struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ReadStaff() map[string]string {
	log := path.Join("config", "staff.log")
	content, err := os.ReadFile(log)
	if err != nil {
		fmt.Println(err)
	}

	var s []staff
	err = json.Unmarshal([]byte(content), &s)
	if err != nil {
		fmt.Println(err)
	}

	staffMap := make(map[string]string)
	for _, v := range s {
		staffMap[v.Id] = v.Name
	}

	return staffMap
}
