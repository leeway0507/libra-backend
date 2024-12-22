package utils

func RemoveEmptyStringInSlice(s []string) []string {
	var newS []string
	for _, a := range s {
		if a != "" {
			newS = append(newS, a)
		}
	}
	return newS
}
