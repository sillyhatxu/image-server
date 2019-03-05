package main

import (
	"flag"
	"github.com/sillyhatxu/mysql-client"
	"image-server/api"
	"image-server/config"
)

func init() {
	cfgFile := flag.String("c", "config.conf", "configuration file")
	flag.Parse()
	config.ParseConfig(*cfgFile)
}

func main() {
	//log "github.com/sillyhatxu/microlog"
	//"github.com/sirupsen/logrus"
	//"net"

	//logFormatter := &logrus.JSONFormatter{
	//	FieldMap: *&logrus.FieldMap{
	//		logrus.FieldKeyMsg: "message",
	//	},
	//}
	//conn, err := net.Dial("tcp", "localhost:51401")
	//if err != nil {
	//	log.Fatal("net.Dial error.", err)
	//}
	//hook := log.New(conn, logFormatter)
	//logrusConfig := log.NewLogrusConfig(logFormatter, logrus.DebugLevel, logrus.Fields{"module":"image-server"}, true, hook)
	//err = logrusConfig.InstallConfig()
	//if err != nil {
	//	log.Fatal("logrus config initial error.", err)
	//}
	dbclient.InitialDBClient(config.Conf.MysqlDB.DataSource, config.Conf.MysqlDB.MaxIdleConns, config.Conf.MysqlDB.MaxOpenConns)
	api.InitialAPI()
}
