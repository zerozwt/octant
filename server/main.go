package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zerozwt/octant/server/bridge"
	"github.com/zerozwt/octant/server/collector"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/handler"
	"github.com/zerozwt/swe"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	logger := swe.CtxLogger(nil)

	// load config
	confFile := ""
	flag.StringVar(&confFile, "conf", "", "config file")
	flag.Parse()
	if len(confFile) > 0 {
		if err := LoadConfig(confFile); err != nil {
			logger.Error("load config file %s failed: %v", confFile, err)
			return
		}
	} else {
		logger.Info("starting with default config ...")
	}

	// init log level & writer
	swe.SetDefaultLogLevel(gConfig.LogLevel())
	if len(gConfig.Log.File) != 0 {
		swe.SetDefaultLogWriter(&FileLogWriter{filename: gConfig.Log.File})
	}
	logger = swe.CtxLogger(nil)
	logger.Info("config %s loaded", confFile)

	// init db
	logger.Info("init db ...")
	if err := InitDB(); err != nil {
		logger.Error("init db failed: %v", err)
		return
	}

	// setting admin password
	if err := TrySetAdminPassword(); err != nil {
		logger.Error("init admin password failed: %v", err)
		return
	}

	// init collector bridge
	logger.Info("init collector bridge ...")
	collectorBridge, client, err := InitCollectorBridge()
	if err != nil {
		logger.Error("init collector bridge failed: %v", err)
		return
	}

	defer func() {
		collectorBridge.Stop()
		if client != nil {
			client.Close()
		}
	}()

	// init collector if needed
	if gConfig.Service.Collector {
		collectorBridge.SetReceiver(collector.GetCollector())
		logger.Info("data collector started")
	}
	defer collector.GetCollector().Stop()

	// init api engine if needed
	engine := InitEngine()
	if engine != nil {
		go engine.Serve(gConfig.WebAddr())
		logger.Info("web app service started on %s", gConfig.WebAddr())
	}

	// load streamers from db
	if err := LoadStreamers(); err != nil {
		logger.Error("load streamers from db failed: %v", err)
		return
	}

	chStop := make(chan bool, 1)
	go func() {
		defer close(chStop)
		tmp := make(chan os.Signal, 1)
		signal.Notify(tmp, syscall.SIGINT, syscall.SIGTERM)
		<-tmp
		logger.Info("octant shutting down ...")
		if engine != nil {
			engine.Close()
		}
	}()
	<-chStop
	logger.Info("octant shutdown")
}

func InitDB() error {
	logger := swe.CtxLogger(nil)
	if gConfig.IsMySQL() {
		if err := db.InitMySQL(gConfig.MySQL); err != nil {
			logger.Error("init mysql [%s] failed: %v", gConfig.MySQL, err)
			return err
		}
	} else if gConfig.IsSQLite() {
		if err := db.InitSQLite(gConfig.SQLite); err != nil {
			logger.Error("init sqlite [%s] failed: %v", gConfig.SQLite, err)
			return err
		}
	} else {
		logger.Error("init db failed: db type not specified")
		return errors.New("init db failed: db type not specified")
	}
	return nil
}

func InitCollectorBridge() (bridge.Bridge, *clientv3.Client, error) {
	var collectorBridge bridge.Bridge
	var client *clientv3.Client

	if gConfig.Service.Core && gConfig.Service.Collector {
		collectorBridge = bridge.CreateLocalBridge()
	} else {
		if len(gConfig.Etcd) == 0 {
			return nil, nil, errors.New("no etcd endpoints specified")
		}
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:            gConfig.Etcd,
			DialTimeout:          time.Second * 5,
			DialKeepAliveTime:    time.Second * 2,
			DialKeepAliveTimeout: time.Second * 5,
		})
		if err != nil {
			return nil, nil, err
		}
		client = cli
		collectorBridge = bridge.CreateEtcdBridge(client)
	}

	if err := collectorBridge.Start(); err != nil {
		return nil, nil, err
	}
	bridge.SetBridge(collectorBridge)

	return collectorBridge, client, nil
}

func TrySetAdminPassword() error {
	pass, err := db.GetSysConfigDAL().GetConfig(nil, db.DB_SYSCONF_ADMIN_PASS)
	if err != nil {
		return err
	}

	if len(pass) > 0 {
		return nil
	}

	swe.CtxLogger(nil).Info("admin password need to be set ...")

	fmt.Print("Enter admin password: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSuffix(line, "\n")

	fmt.Print("Enter admin password again: ")
	line2, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	line2 = strings.TrimSuffix(line2, "\n")

	if line != line2 {
		return errors.New("admin password not same")
	}

	line = db.GetSysConfigDAL().EncodeAdminPassword(line)

	return db.GetSysConfigDAL().SetConfig(nil, db.DB_SYSCONF_ADMIN_PASS, line)
}

func InitEngine() *swe.Engine {
	if !gConfig.Service.Core {
		return nil
	}

	api := swe.NewAPIServer()
	handler.InitAPIServer(api)

	file := swe.NewFileServer(gConfig.WebDir, "/", true)

	swe.CtxLogger(nil).Info("initializing web app service, root=%s", gConfig.WebDir)

	return swe.NewEngine(handler.API_PREFIX, api, file)
}

func LoadStreamers() error {
	streamers, err := db.GetStreamerDAL().All(nil)
	if err != nil {
		return err
	}

	br := bridge.GetBridge()
	logger := swe.CtxLogger(nil)

	for _, item := range streamers {
		logger.Info("start tracking live room %d", item.RoomID)
		if err := br.AddRoom(item.RoomID); err != nil {
			logger.Error("track live room %d failed: %v", item.RoomID, err)
			return err
		}
	}

	return nil
}
