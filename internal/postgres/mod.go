package postgres

import queries "irj/internal/postgres/_generated"

const (
	PendingPublicationStatus queries.PublicationStatus = "PENDING"
	DraftPublicationStatus   queries.PublicationStatus = "DRAFT"
)
