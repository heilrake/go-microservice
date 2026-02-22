package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	pb "ride-sharing/shared/proto/user"
)

// googleTokenResponse is the response from Google token endpoint
type googleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// googleUserInfo is the response from Google UserInfo API
type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// exchangeGoogleCode exchanges the OAuth code for an access token
func exchangeGoogleCode(code, redirectURI, clientID, clientSecret string) (string, error) {
	params := url.Values{}
	params.Set("code", code)
	params.Set("redirect_uri", redirectURI)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("grant_type", "authorization_code")

	resp, err := http.Post(
		"https://oauth2.googleapis.com/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp googleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// fetchGoogleUserInfo fetches user profile from Google
func fetchGoogleUserInfo(accessToken string) (*googleUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed (status %d): %s", resp.StatusCode, string(body))
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo: %w", err)
	}

	return &userInfo, nil
}

// oauthRequest is the request body for POST /auth/oauth
type oauthRequest struct {
	Code        string `json:"code"`
	Provider    string `json:"provider"`
	Role        string `json:"role"`
	RedirectURI string `json:"redirect_uri"`
}

// OAuthHandler handles OAuth login (Google, etc.)
func OAuthHandler(userClient *UserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req oauthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Code == "" || req.Provider == "" || req.Role == "" {
			http.Error(w, "code, provider and role required", http.StatusBadRequest)
			return
		}

		if req.Role != "rider" && req.Role != "driver" {
			http.Error(w, "invalid role: must be rider or driver", http.StatusBadRequest)
			return
		}

		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			log.Printf("auth-service: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET not set")
			http.Error(w, "OAuth not configured", http.StatusInternalServerError)
			return
		}

		if req.RedirectURI == "" {
			req.RedirectURI = os.Getenv("OAUTH_REDIRECT_URI")
		}
		if req.RedirectURI == "" {
			http.Error(w, "redirect_uri required", http.StatusBadRequest)
			return
		}

		switch req.Provider {
		case "google":
			// 1. Exchange code for access token
			accessToken, err := exchangeGoogleCode(req.Code, req.RedirectURI, clientID, clientSecret)
			if err != nil {
				log.Printf("auth-service: Google token exchange failed: %v", err)
				http.Error(w, "OAuth token exchange failed", http.StatusBadRequest)
				return
			}

			// 2. Fetch user profile from Google
			googleUser, err := fetchGoogleUserInfo(accessToken)
			if err != nil {
				log.Printf("auth-service: Google userinfo failed: %v", err)
				http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
				return
			}

			if googleUser.Email == "" {
				http.Error(w, "Google account has no email", http.StatusBadRequest)
				return
			}

			username := googleUser.Name
			if username == "" {
				username = googleUser.Email
			}

			// 3. Get or create user in user-service
			resp, err := userClient.Client.GetOrCreateUserByOAuth(r.Context(), &pb.GetOrCreateUserByOAuthRequest{
				Email:          googleUser.Email,
				Username:       username,
				ProfilePicture: googleUser.Picture,
				Role:           req.Role,
			})
			if err != nil {
				log.Printf("auth-service: GetOrCreateUserByOAuth failed: %v", err)
				http.Error(w, "Failed to create or find user", http.StatusInternalServerError)
				return
			}

			// 4. Sign JWT (rider -> "user", driver -> "driver")
			jwtRole := "user"
			if req.Role == "driver" {
				jwtRole = "driver"
			}
			token, err := SignToken(resp.User.Id, jwtRole)
			if err != nil {
				log.Printf("auth-service: SignToken failed: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			// 5. Return same format as email/password login
			out := loginResponse{}
			out.Data.User = resp.User
			out.Data.Token = token
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(out)

		default:
			http.Error(w, "unsupported provider: "+req.Provider, http.StatusBadRequest)
		}
	}
}
