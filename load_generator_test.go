package webhp

import (
	"net/http"
	"testing"
	"time"
)

func TestLoadGenerator(t *testing.T){
	lg := NewLoadGenerator(http.MethodGet,"http://127.0.0.1:8080", 10000)
	go func (){
		<- time.Tick(3 * time.Second)
		lg.Done()
	}()

	lg.Execute()
}