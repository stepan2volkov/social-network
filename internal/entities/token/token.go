package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/stepan2volkov/social-network/internal/entities/user"
)

var profileKey = struct{}{}

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
)

type Claims struct {
	user.UserProfile
	jwt.RegisteredClaims
}

func CreateClaims(issuer string, expiredIn time.Duration, u user.UserProfile) Claims {
	return Claims{
		UserProfile: u,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(expiredIn),
			},
		},
	}
}

func GetClaims(token string, secretKey []byte) (Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
	if err != nil {
		return Claims{}, ErrTokenInvalid
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return Claims{}, ErrTokenInvalid
	}
	if claims == nil {
		return Claims{}, ErrTokenInvalid
	}

	return *claims, nil
}

func CreateToken(claims Claims, secretKey []byte) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func UpdateToken(
	token string,
	secretKey []byte,
	expiredIn time.Duration,
) (
	newToken string,
	err error,
) {
	claims, err := GetClaims(token, secretKey)
	if err != nil {
		return "", err
	}
	claims.ExpiresAt = &jwt.NumericDate{Time: time.Now().Add(expiredIn)}

	return CreateToken(claims, secretKey)
}

func CreateCtxWithProfile(ctx context.Context, u user.UserProfile) context.Context {
	return context.WithValue(ctx, profileKey, u)
}

func GetProfileFromCtx(ctx context.Context) (user.UserProfile, error) {
	profile, ok := ctx.Value(profileKey).(user.UserProfile)
	if !ok {
		return user.UserProfile{}, ErrTokenInvalid
	}

	return profile, nil
}
