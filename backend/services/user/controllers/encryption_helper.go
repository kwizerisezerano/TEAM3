package controllers

import (
	"github.com/kwizera/university-invoices/internal/encryption"
)

// DecryptEmail decrypts an encrypted email address
func DecryptEmail(encryptedEmail string) (string, error) {
	key := encryption.GetEncryptionKey()
	return encryption.Decrypt(encryptedEmail, key)
}

// EncryptEmail encrypts an email address
func EncryptEmail(email string) (string, error) {
	key := encryption.GetEncryptionKey()
	return encryption.Encrypt(email, key)
}
