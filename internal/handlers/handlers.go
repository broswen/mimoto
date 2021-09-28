package handlers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/broswen/mimoto/internal/user"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (sr *SignupRequest) Bind(r *http.Request) error {
	if sr.Email == "" {
		return errors.New("missing email")
	}
	if sr.Name == "" {
		return errors.New("missing name")
	}
	if sr.Password == "" {
		return errors.New("missing password")
	}
	return nil
}

func SignupHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := SignupRequest{}
		if err := render.Bind(r, &data); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		err := userService.Signup(data.Email, data.Name, data.Password)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ConfirmHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		email := r.URL.Query().Get("email")
		if code == "" {
			render.Render(w, r, ErrBadRequest(errors.New("missing confirmation code param")))
			return
		}
		if email == "" {
			render.Render(w, r, ErrBadRequest(errors.New("missing email param")))
			return
		}

		err := userService.Confirm(email, code)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sr *LoginRequest) Bind(r *http.Request) error {
	return nil
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (sr *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func LoginHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := LoginRequest{}
		if err := render.Bind(r, &data); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		token, refreshToken, err := userService.Login(data.Email, data.Password)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		render.Render(w, r, &LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

type RefreshResponse struct {
	Token string `json:"token"`
}

func (sr *RefreshResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func RefreshHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// refreshToken := r.Context().Value("token").(*jwt.Token)
		claims := r.Context().Value("claims").(jwt.StandardClaims)
		refreshTokenString := r.Context().Value("tokenString").(string)

		token, err := userService.Refresh(claims.Subject, refreshTokenString)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		render.Render(w, r, &RefreshResponse{
			Token: token,
		})
	}
}

func LogoutHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(jwt.StandardClaims)

		err := userService.Logout(claims.Subject)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type EmailRequest struct {
	Email string `json:"email"`
}

func (er *EmailRequest) Bind(r *http.Request) error {
	if er.Email == "" {
		return errors.New("missing email")
	}
	return nil
}

func SendResetHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := EmailRequest{}
		if err := render.Bind(r, &data); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		err := userService.SendReset(data.Email)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

func (rpr *ResetPasswordRequest) Bind(r *http.Request) error {
	if rpr.Password == "" {
		return errors.New("missing password")
	}
	return nil
}

func ResetHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ResetPasswordRequest{}
		if err := render.Bind(r, &data); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
		email := r.URL.Query().Get("email")
		code := r.URL.Query().Get("code")
		if email == "" {
			render.Render(w, r, ErrBadRequest(errors.New("missing email param")))
			return
		}
		if code == "" {
			render.Render(w, r, ErrBadRequest(errors.New("missing reset code param")))
			return
		}

		err := userService.ResetPassword(email, data.Password, code)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func JWTAuthorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("authorization"), " ")
		if len(parts) != 2 {
			render.Render(w, r, ErrBadRequest(errors.New("malformed authorization header")))
			return
		}
		tokenString := parts[1]
		if tokenString == "" {
			render.Render(w, r, ErrBadRequest(errors.New("missing jwt")))
			return
		}
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			render.Render(w, r, ErrUnauthorized(err))
			return
		}
		ctx := context.WithValue(r.Context(), "tokenString", tokenString)
		ctx = context.WithValue(ctx, "token", token)
		ctx = context.WithValue(ctx, "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
