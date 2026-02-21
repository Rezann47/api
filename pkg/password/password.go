package password

import "golang.org/x/crypto/bcrypt"

const cost = 12

// Hash bcrypt ile şifre hash'ler
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(b), err
}

// Compare hash ile plain şifreyi karşılaştırır
func Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
