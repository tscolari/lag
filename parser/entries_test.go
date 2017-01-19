package parser_test

import (
	"time"

	"code.cloudfoundry.org/lager/chug"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/lag/parser"
)

var _ = Describe("Entries", func() {
	var entries parser.Entries

	BeforeEach(func() {
		entries = parser.Entries{}
	})

	Describe("ErroredOnly", func() {
		BeforeEach(func() {
			entries = parser.Entries{
				&parser.Entry{Errored: true, Data: chug.LogEntry{Message: "1"}},
				&parser.Entry{Errored: false, Data: chug.LogEntry{Message: "2"}},
				&parser.Entry{Errored: true, Data: chug.LogEntry{Message: "3"}},
				&parser.Entry{Errored: false, Data: chug.LogEntry{Message: "4"}},
			}
		})

		It("filters the current list and returns only entries that contain an error", func() {
			filteredEntries := entries.ErroredOnly()
			Expect(len(filteredEntries)).To(Equal(2))

			Expect(filteredEntries[0].Data.Message).To(Equal("1"))
			Expect(filteredEntries[1].Data.Message).To(Equal("3"))
		})
	})

	Describe("SessionDuration", func() {
		BeforeEach(func() {
			timeNow := time.Now()
			entries = parser.Entries{
				&parser.Entry{Data: chug.LogEntry{Session: "2", Timestamp: timeNow.Add(-5 * time.Minute)}},
				&parser.Entry{Data: chug.LogEntry{Session: "1", Timestamp: timeNow.Add(-2 * time.Minute)}},
				&parser.Entry{Data: chug.LogEntry{Session: "2", Timestamp: timeNow}},
				&parser.Entry{Data: chug.LogEntry{Session: "3", Timestamp: timeNow.Add(10 * time.Second)}},
			}
		})

		It("returns the time difference between first and last entries of a session", func() {
			Expect(entries.SessionDuration("2")).To(Equal(5 * time.Minute))
		})

		Context("when the session has only one entry", func() {
			It("returns 0", func() {
				Expect(entries.SessionDuration("1")).To(BeZero())
			})
		})
	})

	Describe("RemoveSimilar", func() {
		BeforeEach(func() {
			entries = parser.Entries{
				&parser.Entry{Data: chug.LogEntry{Message: "Annoying log message"}},
				&parser.Entry{Data: chug.LogEntry{Message: "Annoying log message"}},
				&parser.Entry{Data: chug.LogEntry{Message: "Important log message"}},
				&parser.Entry{Data: chug.LogEntry{Message: "Annoying log message"}},
				&parser.Entry{Data: chug.LogEntry{Message: "Important log message"}},
			}
		})

		It("removes entries with the same message from the list", func() {
			filteredEntries := entries.RemoveSimilar(parser.Entry{Data: chug.LogEntry{Message: "Annoying log message"}})
			Expect(len(filteredEntries)).To(Equal(2))

			Expect(filteredEntries[0].Data.Message).To(Equal("Important log message"))
			Expect(filteredEntries[1].Data.Message).To(Equal("Important log message"))
		})
	})
})
