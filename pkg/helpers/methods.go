package helpers

const (
	get  = "GET"
	post = "POST"
	put  = "PUT"
	del  = "DELETE"
)

// IsValidMethod check if provided method is allowed
// it's used for update/create requests for /permissions endpoint
func IsValidMethod(m string) bool {
	var methods = []string{get, post, put, del}

	for i := range methods {
		if m == methods[i] {
			return true
		}
	}

	return false
}
