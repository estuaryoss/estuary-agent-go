package utils

import "strings"

func TrimSpacesAndLineEndings(data []string) []string {
	var compiledData []string
	for _, elem := range data {
		trimElem := strings.TrimSuffix(elem, "\r")
		trimElem = strings.TrimSpace(trimElem)
		if trimElem != "" {
			compiledData = append(compiledData, trimElem)
		}
	}

	return compiledData
}
