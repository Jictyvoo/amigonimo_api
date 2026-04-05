package smtp

import "fmt"

type Encryption uint8

const (
	EncryptionTLS  Encryption = iota // STARTTLS
	EncryptionSSL                    // SSL/TLS on connect
	EncryptionNone                   // Plain SMTP
)

// Config holds the credentials and connection settings for an SMTP server.
type Config struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	Encryption  Encryption
	FromAddress string
	FromName    string
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) SenderHeader() string {
	if c.FromName == "" {
		return c.FromAddress
	}
	return fmt.Sprintf("%s <%s>", c.FromName, c.FromAddress)
}
