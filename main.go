package main

import (
	"flag"
	"github.com/sillyhatxu/mysql-client"
	log "github.com/sirupsen/logrus"
	"image-server/api"
	"image-server/config"
	"io"
	"net"
	"os"
	"time"
)

type DefaultFieldHook struct {
	GetValue func() string
}

func (h *DefaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *DefaultFieldHook) Fire(e *log.Entry) error {
	e.Data["project"] = h.GetValue()
	return nil
}

func GetValueImpl() string {
	return "image-server"
}

type LogstashHook struct {
	writer    io.Writer
	formatter log.Formatter
}

func (h LogstashHook) Fire(e *log.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)
	return err
}

func (h LogstashHook) Levels() []log.Level {
	return log.AllLevels
}

func init() {
	logFormatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		//TimestampFormat:string("2006-01-02 15:04:05"),
		FieldMap: *&log.FieldMap{
			log.FieldKeyMsg:  "message",
			log.FieldKeyTime: "@timestamp",
		},
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(logFormatter)
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}
	log.AddHook(&LogstashHook{writer: conn, formatter: logFormatter})
	log.AddHook(&DefaultFieldHook{GetValue: GetValueImpl})

	cfgFile := flag.String("c", "config.conf", "configuration file")
	flag.Parse()
	config.ParseConfig(*cfgFile)
}

func main() {
	dbclient.InitialDBClient(config.Conf.MysqlDB.DataSource, config.Conf.MysqlDB.MaxIdleConns, config.Conf.MysqlDB.MaxOpenConns)
	api.InitialAPI()
}
