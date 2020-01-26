package webhp

import (
	"net/http"
	"testing"
	"time"
)

func TestLoadGenerator(t *testing.T){
	lg := NewLoadGenerator(http.MethodGet,"http://127.0.0.1:8080", 755)
	go func (){
		<- time.Tick(60 * time.Second)
		lg.Done()
	}()

	lg.Execute()
}