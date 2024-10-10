package liblogger

import (
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {
	// f, err := os.OpenFile(fmt.Sprintf("keyswitch-%s.log", time.Now().Format("02-01-2006-15-04-05")), os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	panic(fmt.Errorf("error opening output log file. error = %v", err))
	// }
	Log = logrus.New()
	// Log.SetFormatter(&logrus.JSONFormatter{})
	//Log.SetOutput(f)
}
