package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	pb "ride-sharing/shared/proto/user"
)

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

type loginResponse struct {
	Data struct {
		User  *pb.User `json:"user"`
		Token string   `json:"token,omitempty"`
	} `json:"data"`
}

// devLoginRequest is the request body for POST /dev/login
type devLoginRequest struct {
	Role string `json:"role"` // "driver" or "rider"
	Seed int    `json:"seed"` // 1..N — creates separate test users per seed
}

// DevLoginHandler creates or reuses a fake user and returns a signed JWT.
// Only works when ENVIRONMENT=development.
func DevLoginHandler(userClient *UserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENVIRONMENT") != "development" {
			http.Error(w, "dev login only available in development", http.StatusForbidden)
			return
		}

		var req devLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Role != "rider" && req.Role != "driver" {
			http.Error(w, "role must be 'driver' or 'rider'", http.StatusBadRequest)
			return
		}
		if req.Seed < 1 {
			req.Seed = 1
		}

		email := fmt.Sprintf("dev-%s-%d@dev.local", req.Role, req.Seed)
		name := fmt.Sprintf("Dev %s %d", capitalize(req.Role), req.Seed)

		resp, err := userClient.Client.GetOrCreateUserByOAuth(r.Context(), &pb.GetOrCreateUserByOAuthRequest{
			Email:    email,
			Username: name,
			Role:     req.Role,
		})
		if err != nil {
			log.Printf("dev-login: GetOrCreateUserByOAuth: %v", err)
			http.Error(w, "failed to get or create dev user", http.StatusInternalServerError)
			return
		}

		jwtRole := "user"
		if req.Role == "driver" {
			jwtRole = "driver"
		}
		token, err := SignToken(resp.User.Id, jwtRole)
		if err != nil {
			log.Printf("dev-login: SignToken: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		out := loginResponse{}
		out.Data.User = resp.User
		out.Data.Token = token
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
