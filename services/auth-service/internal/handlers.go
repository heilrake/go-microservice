package internal

import (
	"encoding/json"
	"log"
	"net/http"

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

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Data struct {
		User  *pb.User `json:"user"`
		Token string   `json:"token,omitempty"`
	} `json:"data"`
}

func UserLoginHandler(userClient *UserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req userLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}

		resp, err := userClient.Client.LoginUser(r.Context(), &pb.LoginUserRequest{
			Email:    req.Email,
			Password: req.Password,
			Role:     "rider",
		})
		if err != nil {
			log.Printf("auth-service: LoginUser: %v", err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := SignToken(resp.User.Id, "user")
		if err != nil {
			log.Printf("auth-service: SignToken: %v", err)
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

type driverLoginResponse struct {
	Data struct {
		Driver *driverResponse `json:"driver"`
		Token  string          `json:"token,omitempty"`
	} `json:"data"`
}

type driverResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

func DriverLoginHandler(userClient *UserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req userLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}

		resp, err := userClient.Client.LoginUser(r.Context(), &pb.LoginUserRequest{
			Email:    req.Email,
			Password: req.Password,
			Role:     "driver",
		})
		if err != nil {
			log.Printf("auth-service: DriverLoginUser: %v", err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := SignToken(resp.User.Id, "driver")
		if err != nil {
			log.Printf("auth-service: SignToken: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		driver := &driverResponse{
			ID:             resp.User.Id,
			Name:           resp.User.Username,
			Email:          resp.User.Email,
			ProfilePicture: resp.User.ProfilePicture,
		}
		out := driverLoginResponse{}
		out.Data.Driver = driver
		out.Data.Token = token
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}
