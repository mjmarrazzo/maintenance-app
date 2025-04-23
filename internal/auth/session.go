package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
)

type SessionKey string

const (
	sessionKey = SessionKey("user-session")
)

func getUserFromSession(c echo.Context) (*domain.User, error) {
	sess, err := session.Get(string(sessionKey), c)
	if err != nil {
		return nil, err
	}

	userID, ok := parseInt(sess.Values["ID"])
	if !ok {
		return nil, errors.New("user ID not found in session")
	}
	firstName, ok := parseString(sess.Values["FirstName"])
	if !ok {
		return nil, errors.New("first name not found in session")
	}
	lastName, ok := parseString(sess.Values["LastName"])
	if !ok {
		return nil, errors.New("last name not found in session")
	}
	email, ok := parseString(sess.Values["Email"])
	if !ok {
		return nil, errors.New("email not found in session")
	}

	role, ok := sess.Values["Role"].(string)
	if !ok {
		return nil, errors.New("role not found in session")
	}
	user := &domain.User{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      domain.UserRole(role),
	}

	return user, nil
}

func parseString(value any) (string, bool) {
	if value == nil {
		return "", false
	}
	str, ok := value.(string)
	if !ok {
		return "", false
	}
	return str, true
}

func parseInt(value any) (int64, bool) {
	if value == nil {
		return 0, false
	}
	str, ok := value.(string)
	if !ok {
		return 0, false
	}
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, false
	}
	return intValue, true
}

func getExpirationTime(c echo.Context) (int64, error) {
	s, err := session.Get(string(sessionKey), c)
	if err != nil {
		return 0, err
	}

	expiresAt, ok := parseInt(s.Values["ExpiresAt"])
	if !ok {
		return 0, errors.New("expiresAt not found in session")
	}

	return expiresAt, nil
}

func SaveUserToSession(c echo.Context, user *domain.User) error {
	s, err := session.Get(string(sessionKey), c)
	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	s.Values["ID"] = strconv.FormatInt(user.ID, 10)
	s.Values["FirstName"] = user.FirstName
	s.Values["LastName"] = user.LastName
	s.Values["Email"] = user.Email
	s.Values["Role"] = string(user.Role)

	expirationTime := time.Now().Add(time.Hour * 2).Unix()
	s.Values["ExpiresAt"] = strconv.FormatInt(expirationTime, 10)

	return s.Save(c.Request(), c.Response())
}

func ClearSession(c echo.Context) error {
	s, err := session.Get(string(sessionKey), c)
	if err != nil {
		return err
	}

	s.Values = make(map[any]any)
	s.Options.MaxAge = -1
	return s.Save(c.Request(), c.Response())
}
