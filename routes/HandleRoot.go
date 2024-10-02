package routes

import (
	"fmt"
	"net/http"
)

func HandleRoot(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hello, Dunia!")
}
