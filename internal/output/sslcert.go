package output

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
)

// SSLCertMode controls how sslcert measurement certificates are rendered.
type SSLCertMode int

const (
	// SSLCertNames prints only the certificate names (subject/issuer CN) and
	// the Subject Alternative Names.
	SSLCertNames SSLCertMode = iota
	// SSLCertFull prints the base64 (PEM) certificate plus a fully decoded
	// human-readable text dump.
	SSLCertFull
)

// ParseSSLCertMode maps a user supplied flag value to an SSLCertMode.
// Unknown values fall back to the concise names-only mode.
func ParseSSLCertMode(s string) SSLCertMode {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "full":
		return SSLCertFull
	default:
		return SSLCertNames
	}
}

// parseSSLCert decodes a RIPE Atlas certificate string into an *x509.Certificate.
// Atlas firmware sometimes returns the PEM body with literal "\n" escape
// sequences instead of real newlines, so we normalize those first.
func parseSSLCert(raw string) (string, *x509.Certificate, error) {
	pemText := strings.ReplaceAll(raw, "\\n", "\n")
	block, _ := pem.Decode([]byte(pemText))
	if block == nil {
		return pemText, nil, fmt.Errorf("no PEM block found")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	// Re-encode to a clean, canonical PEM block for display.
	pemText = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: block.Bytes}))
	return pemText, cert, err
}

// collectSANs returns all Subject Alternative Names across the supported types.
func collectSANs(cert *x509.Certificate) []string {
	var sans []string
	for _, dns := range cert.DNSNames {
		sans = append(sans, "DNS:"+dns)
	}
	for _, ip := range cert.IPAddresses {
		sans = append(sans, "IP:"+ip.String())
	}
	for _, email := range cert.EmailAddresses {
		sans = append(sans, "email:"+email)
	}
	for _, uri := range cert.URIs {
		sans = append(sans, "URI:"+uri.String())
	}
	return sans
}

// formatSSLCert renders a single certificate according to the selected mode.
func formatSSLCert(w io.Writer, idx int, raw string, mode SSLCertMode) {
	pemText, cert, err := parseSSLCert(raw)

	switch mode {
	case SSLCertFull:
		formatSSLCertFull(w, idx, pemText, cert, err)
	default:
		formatSSLCertNames(w, idx, cert, err)
	}
}

func formatSSLCertNames(w io.Writer, idx int, cert *x509.Certificate, err error) {
	if err != nil || cert == nil {
		fmt.Fprintf(w, "  Cert[%d]:     <unparseable certificate: %v>\n", idx, err)
		return
	}
	fmt.Fprintf(w, "  Cert[%d]:\n", idx)
	fmt.Fprintf(w, "    Subject CN: %s\n", cert.Subject.CommonName)
	fmt.Fprintf(w, "    Issuer CN:  %s\n", cert.Issuer.CommonName)
	if sans := collectSANs(cert); len(sans) > 0 {
		fmt.Fprintf(w, "    SANs:       %s\n", strings.Join(sans, ", "))
	}
}

func formatSSLCertFull(w io.Writer, idx int, pemText string, cert *x509.Certificate, err error) {
	fmt.Fprintf(w, "  Cert[%d]:\n", idx)
	fmt.Fprintln(w, "    Base64 (PEM):")
	for _, line := range strings.Split(strings.TrimRight(pemText, "\n"), "\n") {
		fmt.Fprintf(w, "      %s\n", line)
	}

	fmt.Fprintln(w, "    Decoded:")
	if err != nil || cert == nil {
		fmt.Fprintf(w, "      <unparseable certificate: %v>\n", err)
		return
	}

	fmt.Fprintf(w, "      Version:             %d\n", cert.Version)
	fmt.Fprintf(w, "      Serial Number:       %x\n", cert.SerialNumber)
	fmt.Fprintf(w, "      Signature Algorithm: %s\n", cert.SignatureAlgorithm)
	fmt.Fprintf(w, "      Issuer:              %s\n", cert.Issuer)
	fmt.Fprintf(w, "      Validity:\n")
	fmt.Fprintf(w, "        Not Before:        %s\n", cert.NotBefore.UTC().Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "        Not After:         %s\n", cert.NotAfter.UTC().Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "      Subject:             %s\n", cert.Subject)
	fmt.Fprintf(w, "      Public Key:          %s\n", publicKeyDescription(cert))
	if sans := collectSANs(cert); len(sans) > 0 {
		fmt.Fprintf(w, "      SANs:                %s\n", strings.Join(sans, ", "))
	}
	fmt.Fprintf(w, "      Is CA:               %t\n", cert.IsCA)
	if usages := keyUsages(cert.KeyUsage); len(usages) > 0 {
		fmt.Fprintf(w, "      Key Usage:           %s\n", strings.Join(usages, ", "))
	}
	if extUsages := extKeyUsages(cert.ExtKeyUsage); len(extUsages) > 0 {
		fmt.Fprintf(w, "      Ext Key Usage:       %s\n", strings.Join(extUsages, ", "))
	}
	if len(cert.SubjectKeyId) > 0 {
		fmt.Fprintf(w, "      Subject Key ID:      %x\n", cert.SubjectKeyId)
	}
	if len(cert.AuthorityKeyId) > 0 {
		fmt.Fprintf(w, "      Authority Key ID:    %x\n", cert.AuthorityKeyId)
	}
}

func publicKeyDescription(cert *x509.Certificate) string {
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		return fmt.Sprintf("RSA %d bits", pub.N.BitLen())
	case *ecdsa.PublicKey:
		return fmt.Sprintf("ECDSA %s (%d bits)", pub.Curve.Params().Name, pub.Curve.Params().BitSize)
	case ed25519.PublicKey:
		return "Ed25519"
	default:
		return cert.PublicKeyAlgorithm.String()
	}
}

func keyUsages(usage x509.KeyUsage) []string {
	var out []string
	for _, ku := range []struct {
		bit  x509.KeyUsage
		name string
	}{
		{x509.KeyUsageDigitalSignature, "Digital Signature"},
		{x509.KeyUsageContentCommitment, "Content Commitment"},
		{x509.KeyUsageKeyEncipherment, "Key Encipherment"},
		{x509.KeyUsageDataEncipherment, "Data Encipherment"},
		{x509.KeyUsageKeyAgreement, "Key Agreement"},
		{x509.KeyUsageCertSign, "Certificate Sign"},
		{x509.KeyUsageCRLSign, "CRL Sign"},
		{x509.KeyUsageEncipherOnly, "Encipher Only"},
		{x509.KeyUsageDecipherOnly, "Decipher Only"},
	} {
		if usage&ku.bit != 0 {
			out = append(out, ku.name)
		}
	}
	return out
}

func extKeyUsages(usages []x509.ExtKeyUsage) []string {
	names := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageAny:             "Any",
		x509.ExtKeyUsageServerAuth:      "Server Auth",
		x509.ExtKeyUsageClientAuth:      "Client Auth",
		x509.ExtKeyUsageCodeSigning:     "Code Signing",
		x509.ExtKeyUsageEmailProtection: "Email Protection",
		x509.ExtKeyUsageTimeStamping:    "Time Stamping",
		x509.ExtKeyUsageOCSPSigning:     "OCSP Signing",
	}
	var out []string
	for _, u := range usages {
		if name, ok := names[u]; ok {
			out = append(out, name)
		} else {
			out = append(out, fmt.Sprintf("Unknown(%d)", u))
		}
	}
	return out
}
