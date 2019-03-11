package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

type Hook struct {
	writer    io.Writer
	formatter log.Formatter
}

func (h Hook) Fire(e *log.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)
	return err
}

func (h Hook) Levels() []log.Level {
	return log.AllLevels
}

func main() {
	logFormatter := &log.JSONFormatter{
		FieldMap: *&log.FieldMap{
			log.FieldKeyMsg: "message",
		},
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(logFormatter)
	conn, err := net.Dial("tcp", "localhost:51401")
	if err != nil {

	}
	log.AddHook(Hook{
		writer:    conn,
		formatter: logFormatter,
	})

	log.WithFields(log.Fields{"project": "test", "module": "mode-test"}).Debug("This is debug. It has field.")

	log.Debug("This is debug log.HEIHEI")
	log.Info("This is info log.")
}
