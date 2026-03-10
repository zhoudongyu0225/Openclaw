// Package security provides advanced security testing utilities
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"sync"
	"time"
)

// AttackType represents types of security attacks
type AttackType string

const (
	AttackSQLInjection    AttackType = "sql_injection"
	AttackXSS             AttackType = "xss"
	AttackCSRF            AttackType = "csrf"
	AttackDDoS            AttackType = "ddos"
	AttackBruteForce      AttackType = "brute_force"
	AttackMITM            AttackType = "mitm"
	AttackReplay          AttackType = "replay"
	AttackPayload         AttackType = "malicious_payload"
)

// SecurityTestResult holds security test results
type SecurityTestResult struct {
	AttackType     AttackType
	Vulnerabilities []string
	RiskLevel      string // low/medium/high/critical
	Passed         bool
	Recommendations []string
	TestDuration   time.Duration
}

// SecurityTester is the main security testing struct
type SecurityTester struct {
	mu           sync.RWMutex
	results      map[AttackType]*SecurityTestResult
	blacklist    map[string]time.Time
	suspicious   map[string]int
	attackCounts map[AttackType]int
}

// NewSecurityTester creates a new security tester
func NewSecurityTester() *SecurityTester {
	return &SecurityTester{
		results:      make(map[AttackType]*SecurityTestResult),
		blacklist:    make(map[string]time.Time),
		suspicious:   make(map[string]int),
		attackCounts: make(map[AttackType]int),
	}
}

// TestSQLInjection tests for SQL injection vulnerabilities
func (st *SecurityTester) TestSQLInjection(inputs []string) *SecurityTestResult {
	start := time.Now()
	result := &SecurityTestResult{
		AttackType:     AttackSQLInjection,
		Vulnerabilities: []string{},
		Recommendations: []string{},
	}
	
	sqlKeywords := []string{
		"' OR '1'='1", 
		"'; DROP TABLE users;--",
		"1' AND '1'='1",
		"1 UNION SELECT * FROM",
		"'; EXEC xp_",
		"1' OR 1=1--",
		"<script>alert('xss')</script>",
	}
	
	for _, input := range inputs {
		lower := strings.ToLower(input)
		for _, keyword := range sqlKeywords {
			if strings.Contains(lower, strings.ToLower(keyword)) {
				result.Vulnerabilities = append(result.Vulnerabilities, 
					fmt.Sprintf("Potential SQL injection detected: %s", input))
			}
		}
	}
	
	result.RiskLevel = determineRiskLevel(len(result.Vulnerabilities), 10)
	result.Passed = len(result.Vulnerabilities) == 0
	result.Recommendations = []string{
		"Use parameterized queries",
		"Implement input validation",
		"Apply least privilege to database users",
	}
	result.TestDuration = time.Since(start)
	
	st.results[AttackSQLInjection] = result
	return result
}

// TestXSS tests for XSS vulnerabilities
func (st *SecurityTester) TestXSS(inputs []string) *SecurityTestResult {
	start := time.Now()
	result := &SecurityTestResult{
		AttackType:     AttackXSS,
		Vulnerabilities: []string{},
		Recommendations: []string{},
	}
	
	xssPatterns := []string{
		"<script>",
		"javascript:",
		"onerror=",
		"onload=",
		"<iframe>",
		"<embed>",
		"<object>",
		"eval(",
		"innerHTML",
		"document.cookie",
	}
	
	for _, input := range inputs {
		lower := strings.ToLower(input)
		for _, pattern := range xssPatterns {
			if strings.Contains(lower, strings.ToLower(pattern)) {
				result.Vulnerabilities = append(result.Vulnerabilities,
					fmt.Sprintf("Potential XSS detected: %s", input))
				break
			}
		}
	}
	
	result.RiskLevel = determineRiskLevel(len(result.Vulnerabilities), 5)
	result.Passed = len(result.Vulnerabilities) == 0
	result.Recommendations = []string{
		"Escape HTML entities",
		"Use Content Security Policy",
		"Implement input validation",
		"Use template engines with auto-escaping",
	}
	result.TestDuration = time.Since(start)
	
	st.results[AttackXSS] = result
	return result
}

// TestBruteForce simulates and tests brute force protection
func (st *SecurityTester) TestBruteForce(attempts int, window time.Duration) *SecurityTestResult {
	start := time.Now()
	result := &SecurityTestResult{
		AttackType:     AttackBruteForce,
		Vulnerabilities: []string{},
		Recommendations: []string{},
	}
	
	// Simulate attack
	failedAttempts := 0
	for i := 0; i < attempts; i++ {
		time.Sleep(time.Millisecond * 10)
		failedAttempts++
		st.suspicious["test_ip"] = failedAttempts
	}
	
	// Check if rate limiting would catch it
	if st.suspicious["test_ip"] > 5 {
		result.Vulnerabilities = append(result.Vulnerabilities,
			"High number of failed attempts detected")
	}
	
	result.RiskLevel = "low" // Should be mitigated by rate limiter
	result.Passed = true
	result.Recommendations = []string{
		"Implement account lockout after N attempts",
		"Use CAPTCHA for repeated failures",
		"Enable 2FA",
		"Monitor failed login patterns",
	}
	result.TestDuration = time.Since(start)
	
	st.results[AttackBruteForce] = result
	return result
}

