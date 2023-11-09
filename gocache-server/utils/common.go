package utils

type Str_Error string

func (s Str_Error) Error() string {
	return string(s)
}
func VerifyPermit(passwd, db string) bool {

	return false
}
func TransMap[T any](src map[string]T) map[string]any {
	var ans map[string]any = make(map[string]any)
	for k, ele := range src {
		ans[k] = ele
	}
	return ans
}
