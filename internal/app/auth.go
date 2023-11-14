package app

import (
	"context"
	"dryve/internal/dto"
	"dryve/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type contextKey string

// Context Key for user
var ctxKeyUser contextKey = "user"

// Context Keys for JWT claims basic info
var ctxKeyUserId contextKey = "user_id"
var ctxKeyUserEmail contextKey = "user_email"

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var l dto.LoginRequest
	err := DecodeJSONBody(w, r, &l)
	if err != nil {
		HandleDecodeError(w, err)
		return
	}

	user, err := app.UserService.GetUserByEmail(l.Email)
	if err != nil {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	if !utils.VerifyPassword(user.Password, l.Password) {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	jwtStr, err := utils.GenerateJWT(app.Config.JWT, *user)
	if err != nil {
		logrus.Errorf("cannot create jwt with error '%v'", err)
		http.Error(w, "Cannot create JWT", http.StatusInternalServerError)
		return
	}

	EncodeJSONAndSend(w, dto.LoginResponse{
		Token: jwtStr,
	})
}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	var rr dto.RegisterRequest
	err := DecodeJSONBody(w, r, &rr)
	if err != nil {
		HandleDecodeError(w, err)
		return
	}

	user, err := app.UserService.GetUserByEmail(rr.Email)
	if err == nil && user.Email != "" {
		http.Error(w, "This email has already been registered", http.StatusForbidden)
		return
	}

	if !utils.IsValidPassword(rr.Password) {
		http.Error(w, "Password too weak", http.StatusNotAcceptable)
		return
	}

	rr.Password = utils.HashAndSaltPassword(rr.Password)
	user, err = app.UserService.CreateUser(rr)
	if err != nil {
		logrus.Errorf("cannot create user with error '%v'", err)
		http.Error(w, "Cannot register user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Succesfully created user %s %s [%s] with ID %d", user.FirstName, user.LastName, user.Email, user.ID)))
}

func (app *App) JWTMiddleware(next http.Handler) http.Handler {
	hfn := func(w http.ResponseWriter, r *http.Request) {
		encodedJWT := utils.ExtractJWT(r)
		if encodedJWT == "" {
			logrus.Errorf("cannot extract JWT from the request")
			http.Error(w, "cannot extract JWT from the request", http.StatusBadRequest)
			return
		}

		claims, err := utils.VerifyJWT(app.Config.JWT, encodedJWT)
		if err != nil {
			logrus.Errorf("cannot verify JWT: %v", err)
			http.Error(w, "cannot verify JWT", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Need to first assert the accurate dynamic type of the interface value and then to desired.
		// float64 is because is default for json encoding numbers.
		ctx = context.WithValue(ctx, ctxKeyUserId, uint(claims["UserID"].(float64)))
		ctx = context.WithValue(ctx, ctxKeyUserEmail, claims["UserEmail"])

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(hfn)
}

func (app *App) AuthMiddleware(next http.Handler) http.Handler {
	hfn := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxKeyUserId).(uint)

		user, err := app.UserService.GetUser(id)
		if err != nil {
			logrus.Errorf(err.Error())
			http.Error(w, "error getting user", http.StatusInternalServerError)
			return
		}

		if !user.Verified {
			http.Error(w, "user not verified", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(hfn)
}

func (app *App) EmailVerifyStep1(w http.ResponseWriter, r *http.Request) {
	to := r.Context().Value(ctxKeyUserEmail).(string)
	id := r.Context().Value(ctxKeyUserId).(uint)

	user, err := app.UserService.GetUser(id)
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	code, err := app.UserService.SetEmailConfirmationCode(id)
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, "Error creating a verification code", http.StatusInternalServerError)
		return
	}

	// TODO put here the hostname then
	link := fmt.Sprintf("http://localhost:8666/user/verify/2/email/%d/%s", id, code)

	email := dto.Email{
		From:    app.Config.Email.User,
		To:      to,
		Subject: "Confirmation code for your registration",
		Body:    fmt.Sprintf("Hi %s,<br><br>This is your confirmation link: <b>%s</b>.", user.FirstName, link),
	}

	if err := app.EmailService.SendEmail(email); err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, fmt.Sprintf("Error sending email to %s", email), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email sent"))
}

func (app *App) EmailVerifyStep2(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, "error getting user id in the url param", http.StatusInternalServerError)
		return
	}

	code := chi.URLParam(r, "code")

	user, err := app.UserService.GetUser(uint(id))
	if err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	if code != user.EmailCode {
		logrus.Errorf("wrong verification code %s", code)
		http.Error(w, "wrong verification code", http.StatusUnauthorized)
		return
	}

	if err = app.UserService.VerifyUser(uint(id)); err != nil {
		logrus.Errorf(err.Error())
		http.Error(w, "error verifying user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Email %s succesfully verified", user.Email)))
}
