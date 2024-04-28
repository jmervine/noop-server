package server

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

func addMTLSSupportToServer(server *http.Server, ca string) {
	// Add the cert chain as the intermediate signs both the servers and the clients certificates
	clientCACert := []byte(ca)

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	server.TLSConfig = tlsConfig
}

func tlsConfig() (*tls.Config, error) {
	config := &tls.Config{}
	if !cfg.TLSEnabled() {
		return config, nil
	}

	crt, err := tls.X509KeyPair([]byte(cfg.CertPrivatePath), []byte(cfg.CertKeyPath))
	if err != nil {
		return config, err
	}

	config.NextProtos = []string{"http/1.1"}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = crt
	return config, nil
}
