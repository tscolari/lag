package parser

import (
	"time"

	"code.cloudfoundry.org/lager/chug"
)

type Entries []*Entry

type Entry struct {
	Data     chug.LogEntry
	Children Entries
	Parent   *Entry
	Errored  bool
}

func (es *Entries) ErroredOnly() Entries {
	erroredEntries := []*Entry{}

	for _, entry := range *es {
		if entry.Errored {
			erroredEntries = append(erroredEntries, entry)
		}
	}

	return erroredEntries
}

func (es *Entries) RemoveSimilar(sampleEntry Entry) Entries {
	filteredEntries := []*Entry{}

	for _, entry := range *es {
		if entry.Data.Message != sampleEntry.Data.Message {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries
}

func (es *Entries) SessionDuration(session string) time.Duration {
	var (
		firstEntry *Entry
		lastEntry  *Entry
	)

	for _, entry := range *es {
		if entry.Data.Session == session {
			if firstEntry == nil {
				firstEntry = entry
			} else {
				lastEntry = entry
			}
		}
	}

	if lastEntry != nil {
		return lastEntry.Data.Timestamp.Sub(firstEntry.Data.Timestamp)
	}

	return 0
}
