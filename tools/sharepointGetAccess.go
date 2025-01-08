package tools

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetAccessToken obtiene un access_token usando un certificado desde Azure AD
func GetAccessToken() (string, error) {
	tenantID := os.Getenv("SPAAuth_TENANTID")
	clientID := os.Getenv("SPAAuth_CLIENTID")
	certPath := "SharePointUploaderCertificate.pem"

	if tenantID == "" || clientID == "" {
		return "", fmt.Errorf("missing environment variables: SPAAuth_TENANTID or SPAAuth_CLIENTID")
	}

	certPath, err := filepath.Abs(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve certificate path: %w", err)
	}

	// Leer el certificado desde el archivo
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Decodificar bloques
	var certificate *x509.Certificate
	var privateKey *rsa.PrivateKey

	for {
		block, rest := pem.Decode(certPEM)
		if block == nil {
			break
		}

		log.Printf("Block Type: %s\n", block.Type)

		// Procesar el certificado
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Printf("Failed to parse certificate: %v\n", err)
				return "", fmt.Errorf("failed to parse certificate: %w", err)
			}
			certificate = cert
		}

		// Procesar la clave privada
		if block.Type == "PRIVATE KEY" {
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				log.Printf("Failed to parse private key: %v\n", err)
				return "", fmt.Errorf("failed to parse private key: %w", err)
			}
			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				log.Println("Not an RSA private key")
				return "", fmt.Errorf("not an RSA private key")
			}
			privateKey = rsaKey
		}

		certPEM = rest
	}

	// Validar los resultados
	if certificate == nil {
		log.Println("Certificate not found")
		return "", fmt.Errorf("certificate not found")
	}

	if privateKey == nil {
		log.Println("Private key not found")
		return "", fmt.Errorf("private key not found")
	}

	log.Println("Certificate and private key successfully parsed")
	// Calcular el thumbprint SHA-256 del certificado
	thumbprint := sha256.Sum256(certificate.Raw)
	thumbprintEncoded := base64.RawURLEncoding.EncodeToString(thumbprint[:])

	// Crear el JWT firmado
	now := time.Now()
	claims := jwt.MapClaims{
		"aud": fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		"iss": clientID,
		"sub": clientID,
		"jti": fmt.Sprintf("%d", now.UnixNano()),
		"nbf": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["x5t#S256"] = thumbprintEncoded

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	// Crear la solicitud
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := fmt.Sprintf(
		"client_id=%s&scope=https://%s.sharepoint.com/.default&grant_type=client_credentials&client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer&client_assertion=%s",
		clientID, tenantID, signedToken,
	)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("error creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Enviar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get access token: %s", string(body))
	}

	// Parsear la respuesta
	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", fmt.Errorf("error decoding token response: %w", err)
	}

	log.Printf("Access token: %s\n", tokenResp.AccessToken)
	return tokenResp.AccessToken, nil
}

// package tools

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// )

// type TokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	ExpiresIn   int    `json:"expires_in"`
// 	TokenType   string `json:"token_type"`
// }

// // GetAccessToken obtiene un access_token desde Azure AD
// func GetAccessToken() (string, error) {
// 	tenantID := os.Getenv("SPAAuth_TENANTID")
// 	clientID := os.Getenv("SPAAuth_CLIENTID")
// 	clientSecret := os.Getenv("SPAAuth_CLIENTSECRET")

// 	if tenantID == "" || clientID == "" || clientSecret == "" {
// 		return "", fmt.Errorf("missing environment variables: SPAAuth_TENANTID, SPAAuth_CLIENTID, or SPAAuth_CLIENTSECRET")
// 	}

// 	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
// 	data := fmt.Sprintf(
// 		"client_id=%s&client_secret=%s&grant_type=client_credentials&scope=https://%s.sharepoint.com/.default",
// 		clientID, clientSecret, tenantID,
// 	)

// 	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
// 	if err != nil {
// 		return "", fmt.Errorf("error creating token request: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("error sending token request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		return "", fmt.Errorf("failed to get access token: %s", string(body))
// 	}

// 	var tokenResp TokenResponse
// 	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
// 	if err != nil {
// 		return "", fmt.Errorf("error decoding token response: %w", err)
// 	}

// 	log.Printf("Access token: %s\n", tokenResp.AccessToken)

// 	return tokenResp.AccessToken, nil
// }
