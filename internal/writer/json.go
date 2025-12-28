package writer

import (
	"encoding/json"
	"net/http"
)

/* Reusable function to write a Response providing
   the Content Type Header, Status and JSON encoder
*/

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
