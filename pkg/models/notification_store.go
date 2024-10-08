package models

import "sync"

type NotificationStore struct {
	Data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string, notification Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.Data[userID] = append(ns.Data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.Data[userID]
}
