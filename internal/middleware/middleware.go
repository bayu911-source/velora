package middleware

import (
    "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "go.uber.org/zap"

    "velora/internal/services"
    "velora/internal/pkg/utils/response"
)

func Logger(logger *zap.Logger) fiber.Handler {
    return fiberlogger.New(fiberlogger.Config{
        Format: "${time} ${status} - ${method} ${path} ${latency} ${locals:tenant_id}\n",
    })
}

func CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Content-Type, Authorization, X-Tenant-ID",
        AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
    })
}

func RateLimiter(max int, window time.Duration) fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        max,
        Expiration: window,
        KeyGenerator: func(c *fiber.Ctx) string {
            if auth := c.Get("Authorization"); auth != "" {
                return auth
            }
            if tenantID := c.Get("X-Tenant-ID"); tenantID != "" {
                return tenantID
            }
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return response.Error(c, fiber.StatusTooManyRequests, fiber.ErrTooManyRequests)
        },
    })
}

func Auth(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return response.Error(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "missing authorization header"))
        }
        tokenText := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.ParseWithClaims(tokenText, &services.TokenClaims{}, func(token *jwt.Token) (any, error) {
            if token.Method != jwt.SigningMethodHS256 {
                return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return response.Error(c, fiber.StatusUnauthorized, err)
        }

        claims, ok := token.Claims.(*services.TokenClaims)
        if !ok {
            return response.Error(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "invalid token claims"))
        }
        c.Locals("user_id", claims.UserID)
        c.Locals("tenant_id", claims.TenantID)
        c.Locals("role", claims.Role)
        c.Locals("user_claims", claims)
        return c.Next()
    }
}

func Tenant() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tenantID := c.Get("X-Tenant-ID")
        if tenantID == "" {
            if userTenant := c.Locals("tenant_id"); userTenant != nil {
                tenantID = userTenant.(string)
            }
        }
        if tenantID == "" {
            return response.Error(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "missing X-Tenant-ID header"))
        }
        if _, err := uuid.Parse(tenantID); err != nil {
            return response.Error(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid X-Tenant-ID"))
        }
        c.Locals("tenant_id", tenantID)
        return c.Next()
    }
}

func Role(requiredRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("role").(string)
        if !ok || role == "" {
            return response.Error(c, fiber.StatusForbidden, fiber.NewError(fiber.StatusForbidden, "role required"))
        }
        for _, required := range requiredRoles {
            if strings.EqualFold(required, role) {
                return c.Next()
            }
        }
        return response.Error(c, fiber.StatusForbidden, fiber.NewError(fiber.StatusForbidden, "insufficient role permissions"))
    }
}
