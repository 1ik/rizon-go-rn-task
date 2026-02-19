package email

import (
	"context"
)

// EmailSender is the email sending contract used by the domain.
// Callers depend only on this interface; the implementation may be SMTP or anything else.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}
