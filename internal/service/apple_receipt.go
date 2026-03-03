package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	appleProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
	appleSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	appleAPIProduction = "https://api.storekit.itunes.apple.com"
	appleAPISandbox    = "https://api.storekit-sandbox.itunes.apple.com"
)

// VerifyAppleReceipt — JWS token veya eski base64 receipt'i handle eder
func VerifyAppleReceipt(ctx context.Context, receiptData, sharedSecret string) (bool, *time.Time, string, error) {
	if receiptData == "" {
		return false, nil, "", fmt.Errorf("receipt boş")
	}
	if isJWSToken(receiptData) {
		fmt.Println("ℹ️ [Apple] JWS token tespit edildi, App Store API ile doğrulanıyor")
		txID, err := extractTransactionID(receiptData)
		if err != nil {
			return false, nil, "", err
		}
		return verifyWithAppStoreAPI(ctx, txID)
	}
	fmt.Println("ℹ️ [Apple] Legacy receipt tespit edildi, verifyReceipt ile doğrulanıyor")
	return verifyLegacyReceipt(ctx, receiptData, sharedSecret)
}

func isJWSToken(data string) bool {
	return len(strings.Split(data, ".")) == 3
}

// JWS'den sadece transactionId oku (güvenlik riski yok, sadece ID alıyoruz)
func extractTransactionID(token string) (string, error) {
	parts := strings.Split(token, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("JWS decode hatası: %w", err)
	}
	var p struct {
		TransactionID string `json:"transactionId"`
	}
	if err := json.Unmarshal(decoded, &p); err != nil {
		return "", fmt.Errorf("JWS parse hatası: %w", err)
	}
	if p.TransactionID == "" {
		return "", fmt.Errorf("transactionId boş")
	}
	fmt.Printf("ℹ️ [Apple] transactionId çıkarıldı: %s\n", p.TransactionID)
	return p.TransactionID, nil
}

// Apple App Store Server API ile doğrula — önce sandbox, sonra production
func verifyWithAppStoreAPI(ctx context.Context, transactionID string) (bool, *time.Time, string, error) {
	for _, isSandbox := range []bool{true, false} {
		ok, expires, txID, err := callAppStoreAPI(ctx, transactionID, isSandbox)
		if err != nil {
			fmt.Printf("⚠️ [Apple API] sandbox=%v hata: %v\n", isSandbox, err)
			continue
		}
		return ok, expires, txID, nil
	}
	return false, nil, "", fmt.Errorf("apple API doğrulaması başarısız")
}

func callAppStoreAPI(ctx context.Context, transactionID string, isSandbox bool) (bool, *time.Time, string, error) {
	baseURL := appleAPIProduction
	if isSandbox {
		baseURL = appleAPISandbox
	}

	url := fmt.Sprintf("%s/inApps/v1/transactions/%s", baseURL, transactionID)
	fmt.Printf("⏳ [Apple API] istek gönderiliyor — sandbox: %v | txID: %s\n", isSandbox, transactionID)

	appleJWT, err := generateAppleJWT()
	if err != nil {
		return false, nil, "", fmt.Errorf("JWT oluşturulamadı: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+appleJWT)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return false, nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return false, nil, "", fmt.Errorf("transaction bulunamadı")
	}
	if res.StatusCode != 200 {
		fmt.Printf("❌ [Apple API] status: %d\n", res.StatusCode)
		return false, nil, "", fmt.Errorf("apple API status: %d", res.StatusCode)
	}

	var result struct {
		SignedTransactionInfo string `json:"signedTransactionInfo"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return false, nil, "", err
	}

	ok, expires, txID, err := decodeAppleJWS(result.SignedTransactionInfo)
	if err != nil {
		return false, nil, "", err
	}

	fmt.Printf("✅ [Apple API] doğrulama başarılı — expires: %v\n", expires)
	return ok, expires, txID, nil
}

func decodeAppleJWS(token string) (bool, *time.Time, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, nil, "", fmt.Errorf("geçersiz JWS")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, nil, "", err
	}
	var p struct {
		TransactionID  string `json:"transactionId"`
		ExpiresDate    int64  `json:"expiresDate"`
		RevocationDate int64  `json:"revocationDate"`
	}
	if err := json.Unmarshal(decoded, &p); err != nil {
		return false, nil, "", err
	}
	if p.RevocationDate > 0 {
		return false, nil, "", fmt.Errorf("subscription iptal edilmiş")
	}
	if p.ExpiresDate == 0 {
		return false, nil, "", fmt.Errorf("expiresDate yok")
	}
	expires := time.Unix(p.ExpiresDate/1000, 0)
	return true, &expires, p.TransactionID, nil
}

func generateAppleJWT() (string, error) {
	keyID := os.Getenv("APPLE_KEY_ID")
	issuerID := os.Getenv("APPLE_ISSUER_ID")
	privateKeyStr := os.Getenv("APPLE_PRIVATE_KEY")
	bundleID := "com.kaya.yksrota"

	if keyID == "" || issuerID == "" || privateKeyStr == "" {
		return "", fmt.Errorf("APPLE_KEY_ID, APPLE_ISSUER_ID veya APPLE_PRIVATE_KEY eksik")
	}

	privateKeyStr = strings.ReplaceAll(privateKeyStr, `\n`, "\n")

	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "", fmt.Errorf("private key parse edilemedi")
	}

	keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private key parse hatası: %w", err)
	}

	ecKey, ok := keyInterface.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key ECDSA değil")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": issuerID,
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
		"bid": bundleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyID

	return token.SignedString(ecKey)
}

// ── Eski Base64 Receipt ───────────────────────────────────────────────────────

type appleVerifyRequest struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions"`
}

type appleReceiptInfo struct {
	ProductID        string `json:"product_id"`
	TransactionID    string `json:"transaction_id"`
	ExpiresDateMS    string `json:"expires_date_ms"`
	CancellationDate string `json:"cancellation_date"`
}

type appleVerifyResponse struct {
	Status            int                `json:"status"`
	LatestReceiptInfo []appleReceiptInfo `json:"latest_receipt_info"`
}

func verifyLegacyReceipt(ctx context.Context, receiptData, sharedSecret string) (bool, *time.Time, string, error) {
	result, err := callAppleVerify(ctx, appleProductionURL, receiptData, sharedSecret)
	if err != nil {
		return false, nil, "", err
	}
	if result.Status == 21007 {
		result, err = callAppleVerify(ctx, appleSandboxURL, receiptData, sharedSecret)
		if err != nil {
			return false, nil, "", err
		}
	}
	if result.Status != 0 {
		return false, nil, "", fmt.Errorf("apple receipt status: %d", result.Status)
	}

	var latestExpires *time.Time
	var latestTxID string
	for _, info := range result.LatestReceiptInfo {
		if info.CancellationDate != "" || info.ExpiresDateMS == "" {
			continue
		}
		var ms int64
		fmt.Sscanf(info.ExpiresDateMS, "%d", &ms)
		t := time.Unix(ms/1000, 0)
		if latestExpires == nil || t.After(*latestExpires) {
			latestExpires = &t
			latestTxID = info.TransactionID
		}
	}

	if latestExpires == nil {
		return false, nil, "", fmt.Errorf("aktif subscription bulunamadı")
	}
	return true, latestExpires, latestTxID, nil
}

func callAppleVerify(ctx context.Context, url, receiptData, sharedSecret string) (*appleVerifyResponse, error) {
	body, _ := json.Marshal(appleVerifyRequest{
		ReceiptData:            receiptData,
		Password:               sharedSecret,
		ExcludeOldTransactions: true,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var result appleVerifyResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
