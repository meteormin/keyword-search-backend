package utils

import (
	"fmt"
	pathModule "path"
	"strings"
)

func JoinHostPath(host string, path string) string {
	sep := ":/"
	splitString := strings.Split(host, sep)
	return splitString[0] + sep + pathModule.Join(splitString[1], path)

}

func AddQueryString(host string, query map[string]interface{}) string {
	var queryString string
	for key, value := range query {
		if queryString == "" {
			queryString += fmt.Sprintf("?%s=%v", key, value)
		} else {
			queryString += fmt.Sprintf("&%s=%v", key, value)
		}
	}

	return host + queryString
}
