package log_streamer_test

import (
	"fmt"
	"strings"
	"sync"
	"time"

	mfakes "code.cloudfoundry.org/diego-logging-client/testhelpers"
	"code.cloudfoundry.org/executor/depot/log_streamer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("LogStreamer", func() {
	var (
		streamer                           log_streamer.LogStreamer
		fakeClient                         *mfakes.FakeIngressClient
		maxLogLinesPerSecond               int
		logRateLimitExceededReportInterval time.Duration
	)

	guid := "the-guid"
	sourceName := "the-source-name"
	index := 11
	tags := map[string]string{
		"foo": "bar",
		"biz": "baz",
	}

	BeforeEach(func() {
		maxLogLinesPerSecond = 9999
		logRateLimitExceededReportInterval = 5 * time.Minute
		fakeClient = &mfakes.FakeIngressClient{}
		streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)
	})

	Context("when told to emit", func() {
		Context("when given a message that corresponds to one line", func() {
			BeforeEach(func() {
				fmt.Fprintln(streamer.Stdout(), "this is a log")
				fmt.Fprintln(streamer.Stdout(), "this is another log")
			})

			It("should emit that message", func() {
				Expect(fakeClient.SendAppLogCallCount()).To(Equal(2))

				message, sn, tags := fakeClient.SendAppLogArgsForCall(0)
				Expect(tags["source_id"]).To(Equal(guid))
				Expect(sn).To(Equal(sourceName))
				Expect(message).To(Equal("this is a log"))
				Expect(tags["instance_id"]).To(Equal("11"))
				Expect(tags["foo"]).To(Equal("bar"))
				Expect(tags["biz"]).To(Equal("baz"))

				message, sn, tags = fakeClient.SendAppLogArgsForCall(1)
				Expect(tags["source_id"]).To(Equal(guid))
				Expect(sn).To(Equal(sourceName))
				Expect(message).To(Equal("this is another log"))
				Expect(tags["instance_id"]).To(Equal("11"))
				Expect(tags["foo"]).To(Equal("bar"))
				Expect(tags["biz"]).To(Equal("baz"))
			})
		})

		Context("when given a log line arrival rate exceeds the maximum allowed", func() {
			Context("rate limit is applied at a lower threshold", func() {
				BeforeEach(func() {
					maxLogLinesPerSecond = 1
					streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

					for i := 0; i < maxLogLinesPerSecond*3; i++ {
						go fmt.Fprintf(streamer.Stdout(), "this is log # %d\n", i)
					}
				})

				It("should rate limit the messages", func() {
					Eventually(func() int {
						printedMsgs := 0
						for i := 0; i < fakeClient.SendAppLogCallCount(); i++ {
							msg, _, _ := fakeClient.SendAppLogArgsForCall(i)
							if strings.Contains(msg, "this is log") {
								printedMsgs++
							}
						}
						return printedMsgs
					}, 10*time.Second).Should(Equal(3))

					calls := fakeClient.SendAppLogCallCount()
					args := []string{}
					for i := 0; i < calls; i++ {
						msg, _, _ := fakeClient.SendAppLogArgsForCall(i)
						args = append(args, msg)
					}

					expectedArgs := []types.GomegaMatcher{}
					for i := 0; i < 3; i++ {
						expectedArgs = append(expectedArgs, ContainSubstring("this is log"))
					}
					for i := 0; i < calls-3; i++ {
						expectedArgs = append(expectedArgs, ContainSubstring("app instance exceeded log rate limit (1 log-lines/sec) set by platform operator"))
					}

					Expect(args).To(ConsistOf(expectedArgs))
				})

				It("increments an AppInstanceExceededLogRateLimitCount metric", func() {
					Eventually(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(Equal(1))
					metricName := fakeClient.IncrementCounterArgsForCall(0)
					Expect(metricName).To(Equal("AppInstanceExceededLogRateLimitCount"))
				})

				Context("when metric was already incremented in the past metric report interval", func() {
					It("does not increment an AppInstanceExceededLogRateLimitCount metric during report interval", func() {
						Eventually(fakeClient.SendAppLogCallCount, 3*time.Second).Should(BeNumerically(">=", 4))
						Consistently(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(Equal(1))
					})
				})
			})

			Context("when metric was incremented and report interval passed", func() {
				BeforeEach(func() {
					maxLogLinesPerSecond = 1
					logRateLimitExceededReportInterval = time.Second
					streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

					for i := 0; i < 3; i++ {
						go fmt.Fprintf(streamer.Stdout(), "this is log # %d \n", i)
					}
				})

				It("increments an AppInstanceExceededLogRateLimitCount metric during report interval", func() {
					Consistently(fakeClient.SendAppLogCallCount, 1*time.Second).Should(BeNumerically("<", 4))
					Eventually(fakeClient.SendAppLogCallCount, 5*time.Second).Should(BeNumerically(">=", 4))
					Eventually(fakeClient.IncrementCounterCallCount, 5*time.Second).Should(Equal(2))
				})
			})

			Context("rate limit is not applied", func() {
				BeforeEach(func() {
					maxLogLinesPerSecond = 0
					streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

					for i := 0; i < 20; i++ {
						go fmt.Fprintf(streamer.Stdout(), "this is log # %d \n", i)
					}
				})

				It("should not rate limit the messages", func() {
					Eventually(fakeClient.SendAppLogCallCount, time.Second).Should(Equal(20))
				})
			})

			Context("rate limit is bigger than number of log lines", func() {
				BeforeEach(func() {
					maxLogLinesPerSecond = 6
					streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

					for i := 0; i < 3; i++ {
						go fmt.Fprintf(streamer.Stdout(), "this is log # %d \n", i)
					}
				})

				It("should not rate limit the messages", func() {
					Eventually(fakeClient.SendAppLogCallCount, time.Second).Should(Equal(3))
				})
			})
		})

		Describe("WithSource", func() {
			Context("when a new log source is provided", func() {
				var newSourceName string

				BeforeEach(func() {
					newSourceName = "new-source-name"
					streamer = streamer.WithSource(newSourceName)
				})

				It("should emit a message with the new log source", func() {
					fmt.Fprintln(streamer.Stdout(), "this is a log")
					Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))

					_, sn, _ := fakeClient.SendAppLogArgsForCall(0)

					Expect(sn).To(Equal(newSourceName))
				})
			})

			Context("when no log source is provided", func() {
				BeforeEach(func() {
					streamer = streamer.WithSource("")
				})

				It("should emit a message with the existing log source", func() {
					fmt.Fprintln(streamer.Stdout(), "this is a log")

					Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))

					_, sn, _ := fakeClient.SendAppLogArgsForCall(0)

					Expect(sn).To(Equal(sourceName))
				})
			})

			Context("rate limit is applied at a lower threshold", func() {
				var newStreamer log_streamer.LogStreamer

				BeforeEach(func() {
					maxLogLinesPerSecond = 1
					streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

					newStreamer = streamer.WithSource("new-source-name")
				})

				It("should rate limit the messages for original and new log streamers", func() {
					for i := 0; i < maxLogLinesPerSecond*3; i++ {
						go fmt.Fprintf(streamer.Stdout(), "old streamer: this is log # %d \n", i)
						go fmt.Fprintf(newStreamer.Stdout(), "new streamer: this is log # %d \n", i)
					}

					Eventually(func() int {
						printedMsgs := 0
						for i := 0; i < fakeClient.SendAppLogCallCount(); i++ {
							msg, _, _ := fakeClient.SendAppLogArgsForCall(i)
							if strings.Contains(msg, "this is log") {
								printedMsgs++
							}
						}
						return printedMsgs
					}, 10*time.Second).Should(Equal(6))

					calls := fakeClient.SendAppLogCallCount()
					args := []string{}
					for i := 0; i < calls; i++ {
						msg, _, _ := fakeClient.SendAppLogArgsForCall(i)
						args = append(args, msg)
					}

					expectedArgs := []types.GomegaMatcher{}
					for i := 0; i < 6; i++ {
						expectedArgs = append(expectedArgs, ContainSubstring("this is log"))
					}
					for i := 0; i < calls-6; i++ {
						expectedArgs = append(expectedArgs, ContainSubstring("app instance exceeded log rate limit (1 log-lines/sec) set by platform operator"))
					}

					Expect(args).To(ConsistOf(expectedArgs))
				})

				It("increments an AppInstanceExceededLogRateLimitCount metric once per interval for both streamers", func() {
					for i := 0; i < maxLogLinesPerSecond*3; i++ {
						go fmt.Fprintf(streamer.Stdout(), "old streamer: this is log # %d \n", i)
						go fmt.Fprintf(newStreamer.Stdout(), "new streamer: this is log # %d \n", i)
					}

					Eventually(fakeClient.IncrementCounterCallCount, 2*time.Second).Should(Equal(1))
					metricName := fakeClient.IncrementCounterArgsForCall(0)
					Expect(metricName).To(Equal("AppInstanceExceededLogRateLimitCount"))
				})

				Context("when metric was already incremented in the past metric report interval", func() {
					It("does not increment an AppInstanceExceededLogRateLimitCount metric", func() {
						for i := 0; i < maxLogLinesPerSecond*3; i++ {
							go fmt.Fprintf(streamer.Stdout(), "old streamer: this is log # %d \n", i)
						}
						// can either print the last warning or not if not fast
						Eventually(fakeClient.SendAppLogCallCount, 3*time.Second).Should(BeNumerically(">=", 4))
						Consistently(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(Equal(1))

						for i := 0; i < maxLogLinesPerSecond*3; i++ {
							go fmt.Fprintf(newStreamer.Stdout(), "new streamer: this is log # %d \n", i)
						}
						Eventually(fakeClient.SendAppLogCallCount, 3*time.Second).Should(BeNumerically(">=", 7))
						Consistently(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(Equal(1))
					})
				})

				Context("when metric was incremented and report interval passed", func() {
					BeforeEach(func() {
						maxLogLinesPerSecond = 1
						logRateLimitExceededReportInterval = time.Second
						streamer = log_streamer.New(guid, sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)
						newStreamer = streamer.WithSource("new-source-name")
					})

					It("increments an AppInstanceExceededLogRateLimitCount metric", func() {
						for i := 0; i < maxLogLinesPerSecond*3; i++ {
							go fmt.Fprintf(streamer.Stdout(), "old streamer: this is log # %d \n", i)
						}
						Eventually(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(BeNumerically(">=", 1))

						for i := 0; i < maxLogLinesPerSecond*3; i++ {
							go fmt.Fprintf(newStreamer.Stdout(), "new streamer: this is log # %d \n", i)
						}
						Eventually(fakeClient.SendAppLogCallCount, 10*time.Second).Should(BeNumerically(">=", 5))
						Eventually(fakeClient.IncrementCounterCallCount, 3*time.Second).Should(BeNumerically(">=", 3))
					})
				})
			})
		})

		Describe("SourceName", func() {
			It("should return the log streamer's configured source name", func() {
				Expect(streamer.SourceName()).To(Equal(sourceName))
			})
		})

		Context("when given a message with all sorts of fun newline characters", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "A\nB\rC\n\rD\r\nE\n\n\nF\r\r\rG\n\r\r\n\n\n\r")
			})

			It("should do the right thing", func() {
				Expect(fakeClient.SendAppLogCallCount()).To(Equal(7))
				for i, expectedString := range []string{"A", "B", "C", "D", "E", "F", "G"} {
					message, _, _ := fakeClient.SendAppLogArgsForCall(i)
					Expect(message).To(Equal(expectedString))
				}
			})
		})

		Context("when given a series of short messages", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "this is a log")
				fmt.Fprintf(streamer.Stdout(), " it is made of wood")
				fmt.Fprintf(streamer.Stdout(), " - and it is longer")
				fmt.Fprintf(streamer.Stdout(), "than it seems\n")
			})

			It("concatenates them, until a new-line is received, and then emits that", func() {
				Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))
				message, _, _ := fakeClient.SendAppLogArgsForCall(0)
				Expect(message).To(Equal("this is a log it is made of wood - and it is longerthan it seems"))
			})
		})

		Context("when given a message with multiple new lines", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "this is a log\nand this is another\nand this one isn't done yet...")
			})

			It("should break the message up into multiple loggings", func() {
				Expect(fakeClient.SendAppLogCallCount()).To(Equal(2))

				message, _, _ := fakeClient.SendAppLogArgsForCall(0)
				Expect(message).To(Equal("this is a log"))

				message, _, _ = fakeClient.SendAppLogArgsForCall(1)
				Expect(message).To(Equal("and this is another"))
			})
		})

		Describe("message limits", func() {
			var message string
			Context("when the message is just at the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", log_streamer.MAX_MESSAGE_SIZE)
					Expect([]byte(message)).To(HaveLen(log_streamer.MAX_MESSAGE_SIZE), "Ensure that the byte representation of our message is under the limit")

					fmt.Fprintf(streamer.Stdout(), message+"\n")
				})

				It("should not break the message up and send a single messages", func() {
					Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))
					ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
					Expect(ms).To(Equal(message))
				})
			})

			Context("when the message exceeds the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", log_streamer.MAX_MESSAGE_SIZE)
					message += strings.Repeat("8", log_streamer.MAX_MESSAGE_SIZE)
					message += strings.Repeat("9", log_streamer.MAX_MESSAGE_SIZE)
					message += "hello\n"
					fmt.Fprintf(streamer.Stdout(), message)
				})

				It("should break the message up and send multiple messages", func() {
					Expect(fakeClient.SendAppLogCallCount()).To(Equal(4))

					ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
					Expect(ms).To(Equal(strings.Repeat("7", log_streamer.MAX_MESSAGE_SIZE)))
					ms, _, _ = fakeClient.SendAppLogArgsForCall(1)
					Expect(ms).To(Equal(strings.Repeat("8", log_streamer.MAX_MESSAGE_SIZE)))
					ms, _, _ = fakeClient.SendAppLogArgsForCall(2)
					Expect(ms).To(Equal(strings.Repeat("9", log_streamer.MAX_MESSAGE_SIZE)))
					ms, _, _ = fakeClient.SendAppLogArgsForCall(3)
					Expect(ms).To(Equal("hello"))
				})
			})

			Context("when having to deal with byte boundaries and long utf characters", func() {
				BeforeEach(func() {
					message = strings.Repeat("a", log_streamer.MAX_MESSAGE_SIZE-3)
					message += "\U0001F428\n"
				})

				It("should break the message up and send multiple messages without sending error runes", func() {
					fmt.Fprintf(streamer.Stdout(), message)
					Expect(fakeClient.SendAppLogCallCount()).To(Equal(2))

					ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
					Expect(ms).To(Equal(strings.Repeat("a", log_streamer.MAX_MESSAGE_SIZE-3)))
					ms, _, _ = fakeClient.SendAppLogArgsForCall(1)
					Expect(ms).To(Equal("\U0001F428"))
				})

				Context("with an invalid utf8 character in the message", func() {
					var utfChar string

					BeforeEach(func() {
						message = strings.Repeat("9", log_streamer.MAX_MESSAGE_SIZE-4)
						utfChar = "\U0001F428"
					})

					It("emits both messages correctly", func() {
						fmt.Fprintf(streamer.Stdout(), message+utfChar[0:2])
						fmt.Fprintf(streamer.Stdout(), utfChar+"\n")

						Expect(fakeClient.SendAppLogCallCount()).To(Equal(2))

						ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
						Expect(ms).To(Equal(message + utfChar[0:2]))

						ms, _, _ = fakeClient.SendAppLogArgsForCall(1)
						Expect(ms).To(Equal(utfChar))
					})
				})

				Context("when the entire message is invalid utf8 characters", func() {
					var utfChar string

					BeforeEach(func() {
						utfChar = "\U0001F428"
						message = strings.Repeat(utfChar[0:2], log_streamer.MAX_MESSAGE_SIZE/2)
						Expect(len(message)).To(Equal(log_streamer.MAX_MESSAGE_SIZE))
					})

					It("drops the last 3 bytes", func() {
						fmt.Fprintf(streamer.Stdout(), message)

						Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))

						ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
						Expect(ms).To(Equal(message[0 : len(message)-3]))
					})
				})
			})

			Context("while concatenating, if the message exceeds the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", log_streamer.MAX_MESSAGE_SIZE-2)
					fmt.Fprintf(streamer.Stdout(), message)
					fmt.Fprintf(streamer.Stdout(), "778888\n")
				})

				It("should break the message up and send multiple messages", func() {
					Expect(fakeClient.SendAppLogCallCount()).To(Equal(2))

					ms, _, _ := fakeClient.SendAppLogArgsForCall(0)
					Expect(ms).To(Equal(strings.Repeat("7", log_streamer.MAX_MESSAGE_SIZE)))
					ms, _, _ = fakeClient.SendAppLogArgsForCall(1)
					Expect(ms).To(Equal("8888"))
				})
			})
		})
	})

	Context("when told to emit stderr", func() {
		It("should handle short messages", func() {
			fmt.Fprintf(streamer.Stderr(), "this is a log\nand this is another\nand this one isn't done yet...")
			Expect(fakeClient.SendAppErrorLogCallCount()).To(Equal(2))

			msg, sn, _ := fakeClient.SendAppErrorLogArgsForCall(0)
			Expect(msg).To(Equal("this is a log"))
			Expect(sn).To(Equal(sourceName))

			msg, sn, _ = fakeClient.SendAppErrorLogArgsForCall(1)
			Expect(msg).To(Equal("and this is another"))
			Expect(sn).To(Equal(sourceName))
		})

		It("should handle long messages", func() {
			fmt.Fprintf(streamer.Stderr(), strings.Repeat("e", log_streamer.MAX_MESSAGE_SIZE+1)+"\n")
			Expect(fakeClient.SendAppErrorLogCallCount()).To(Equal(2))

			msg, _, _ := fakeClient.SendAppErrorLogArgsForCall(0)
			Expect(msg).To(Equal(strings.Repeat("e", log_streamer.MAX_MESSAGE_SIZE)))

			msg, _, _ = fakeClient.SendAppErrorLogArgsForCall(1)
			Expect(msg).To(Equal("e"))
		})
	})

	Context("when told to flush", func() {
		It("should send whatever log is left in its buffer", func() {
			fmt.Fprintf(streamer.Stdout(), "this is a stdout")
			fmt.Fprintf(streamer.Stderr(), "this is a stderr")

			Expect(fakeClient.SendAppLogCallCount()).To(Equal(0))
			Expect(fakeClient.SendAppErrorLogCallCount()).To(Equal(0))

			streamer.Flush()

			Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))
			Expect(fakeClient.SendAppErrorLogCallCount()).To(Equal(1))
		})
	})

	Context("when there is no app guid", func() {
		It("does nothing when told to emit or flush", func() {
			streamer = log_streamer.New("", sourceName, index, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

			streamer.Stdout().Write([]byte("hi"))
			streamer.Stderr().Write([]byte("hi"))
			streamer.Flush()

			Expect(fakeClient.SendAppLogCallCount()).To(Equal(0))
		})
	})

	Context("when there is no log source", func() {
		It("defaults to LOG", func() {
			streamer = log_streamer.New(guid, "", -1, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

			streamer.Stdout().Write([]byte("hi"))
			streamer.Flush()

			Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))
			_, sn, _ := fakeClient.SendAppLogArgsForCall(0)
			Expect(sn).To(Equal(log_streamer.DefaultLogSource))
		})
	})

	Context("when there is no source index", func() {
		It("defaults to 0", func() {
			streamer = log_streamer.New(guid, sourceName, -1, tags, fakeClient, maxLogLinesPerSecond, logRateLimitExceededReportInterval)

			streamer.Stdout().Write([]byte("hi"))
			streamer.Flush()

			Expect(fakeClient.SendAppLogCallCount()).To(Equal(1))
			_, _, tags := fakeClient.SendAppLogArgsForCall(0)
			Expect(tags["instance_id"]).To(Equal("-1"))
		})
	})

	Context("with multiple goroutines emitting simultaneously", func() {
		var waitGroup *sync.WaitGroup

		BeforeEach(func() {
			waitGroup = new(sync.WaitGroup)

			for i := 0; i < 2; i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					fmt.Fprintln(streamer.Stdout(), "this is a log")
				}()
			}
		})

		AfterEach(func(done Done) {
			defer close(done)
			waitGroup.Wait()
		})

		It("does not trigger data races", func() {
			Eventually(fakeClient.SendAppLogCallCount).Should(Equal(2))
		})
	})

	Describe("Stop", func() {
		Context("stopping the log streamer", func() {
			BeforeEach(func() {
				streamer.Stop()
			})

			It("writes to stdout and stderr should not fail", func() {
				_, stdOutErr := fmt.Fprintln(streamer.Stdout(), "this is a log")
				Expect(stdOutErr).NotTo(HaveOccurred())
				_, stdErrErr := fmt.Fprintln(streamer.Stderr(), "this is another log")
				Expect(stdErrErr).NotTo(HaveOccurred())

			})

			Context("when a 'child' log streamer is present", func() {
				var childStreamer log_streamer.LogStreamer

				BeforeEach(func() {
					childStreamer = streamer.WithSource("CHILD")
				})

				It("writes to the child's stdout and stderr should not fail", func() {
					_, stdOutErr := fmt.Fprintln(childStreamer.Stdout(), "this is a log")
					Expect(stdOutErr).NotTo(HaveOccurred())
					_, stdErrErr := fmt.Fprintln(childStreamer.Stderr(), "this is another log")
					Expect(stdErrErr).NotTo(HaveOccurred())
				})
			})
		})

		Context("when a 'child' log streamer is stopped", func() {
			var childStreamer log_streamer.LogStreamer

			BeforeEach(func() {
				childStreamer = streamer.WithSource("CHILD")
				childStreamer.Stop()
			})

			It("writes to the child's stdout and stderr should not fail", func() {
				_, stdOutErr := fmt.Fprintln(childStreamer.Stdout(), "this is a log")
				Expect(stdOutErr).NotTo(HaveOccurred())
				_, stdErrErr := fmt.Fprintln(childStreamer.Stderr(), "this is another log")
				Expect(stdErrErr).NotTo(HaveOccurred())
			})

			It("writes to the parent's stdout and stderr should continue to succeed", func() {
				_, stdOutErr := fmt.Fprintln(streamer.Stdout(), "this is a log")
				Expect(stdOutErr).NotTo(HaveOccurred())
				_, stdErrErr := fmt.Fprintln(streamer.Stderr(), "this is another log")
				Expect(stdErrErr).NotTo(HaveOccurred())
			})
		})
	})
})
