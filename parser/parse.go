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

		if !AppendEntry(entries, newEntry) {
			entries = append(entries, &newEntry)
		}
	}

	return entries, nil
}

func AppendEntry(entries Entries, newEntry Entry) bool {
	for _, entry := range entries {
		if newEntry.Data.Session == entry.Data.Session {
			return false
		}

		if strings.HasPrefix(newEntry.Data.Session, entry.Data.Session) {
			if ok := AppendEntry(entry.Children, newEntry); ok {
				return true
			} else {
				newEntry.Parent = entry
				if newEntry.Errored {
					setErroredState(newEntry.Parent)
				}

				entry.Children = append(entry.Children, &newEntry)
				return true
			}
		}
	}

	return false
}

func setErroredState(entry *Entry) {
	if entry == nil {
		return
	}

	entry.Errored = true
	setErroredState(entry.Parent)
}
