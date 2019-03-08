package crypter

// Crypter defines the interface for various encryption tools
type Crypter interface {
	Encrypt(data string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
