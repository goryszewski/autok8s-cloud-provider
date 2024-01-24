package autok8s

import (
	"net/http"
)

// newautok8sClient returns a specific HTTP client used when communicating with the autok8s API(s)
func newAutok8sClient() *http.Client {
	return &http.Client{}
}