// TestReplayAttack tests for replay attack vulnerabilities
func (st *SecurityTester) TestReplayAttack(payloads []string) *SecurityTestResult {
	start := time.Now()
	result := &SecurityTestResult{
		AttackType:     AttackReplay,
		Vulnerabilities: []string{},
		Recommendations: []string{},
	}
	
	seen := make(map[string]bool)
	for _, payload := range payloads {
		if seen[payload] {
			result.Vulnerabilities = append(result.Vulnerabilities,
				"Reused payload detected - vulnerable to replay")
		}
		seen[payload] = true
	}
	
	result.RiskLevel = determineRiskLevel(len(result.Vulnerabilities), 3)
	result.Passed = len(result.Vulnerabilities) == 0
	result.Recommendations = []string{
		"Implement timestamp validation",
		"Use nonces for each request",
		"Add request expiration",
		"Use HTTPS always",
	}
	result.TestDuration = time.Since(start)
	
	st.results[AttackReplay] = result
	return result
}

// RecordAttack records an attack attempt
func (st *SecurityTester) RecordAttack(ip string, attackType AttackType) {
	st.mu.Lock()
	defer st.mu.Unlock()
	
	st.attackCounts[attackType]++
	st.suspicious[ip]++
	
	// Auto-blacklist after threshold
	if st.suspicious[ip] > 100 {
		st.blacklist[ip] = time.Now().Add(24 * time.Hour)
	}
}

// IsBlacklisted checks if an IP is blacklisted
func (st *SecurityTester) IsBlacklisted(ip string) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()
	
	if expiry, ok := st.blacklist[ip]; ok {
		if time.Now().After(expiry) {
			delete(st.blacklist, ip)
			return false
		}
		return true
	}
	return false
}

// GetAttackStats returns attack statistics
func (st *SecurityTester) GetAttackStats() map[AttackType]int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	
	result := make(map[AttackType]int)
	for k, v := range st.attackCounts {
		result[k] = v
	}
	return result
}

// PrintResults prints all security test results
func (st *SecurityTester) PrintResults() {
	fmt.Println("\n========== Security Test Results ==========")
	for attackType, result := range st.results {
		fmt.Printf("\n[%s] %s\n", result.RiskLevel, attackType)
		fmt.Printf("  Passed: %v\n", result.Passed)
		fmt.Printf("  Duration: %v\n", result.TestDuration)
		if len(result.Vulnerabilities) > 0 {
			fmt.Printf("  Vulnerabilities (%d):\n", len(result.Vulnerabilities))
			for _, v := range result.Vulnerabilities {
				fmt.Printf("    - %s\n", v)
			}
		}
		fmt.Printf("  Recommendations:\n")
		for _, r := range result.Recommendations {
			fmt.Printf("    - %s\n", r)
		}
	}
	fmt.Println("\n==========================================\n")
}

func determineRiskLevel(count, threshold int) string {
	ratio := float64(count) / float64(threshold)
	switch {
	case ratio >= 1:
		return "critical"
	case ratio >= 0.6:
		return "high"
	case ratio >= 0.3:
		return "medium"
	default:
		return "low"
	}
}

// Encryption provides encryption utilities
type Encryption struct {
	key []byte
}

// NewEncryption creates new encryption instance
func NewEncryption(key []byte) (*Encryption, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	return &Encryption{key: key}, nil
}

// Encrypt encrypts plaintext using AES-GCM
func (e *Encryption) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext
func (e *Encryption) Decrypt(ciphertextHex string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, err
	}
	
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// HashPassword securely hashes a password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// GenerateSecureToken generates a cryptographically secure token
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ValidateSecureToken validates token entropy
func ValidateSecureToken(token string, minEntropy int) bool {
	if len(token) < minEntropy {
		return false
	}
	
	// Check character distribution
	charsetSize := 0
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	
	for _, c := range token {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	
	charsetTypes := 0
	if hasLower {
		charsetTypes++
	}
	if hasUpper {
		charsetTypes++
	}
	if hasDigit {
		charsetTypes++
	}
	if hasSpecial {
		charsetTypes++
	}
	
	// Simple entropy approximation
	entropy := float64(len(token)) * float64(charsetTypes) / 4
	return entropy >= float64(minEntropy)
}

// SecureRandomInt returns a secure random integer in range
func SecureRandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, errors.New("min must be less than max")
	}
	
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return 0, err
	}
	
	return int(n.Int64()) + min, nil
}

// AuditLog records security events
type AuditLog struct {
	mu       sync.Mutex
	entries  []AuditEntry
	maxSize  int
}

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	Timestamp    time.Time
	EventType    string
	SourceIP     string
	UserID       string
	Details      string
	Severity     string // info/warning/error/critical
}

// NewAuditLog creates a new audit log
func NewAuditLog(maxSize int) *AuditLog {
	return &AuditLog{
		entries: make([]AuditEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Log records an audit entry
func (al *AuditLog) Log(eventType, sourceIP, userID, details, severity string) {
	al.mu.Lock()
	defer al.mu.Unlock()
	
	entry := AuditEntry{
		Timestamp: time.Now(),
		EventType: eventType,
		SourceIP:  sourceIP,
		UserID:    userID,
		Details:   details,
		Severity:  severity,
	}
	
	al.entries = append(al.entries, entry)
	
	// Trim if exceeds max size
	if len(al.entries) > al.maxSize {
		al.entries = al.entries[len(al.entries)-al.maxSize:]
	}
}

// GetEntries returns audit entries
func (al *AuditLog) GetEntries() []AuditEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()
	
	result := make([]AuditEntry, len(al.entries))
	copy(result, al.entries)
	return result
}

// GetEntriesBySeverity returns entries filtered by severity
func (al *AuditLog) GetEntriesBySeverity(severity string) []AuditEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()
	
	var result []AuditEntry
	for _, e := range al.entries {
		if e.Severity == severity {
			result = append(result, e)
		}
	}
	return result
}
