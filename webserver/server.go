package webserver

import (
  "context"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  _ "github.com/go-macaron/session/postgres"
  "github.com/plumbie/plumbie/config"
  "github.com/plumbie/plumbie/models"
  apiv1 "github.com/plumbie/plumbie/webserver/api/v1"

  "github.com/go-macaron/session"
  log "github.com/sirupsen/logrus"
  "gopkg.in/macaron.v1"
)

func getHTTPServer() (s *http.Server) {
  m := macaron.New()

  m.SetAutoHead(true)
  m.Use(macaron.Logger())
  m.Use(macaron.Recovery())
  m.Use(macaron.Static("public"))
  m.Use(macaron.Renderer())
  m.Use(session.Sessioner(session.Options{
    // Provider:   "memory",
    Provider:       models.Driver,
    ProviderConfig: models.ConnectionStr,
    CookieName:     "Session",
  }))

  m.Group("/api", func() {
    apiv1.RegisterRoutes(m)
  })

  return &http.Server{
    Handler:      m,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    Addr:         fmt.Sprintf("%s:%d", config.Web.BindIP, config.Web.Port),
  }
}

func handleKillSignals(server *http.Server, idleConnsClosed chan struct{}) {
  sigint := make(chan os.Signal, 1)

  signal.Notify(sigint, syscall.SIGINT)
  signal.Notify(sigint, syscall.SIGTERM)

  <-sigint
  log.Debug("webserver: Crtl-C detected")
  signal.Reset(syscall.SIGINT)
  signal.Reset(syscall.SIGTERM)
  if err := server.Shutdown(context.Background()); err != nil {
    // Error from closing listeners, or context timeout:
    log.Errorf("HTTP server Shutdown: %v", err)
  }
  log.Debug("webserver: Shutdown complete")
  close(idleConnsClosed)
}

func Start() error {
  s := getHTTPServer()

  idleConnsClosed := make(chan struct{})
  go handleKillSignals(s, idleConnsClosed)

  log.Infof("webserver: Starting web server on %s:%d", config.Web.BindIP, config.Web.Port)
  if err := s.ListenAndServe(); err != http.ErrServerClosed {
    log.Errorf("webserver: HTTP server ListenAndServe: %v", err)
    close(idleConnsClosed)
    return err
  }
  <-idleConnsClosed
  return nil
}
