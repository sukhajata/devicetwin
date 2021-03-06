package loggerhelper

import (
	"context"
	"fmt"
	pbLogger "github.com/sukhajata/pplogger"
	"time"
)

// WriteToLog writes to the console with formatted timestamp
func WriteToLog(msg interface{}) {
	fmt.Printf("[%s]: %v\n", time.Now().Format(time.RFC3339), msg)
}

type Helper interface {
	LogError(functionName string, message string, severity pbLogger.ErrorMessage_Severity)
}

type helper struct {
	grpcLoggerClient pbLogger.LoggerServiceClient
	serviceName      string
}

// NewHelper - returns a new helper
func NewHelper(grpcLoggerClient pbLogger.LoggerServiceClient, serviceName string) *helper {
	return &helper{
		grpcLoggerClient: grpcLoggerClient,
		serviceName:      serviceName,
	}
}

// LogError writes to console and sends to logger service
func (h *helper) LogError(functionName string, message string, severity pbLogger.ErrorMessage_Severity) {
	WriteToLog(functionName + ": " + message)

	if h.grpcLoggerClient != nil {
		request := &pbLogger.ErrorMessage{
			Service:  h.serviceName,
			Function: functionName,
			Severity: severity,
			Message:  message,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := h.grpcLoggerClient.LogError(ctx, request)
		if err != nil {
			WriteToLog(fmt.Sprintf("Failed to send error message: %v", err))
		}
	}
}
