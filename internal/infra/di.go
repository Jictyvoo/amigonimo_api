package infra

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/smtp"
	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/stub"
	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/webhook"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func RegisterInfraServices(inj remy.Injector) {
	repositories.RegisterRepositories(inj)

	// Sender is stateless/connection-reusable - register as singleton.
	remy.RegisterSingleton(
		inj, func(retriever remy.DependencyRetriever) (mailer.Sender, error) {
			conf, _ := remy.Get[config.Config](retriever)
			return buildSender(conf.Mailer), nil
		},
	)

	// MailerImpl is a Factory: each request gets its own instance carrying the
	// request context so tracing and cancellation propagate through Send calls.
	remy.RegisterConstructorArgs2(
		inj, remy.Factory[*mailer.MailerImpl], mailer.NewMailerImpl,
	)
}

func buildSender(cfg config.Mailer) mailer.Sender {
	switch mailer.Driver(cfg.Driver) {
	case mailer.DriverSMTP:
		return smtp.New(
			smtp.Config{
				Host:        cfg.Smtp.Host,
				Port:        cfg.Smtp.Port,
				User:        cfg.Smtp.User,
				Password:    cfg.Smtp.Password,
				Encryption:  smtp.EncryptionTLS,
				FromAddress: cfg.From,
				FromName:    cfg.FromName,
			},
		)
	case mailer.DriverWebhook:
		return webhook.New(
			webhook.Config{
				URL:    cfg.Webhook.URL,
				APIKey: cfg.Webhook.APIKey,
			},
		)
	default:
		return stub.New()
	}
}
