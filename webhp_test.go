package webhp

import (
	"testing"
)

func TestLoadGenerator(t *testing.T){
	lg := NewLoadGenerator("http://127.0.0.1:8080", 10000)
	lg.Test()
}