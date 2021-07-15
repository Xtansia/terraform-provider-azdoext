package utils

func MapStrings(vs []string, f func(string) string) []string {
	vso := make([]string, len(vs))
	for i, v := range vs {
		vso[i] = f(v)
	}
	return vso
}
