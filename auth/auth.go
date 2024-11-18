package auth

import (
	"log"
    "context"
	"os"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ProviderURL  string
}

var (
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
)

func InitOIDC() {
    clientID := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    redirectURI := os.Getenv("REDIRECT_URI")
    providerURL := os.Getenv("PROVIDER_URL")

    log.Printf("Loaded config: ClientID: %s, ClientSecret: %s, RedirectURI: %s", clientID, clientSecret, redirectURI)

    if clientID == "" || clientSecret == "" || redirectURI == "" {
        log.Fatal("Missing CLIENT_ID, CLIENT_SECRET, or REDIRECT_URI environment variable")
    }

    if providerURL == "" {
        providerURL = "https://accounts.google.com"
    }

    oauth2Config = &oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURI,
        Scopes:       []string{"openid", "email", "profile"},
        Endpoint:     google.Endpoint,
    }

    var err error
    provider, err = oidc.NewProvider(context.Background(), providerURL)
    if err != nil {
        log.Fatalf("Failed to get OIDC provider: %v", err)
    }

    verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
}

// GetOAuth2Config returns the OAuth2 configuration
func GetOAuth2Config() *oauth2.Config {
	return oauth2Config
}

// GetVerifier returns the IDTokenVerifier
func GetVerifier() *oidc.IDTokenVerifier {
	return verifier
}

// OIDCAuthMiddleware is a middleware that ensures the user is authenticated
func OIDCAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := verifier.Verify(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		log.Printf("Authorization header: %s", authHeader)
		log.Printf("Extracted token: %s", token)

		c.Next()
	}
}
