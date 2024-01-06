package random

import "crypto/rand"

// Generate is a helper function for getting []byte particular size
func Generate(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
