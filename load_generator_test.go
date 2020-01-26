package webhp

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLoadGenerator(t *testing.T) {
	lg := CreatePLoadGenerator()
	lg.Execute()
	fmt.Println("===========================")
}

func CreatePLoadGenerator() *LoadGenerator {
	lg := NewLoadGenerator(http.MethodPost, "http://127.0.0.1:8080", 10000, 10)
	return &lg
}
