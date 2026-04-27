package middleware

import (
	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Protected(cfg *config.Config) fiber.Handler {
	secret := cfg.SupabaseJWTSecret

	if secret == "" {
		panic("SUPABASE_JWT_SECRET env variable is required.")
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(secret),
			JWTAlg: "ES256", // Use ES256 for Supabase JWTs
		},
		JWKSetURLs: []string{cfg.SupabaseURL + "/auth/v1/.well-known/jwks.json"},
		Extractor:  extractors.FromAuthHeader("Bearer"),
		SuccessHandler: func(c fiber.Ctx) error {
			// extract claims and store userID in locals
			token := jwtware.FromContext(c)
			fiberlog.Debug("Token claims: ", token.Claims)
			claims := token.Claims.(jwt.MapClaims)

			// make sure "sub" exists before casting
			sub, ok := claims["sub"].(string)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "invalid token claims",
				})
			}

			user := &domain.User{}
			user.ID = uuid.MustParse(sub)

			c.Locals("user", user)
			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			fiberlog.Debug("JWT error: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid or expired token",
			})
		},
	})
}
