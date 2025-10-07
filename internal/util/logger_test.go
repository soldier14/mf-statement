package util_test

import (
	"bytes"
	"mf-statement/internal/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	var (
		logger *util.Logger
		buf    *bytes.Buffer
	)

	BeforeEach(func() {
		buf = &bytes.Buffer{}
	})

	Describe("NewLogger", func() {
		It("should create a logger with specified level and output", func() {
			logger = util.NewLogger(util.LogLevelDebug, buf)
			Expect(logger).ToNot(BeNil())
		})

		It("should use stderr when output is nil", func() {
			logger = util.NewLogger(util.LogLevelInfo, nil)
			Expect(logger).ToNot(BeNil())
		})
	})

	Describe("NewDefaultLogger", func() {
		It("should create a logger with default settings", func() {
			logger = util.NewDefaultLogger()
			Expect(logger).ToNot(BeNil())
		})
	})

	Describe("Logging methods", func() {
		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelDebug, buf)
		})

		It("should log debug messages when level allows", func() {
			logger.Debug("debug message")
			output := buf.String()
			Expect(output).To(ContainSubstring("DEBUG"))
			Expect(output).To(ContainSubstring("debug message"))
		})

		It("should log info messages when level allows", func() {
			logger.Info("info message")
			output := buf.String()
			Expect(output).To(ContainSubstring("INFO"))
			Expect(output).To(ContainSubstring("info message"))
		})

		It("should log warn messages when level allows", func() {
			logger.Warn("warn message")
			output := buf.String()
			Expect(output).To(ContainSubstring("WARN"))
			Expect(output).To(ContainSubstring("warn message"))
		})

		It("should log error messages when level allows", func() {
			logger.Error("error message")
			output := buf.String()
			Expect(output).To(ContainSubstring("ERROR"))
			Expect(output).To(ContainSubstring("error message"))
		})

		It("should format messages with arguments", func() {
			logger.Info("User %s has %d items", "john", 5)
			output := buf.String()
			Expect(output).To(ContainSubstring("User john has 5 items"))
		})

		It("should include timestamp in log messages", func() {
			logger.Info("test message")
			output := buf.String()
			Expect(output).To(MatchRegexp(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\]`))
		})
	})

	Describe("Log level filtering", func() {
		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelWarn, buf)
		})

		It("should not log debug messages when level is warn", func() {
			logger.Debug("debug message")
			Expect(buf.String()).To(BeEmpty())
		})

		It("should not log info messages when level is warn", func() {
			logger.Info("info message")
			Expect(buf.String()).To(BeEmpty())
		})

		It("should log warn messages when level is warn", func() {
			logger.Warn("warn message")
			output := buf.String()
			Expect(output).To(ContainSubstring("WARN"))
			Expect(output).To(ContainSubstring("warn message"))
		})

		It("should log error messages when level is warn", func() {
			logger.Error("error message")
			output := buf.String()
			Expect(output).To(ContainSubstring("ERROR"))
			Expect(output).To(ContainSubstring("error message"))
		})
	})

	Describe("WithFields", func() {
		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelInfo, buf)
		})

		It("should return the same logger", func() {
			fields := map[string]interface{}{
				"user": "john",
				"id":   123,
			}
			newLogger := logger.WithFields(fields)
			Expect(newLogger).To(Equal(logger))
		})
	})

	Describe("SetLevel", func() {
		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelInfo, buf)
		})

		It("should change the log level", func() {
			logger.SetLevel(util.LogLevelError)
			logger.Info("info message")
			Expect(buf.String()).To(BeEmpty())
		})
	})

	Describe("SetOutput", func() {
		var newBuf *bytes.Buffer

		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelInfo, buf)
			newBuf = &bytes.Buffer{}
		})

		It("should change the output writer", func() {
			logger.SetOutput(newBuf)
			logger.Info("test message")
			Expect(buf.String()).To(BeEmpty())
			Expect(newBuf.String()).To(ContainSubstring("test message"))
		})
	})

	Describe("Log level constants", func() {
		It("should have correct log level values", func() {
			Expect(util.LogLevelDebug).To(Equal(util.LogLevel(0)))
			Expect(util.LogLevelInfo).To(Equal(util.LogLevel(1)))
			Expect(util.LogLevelWarn).To(Equal(util.LogLevel(2)))
			Expect(util.LogLevelError).To(Equal(util.LogLevel(3)))
		})
	})

	Describe("Log formatting", func() {
		BeforeEach(func() {
			logger = util.NewLogger(util.LogLevelInfo, buf)
		})

		It("should handle empty format string", func() {
			logger.Info("")
			output := buf.String()
			Expect(output).To(ContainSubstring("INFO"))
		})

		It("should handle format string with no arguments", func() {
			logger.Info("simple message")
			output := buf.String()
			Expect(output).To(ContainSubstring("simple message"))
		})

		It("should handle multiple arguments", func() {
			logger.Info("User %s has %d items and balance %.2f", "john", 5, 123.45)
			output := buf.String()
			Expect(output).To(ContainSubstring("User john has 5 items and balance 123.45"))
		})
	})
})
