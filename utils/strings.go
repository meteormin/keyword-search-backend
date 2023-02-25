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

type realUrlMaker struct {
	Host     string
	Path     string
	QueryKey string
	Query    string
}

func (rm *realUrlMaker) makeRealUrl() string {
	var realUrl string
	host := JoinHostPath(rm.Host, rm.Path)

	if rm.QueryKey == "" {
		realUrl = fmt.Sprintf("%s/%s", host, rm.Query)
	} else {
		realUrl = AddQueryString(host, map[string]interface{}{
			rm.QueryKey: rm.Query,
		})
	}

	return realUrl
}

func MakeRealUrl(host string, path string, queryKey string, query string) string {
	return (&realUrlMaker{Host: host, Path: path, QueryKey: queryKey, Query: query}).makeRealUrl()
}
