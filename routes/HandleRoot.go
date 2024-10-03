package routes

import (
	"fmt"
	"net/http"
)

// nothing to see here... hehehe
func HandleRoot(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hello, Dunia! :D")
}
