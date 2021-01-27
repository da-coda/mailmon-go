package store

import (
	"mailpie/pkg/event"
	"mailpie/pkg/instances"
)

const NewMailStoredEvent event.Event = "newMailStored"
const EventDispatcher = "MailStore"

type MailStore struct {
	mails        map[string]instances.Mail
	messageQueue event.Dispatcher
}

func CreateMailStore(messageQueue event.Dispatcher) *MailStore {
	var store *MailStore
	store = &MailStore{messageQueue: messageQueue}
	store.mails = make(map[string]instances.Mail)
	return store
}

func (store *MailStore) Add(key string, mailData instances.Mail) error {
	_, exists := store.mails[key]
	if exists {
		return AlreadyExistsError
	}

	return store.Set(key, mailData)
}

func (store *MailStore) Set(key string, data instances.Mail) error {
	store.mails[key] = data
	store.messageQueue.Dispatch(NewMailStoredEvent, EventDispatcher, data)
	return nil
}

func (store *MailStore) GetSingle(key string) (instances.Mail, error) {
	mail, exists := store.mails[key]
	if !exists {
		return instances.Mail{}, KeyNotExistsError
	}
	return mail, nil
}

func (store *MailStore) GetMultiple(keys []string) (mails map[string]instances.Mail, err error, notFoundKeys []string) {
	mails = make(map[string]instances.Mail)
	for _, key := range keys {
		mail, returnedErr := store.GetSingle(key)
		if returnedErr == KeyNotExistsError {
			err = KeyNotExistsError
			notFoundKeys = append(notFoundKeys, key)
			continue
		}
		mails[key] = mail
	}
	return
}
