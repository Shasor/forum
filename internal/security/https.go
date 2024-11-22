package security

import (
	"crypto/tls"
)

// LoadTLSConfig sets up secure cipher suites and forces TLS 1.2 or higher.
func LoadTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, // Specify your desired cipher
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
