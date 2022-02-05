package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/pyroscope-io/pyroscope/pkg/model"
)

type AuthService struct {
	db              *gorm.DB
	userService     UserService
	apiKeyService   APIKeyService
	jwtTokenService JWTTokenService
}

func NewAuthService(db *gorm.DB, jwtTokenService JWTTokenService) AuthService {
	return AuthService{
		db:              db,
		userService:     NewUserService(db),
		apiKeyService:   NewAPIKeyService(db, jwtTokenService),
		jwtTokenService: jwtTokenService,
	}
}

// AuthenticateUser returns User with the given login if its password hash
// matches the given one. If user cannot be found or the password does not
// match the function returns ErrInvalidCredentials.
//
// External users are not allowed to use password authentication.
// TODO(kolesnikovae): It's true for "some" authentication providers.
//  Others may need us to pass through the credentials (e.g. LDAP).
func (svc AuthService) AuthenticateUser(ctx context.Context, name string, password string) (model.User, error) {
	user, err := svc.userService.FindUserByName(ctx, name)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrUserNotFound):
		return model.User{}, model.ErrInvalidCredentials
	default:
		return model.User{}, err
	}
	if model.IsUserExternal(user) {
		// TODO(kolesnikovae): We could be a bit more specific here
		//  disclosing some info: e.g. tell which auth provider to use.
		return model.User{}, model.ErrInvalidCredentials
	}
	if model.IsUserDisabled(user) {
		return model.User{}, model.ErrUserDisabled
	}
	if err = model.VerifyPassword(user.PasswordHash, password); err != nil {
		return model.User{}, model.ErrInvalidCredentials
	}
	return user, nil
}

func (svc AuthService) APIKeyTokenFromJWTToken(ctx context.Context, t string) (model.APIKeyToken, error) {
	token, err := svc.jwtTokenService.Parse(t)
	if err != nil {
		return model.APIKeyToken{}, fmt.Errorf("invalid jwt token: %w", err)
	}
	keyToken, ok := svc.jwtTokenService.APIKeyFromJWTToken(token)
	if !ok {
		return model.APIKeyToken{}, fmt.Errorf("api key is invalid")
	}
	apiKey, err := svc.apiKeyService.FindAPIKeyByName(ctx, keyToken.Name)
	if err != nil {
		return model.APIKeyToken{}, err
	}
	if !apiKey.VerifySignature(token) {
		return model.APIKeyToken{}, fmt.Errorf("api key signature mismatch")
	}
	return keyToken, nil
}

func (svc AuthService) UserFromJWTToken(ctx context.Context, t string) (model.User, error) {
	token, err := svc.jwtTokenService.Parse(t)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid jwt token: %w", err)
	}
	userToken, ok := svc.jwtTokenService.UserFromJWTToken(token)
	if !ok {
		return model.User{}, fmt.Errorf("user token is invalid")
	}
	user, err := svc.userService.FindUserByName(ctx, userToken.Name)
	if err != nil {
		return model.User{}, err
	}
	if model.IsUserDisabled(user) {
		return model.User{}, model.ErrUserDisabled
	}
	if user.Role != userToken.Role {
		return model.User{}, fmt.Errorf("user role has changed")
	}
	return user, nil
}
