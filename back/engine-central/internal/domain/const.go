package domain

import "strings"

var countriesNormalizations = map[int][]string{
	1: {"colombia", "co"},
	2: {"mexico", "mx"},
	3: {"chile", "cl"},
	4: {"estados unidos", "us"},
	5: {"costa rica", "cri"},
}

// normalizeCountry returns the country id from the country name
func NormalizeCountry(countryName string) int {
	countryName = strings.ToLower(countryName)
	countryName = strings.ReplaceAll(countryName, " ", "")
	countryName = strings.ReplaceAll(countryName, "-", "")
	countryName = strings.ReplaceAll(countryName, "_", "")

	for id, names := range countriesNormalizations {
		for _, name := range names {
			if name == countryName {
				return id
			}
		}
	}

	return 1
}
