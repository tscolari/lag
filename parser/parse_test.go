package parser_test

import (
	"github.com/tscolari/lag/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Parse", func() {
	var (
		rawData       []byte
		rawDataStream *gbytes.Buffer
	)

	JustBeforeEach(func() {
		rawDataStream = gbytes.NewBuffer()
		_, err := rawDataStream.Write(rawData)
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		rawData = []byte(`{"timestamp":"1476218925.979571342","source":"guardian","message":"guardian.run.started","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179"}}
											{"timestamp":"1476218925.979685307","source":"guardian","message":"guardian.run.exec.start","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2"}}
											{"timestamp":"1476218925.979759455","source":"guardian","message":"guardian.run.exec.prepare.start","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.1"}}
											{"timestamp":"1476218927.086366415","source":"guardian","message":"guardian.run.exec.prepare.finished","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.1"}}
											{"timestamp":"1476218927.086794615","source":"guardian","message":"guardian.run.exec.execrunner.start","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2"}}
											{"timestamp":"1476218927.582846642","source":"guardian","message":"guardian.run.exec.execrunner.read-exit-fd","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2"}}
											{"timestamp":"1476218927.937837839","source":"guardian","message":"guardian.run.exec.execrunner.runc-exit-status","log_level":2,"data":{"error": "text error", "handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2","status":0}}
											{"timestamp":"1476218927.937957525","source":"guardian","message":"guardian.run.exec.execrunner.done","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2"}}
											{"timestamp":"1476218927.938017845","source":"guardian","message":"guardian.run.exec.finished","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2"}}
											{"timestamp":"1476218927.938051939","source":"guardian","message":"guardian.run.finished","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179"}}`)
	})

	It("parses the top level entries correctly", func() {
		entries, err := parser.Parse(rawDataStream)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(entries)).To(Equal(2))
		Expect(entries[0].Data.Message).To(Equal("guardian.run.started"))
		Expect(entries[1].Data.Message).To(Equal("guardian.run.finished"))
	})

	It("doesn't add a parent to top level messages", func() {
		entries, err := parser.Parse(rawDataStream)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(entries)).To(Equal(2))
		Expect(entries[0].Parent).To(BeNil())
		Expect(entries[1].Parent).To(BeNil())
	})

	It("correctly add the children logs", func() {
		entries, err := parser.Parse(rawDataStream)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(entries)).To(Equal(2))

		Expect(len(entries[0].Children)).To(Equal(2))
		Expect(len(entries[1].Children)).To(Equal(0))

		Expect(entries[0].Children[0].Data.Message).To(Equal("guardian.run.exec.start"))
		Expect(entries[0].Children[1].Data.Message).To(Equal("guardian.run.exec.finished"))
	})

	Context("nested children", func() {
		It("correctly parses children of children logs", func() {
			entries, err := parser.Parse(rawDataStream)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(entries)).To(Equal(2))

			Expect(len(entries[0].Children)).To(Equal(2))
			Expect(len(entries[0].Children[0].Children)).To(Equal(6))
			Expect(entries[0].Children[0].Children[0].Data.Message).To(Equal("guardian.run.exec.prepare.start"))
			Expect(entries[0].Children[0].Children[1].Data.Message).To(Equal("guardian.run.exec.prepare.finished"))
			Expect(entries[0].Children[0].Children[2].Data.Message).To(Equal("guardian.run.exec.execrunner.start"))
			Expect(entries[0].Children[0].Children[3].Data.Message).To(Equal("guardian.run.exec.execrunner.read-exit-fd"))
			Expect(entries[0].Children[0].Children[4].Data.Message).To(Equal("guardian.run.exec.execrunner.runc-exit-status"))
			Expect(entries[0].Children[0].Children[5].Data.Message).To(Equal("guardian.run.exec.execrunner.done"))
		})

		Context("complex nesting", func() {
			BeforeEach(func() {
				rawData = []byte(`
											{"timestamp":"1476218925.979571342","source":"g","message":"a","log_level":1,"data":{"session":"1"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab1","log_level":1,"data":{"session":"1.1"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab2","log_level":1,"data":{"session":"1.1"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab2c1","log_level":1,"data":{"session":"1.1.2"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab3","log_level":1,"data":{"session":"1.1"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab3c1","log_level":1,"data":{"session":"1.1.2"}}
											{"timestamp":"1476218925.979571342","source":"g","message":"ab3c2","log_level":1,"data":{"session":"1.1.2"}}
											`)
			})

			It("deals with it", func() {
				entries, err := parser.Parse(rawDataStream)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(entries)).To(Equal(1))
				Expect(entries[0].Data.Message).To(Equal("a"))

				Expect(len(entries[0].Children)).To(Equal(3))
				Expect(entries[0].Children[0].Data.Message).To(Equal("ab1"))
				Expect(entries[0].Children[1].Data.Message).To(Equal("ab2"))
				Expect(entries[0].Children[2].Data.Message).To(Equal("ab3"))

				Expect(len(entries[0].Children[0].Children)).To(Equal(0))

				Expect(len(entries[0].Children[1].Children)).To(Equal(1))
				Expect(entries[0].Children[1].Children[0].Data.Message).To(Equal("ab2c1"))

				Expect(len(entries[0].Children[2].Children)).To(Equal(2))
				Expect(entries[0].Children[2].Children[0].Data.Message).To(Equal("ab3c1"))
				Expect(entries[0].Children[2].Children[1].Data.Message).To(Equal("ab3c2"))
			})
		})
	})

	Context("when it's an error log", func() {
		BeforeEach(func() {
			rawData = []byte(`{"timestamp":"1476218927.937837839","source":"guardian","message":"guardian.run.exec.execrunner.runc-exit-status","log_level":2,"data":{"error": "text error", "handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2","status":0}}`)
		})

		It("marks it as errored", func() {
			entries, err := parser.Parse(rawDataStream)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(entries)).To(Equal(1))
			Expect(entries[0].Errored).To(BeTrue())
		})

		Context("when a children is an error log message", func() {
			BeforeEach(func() {
				rawData = []byte(`{"timestamp":"1476218925.979571342","source":"guardian","message":"guardian.run.started","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179"}}
											{"timestamp":"1476218925.979685307","source":"guardian","message":"guardian.run.exec.start","log_level":1,"data":{"handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2"}}
											{"timestamp":"1476218927.937837839","source":"guardian","message":"guardian.run.exec.execrunner.runc-exit-status","log_level":2,"data":{"error": "text error", "handle":"d1012287-83df-4456-5dcb-63d94b07a305","id":"d1012287-83df-4456-5dcb-63d94b07a305","path":"/tmp/lifecycle/healthcheck","session":"1172179.2.2","status":0}}`)
			})

			It("marks all the tree to the parent as errored", func() {
				entries, err := parser.Parse(rawDataStream)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(entries)).To(Equal(1))

				Expect(entries[0].Errored).To(BeTrue())

				Expect(len(entries[0].Children)).To(Equal(1))
				Expect(entries[0].Children[0].Errored).To(BeTrue())

				Expect(len(entries[0].Children[0].Children)).To(Equal(1))
				Expect(entries[0].Children[0].Children[0].Errored).To(BeTrue())
			})
		})
	})

})
