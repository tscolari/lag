package parser

import (
	"io"
	"strings"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/chug"
)

func Parse(data io.Reader) (Entries, error) {
	entries := Entries{}
	lagerEntries := make(chan chug.Entry)
	go chug.Chug(data, lagerEntries)

	for entry := range lagerEntries {
		if !entry.IsLager {
			continue
		}

		newEntry := Entry{
			Data:     entry.Log,
			Children: []*Entry{},
			Errored:  (entry.Log.LogLevel == lager.ERROR),
		}

		if !appendEntry(entries, newEntry) {
			entries = append(entries, &newEntry)
		}
	}

	return entries, nil
}

func appendEntry(entries Entries, newEntry Entry) bool {
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if newEntry.Data.Session == entry.Data.Session {
			return false
		}

		if strings.HasPrefix(newEntry.Data.Session, entry.Data.Session) {
			if ok := appendEntry(entry.Children, newEntry); ok {
				return true
			} else {
				newEntry.Parent = entry
				if newEntry.Errored {
					newEntry.Errored = true
					propagateErroredState(&newEntry)
				}

				entry.Children = append(entry.Children, &newEntry)
				return true
			}
		}
	}

	return false
}

func propagateErroredState(entry *Entry) {
	for entry := entry.Parent; entry != nil; entry = entry.Parent {
		entry.Errored = true
	}
}
