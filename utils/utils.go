// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"net"
	"net/http"
	"os"

	"github.com/cetessai2501/gasilo/model"
)

func StringArrayIntersection(arr1, arr2 []string) []string {
	arrMap := map[string]bool{}
	result := []string{}

	for _, value := range arr1 {
		arrMap[value] = true
	}

	for _, value := range arr2 {
		if arrMap[value] {
			result = append(result, value)
		}
	}

	return result
}

func FileExistsInConfigFolder(filename string) bool {
	if len(filename) == 0 {
		return false
	}

	if _, err := os.Stat(FindConfigFile(filename)); err == nil {
		return true
	}
	return false
}

func RemoveDuplicatesFromStringArray(arr []string) []string {
	result := make([]string, 0, len(arr))
	seen := make(map[string]bool)

	for _, item := range arr {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

func GetIpAddress(r *http.Request) string {
	address := r.Header.Get(model.HEADER_FORWARDED)

	if len(address) == 0 {
		address = r.Header.Get(model.HEADER_REAL_IP)
	}

	if len(address) == 0 {
		address, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return address
}
