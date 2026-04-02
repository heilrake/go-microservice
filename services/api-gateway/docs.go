package main

// @title           Ride Sharing API Gateway
// @version         1.0
// @description     HTTP entry point for the Ride Sharing microservice platform.
// @host            localhost:8081
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT token (format: "Bearer <token>"). Required for /driver/* endpoints.

// postAuthOAuth is a swagger documentation stub for the OAuth login route (proxied to auth-service).
// @Summary      OAuth login
// @Description  Authenticate via Google OAuth. On success returns a signed JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body oauthLoginRequest true "OAuth code and provider info"
// @Success      200  {object}  map[string]string  "{ token: <jwt> }"
// @Failure      400  {string}  string             "Bad request"
// @Failure      502  {string}  string             "Auth service unavailable"
// @Router       /auth/oauth [post]
func postAuthOAuth() {} //nolint:deadcode,unused

// postDevLogin is a swagger documentation stub for the dev login route (proxied to auth-service).
// @Summary      Dev login (development only)
// @Description  Returns a signed JWT for a seeded test account. Only available in dev environment.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body devLoginRequest true "Role and seed number"
// @Success      200  {object}  map[string]string  "{ token: <jwt> }"
// @Failure      400  {string}  string             "Bad request"
// @Router       /dev/login [post]
func postDevLogin() {} //nolint:deadcode,unused
