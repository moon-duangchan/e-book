package controller

import (
    "Backend/auth"
    "Backend/database"
    "Backend/models"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "net/mail"
    "os"
    "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type registerInput struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Register creates a new user, stores verification token, and emails verification link
func Register(c *fiber.Ctx) error {
    var in registerInput
    if err := c.BodyParser(&in); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
    }
    in.Email = strings.TrimSpace(strings.ToLower(in.Email))
    if in.Email == "" || in.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
    }
    if _, err := mail.ParseAddress(in.Email); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email"})
    }

    db := database.DBConn

    // Check if email already exists
    var count int64
    if err := db.Model(&models.User{}).Where("email = ?", in.Email).Count(&count).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
    }
    if count > 0 {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
    }

    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
    }

    // Generate verification token
    token, err := generateToken(32)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
    }

    user := models.User{
        Name:                       in.Name,
        Email:                      in.Email,
        PasswordHash:               string(hash),
        Verified:                   false,
        VerificationToken:          token,
        VerificationTokenExpiresAt: time.Now().Add(24 * time.Hour),
    }

    if err := db.Create(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
    }

    // Build verification link
    baseURL := os.Getenv("APP_BASE_URL")
    if baseURL == "" {
        port := os.Getenv("PORT")
        if port == "" { port = "3001" }
        baseURL = fmt.Sprintf("http://localhost:%s", port)
    }
    link := fmt.Sprintf("%s/auth/verify?token=%s", strings.TrimRight(baseURL, "/"), token)

    // Send email (best-effort; do not leak internal errors)
    if err := auth.SendVerificationEmail(user.Email, link); err != nil {
        // Log-like response without exposing secrets
        fmt.Println("mail send error:", err)
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "registration successful; please verify your email",
    })
}

// VerifyEmail confirms the user's email using a token query parameter
func VerifyEmail(c *fiber.Ctx) error {
    token := c.Query("token")
    if token == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token required"})
    }

    db := database.DBConn
    var user models.User
    if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
    }
    if time.Now().After(user.VerificationTokenExpiresAt) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
    }

    user.Verified = true
    user.VerificationToken = ""
    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to verify"})
    }
    return c.JSON(fiber.Map{"message": "email verified"})
}

// Login authenticates a user and returns a JWT if verified
func Login(c *fiber.Ctx) error {
    var in loginInput
    if err := c.BodyParser(&in); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
    }
    email := strings.TrimSpace(strings.ToLower(in.Email))
    if email == "" || in.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
    }

    db := database.DBConn
    var user models.User
    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials or someone already login"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
    }
    if !user.Verified {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "please verify your email"})
    }

    token, err := signJWT(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to sign token"})
    }

    return c.JSON(fiber.Map{"token": token})
}

func generateToken(n int) (string, error) {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func signJWT(u models.User) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("missing JWT_SECRET")
    }
    claims := jwt.MapClaims{
        "sub": u.ID,
        "email": u.Email,
        "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

