package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthorizationHandler struct{}

//type httpOIDCAuthenticationRequest {
//
//}

func (h *AuthorizationHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Print(query)
	json.NewEncoder(w).Encode(query)
}
