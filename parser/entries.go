package parser

import "code.cloudfoundry.org/lager/chug"

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
