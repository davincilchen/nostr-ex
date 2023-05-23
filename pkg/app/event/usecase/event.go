package usecase

import (
	repo "nostr-ex/pkg/app/event/repo/postgre"
	"nostr-ex/pkg/models"
)

type Handler struct {
}

func NewEventHandler() *Handler {
	return &Handler{}
}

func (t *Handler) SaveEvent(data models.Event) error {
	return repo.SaveEvent(data)
}

func (t *Handler) GetEvent(limit int) []models.Event {
	return repo.GetEvent(limit)
}
