package errorhelper

import (
	"github.com/sukhajata/devicetwin/pkg/loggerhelper"
	"time"
)

// PanicOnError logs and panics if err is not nil
func PanicOnError(err error) {
	if err != nil {
		loggerhelper.WriteToLog(err.Error())
		panic(err)
	}
}

// StartUpError writes error to log and waits 2 seconds
func StartUpError(err error) {
	loggerhelper.WriteToLog(err)
	loggerhelper.WriteToLog("Waiting 2 seconds to retry...")
	time.Sleep(2 * time.Second)
}
