package webhp

import (
	"net/http"
	"testing"
)

func TestLoadGenerator(t *testing.T){
	lg := NewLoadGenerator(http.MethodGet,"http://127.0.0.1:8080", 10000)
	lg.Execute()
}