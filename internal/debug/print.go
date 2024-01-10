package debug

import (
	"encoding/json"
	"log"
)

func Print(data any) {
	marshal, _ := json.MarshalIndent(data, "", " ")
	log.Printf("%s", string(marshal))
}
