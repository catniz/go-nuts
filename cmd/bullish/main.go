package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-nuts/internal"
	"go-nuts/pkg/config"
	"go-nuts/pkg/mysql"
	"go-nuts/pkg/route"
	"os"
	"os/signal"
	"runtime/debug"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	defer func() {
		if err := recover(); err != nil {
			log.WithError(err.(error)).WithField("stack", string(debug.Stack())).Fatal()
		}
	}()

	m := NewMain()
	m.logInit()

	defer func(m *Main) {
		_ = m.Close()
	}(m)

	if err := m.Run(ctx); err != nil {
		log.WithError(err.(error)).WithField("stack", string(debug.Stack())).Fatal()
	}
}

type Main struct {
	cfg config.Config
	db  *mysql.DB

	server *route.Server
	beans  *Beans
}

func NewMain() *Main {
	return &Main{}
}

func (m *Main) Close() error {
	if m.db != nil {
		if err := m.db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Main) Run(ctx context.Context) (err error) {
	m.cfg, err = config.LoadConfig("./")
	if err != nil {
		return err
	}

	m.db = mysql.NewDB(m.cfg.MySql)
	if err := m.db.Open(); err != nil {
		return err
	}

	m.beans = injectBeans()

	m.server = route.NewServer()
	m.registerRoutes()

	log.Info("Server is running on", m.cfg.HttpPort)
	if err := m.server.Router.Run(m.cfg.HttpPort); err != nil {
		return err
	}

	return nil
}

func (m *Main) logInit() {
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		DisableQuote:    true,
	})
}

type Beans struct {
	handlers []internal.Handler
}

func injectBeans() *Beans {
	return &Beans{
		handlers: []internal.Handler{},
	}
}

func (m *Main) registerRoutes() {
	for _, handler := range m.beans.handlers {
		handler.RegisterRoutes(m.server.Router)
	}
}
