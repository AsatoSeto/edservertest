package server

import (
	"edServer/sse"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var StatusBoard = make(map[string]Status)
var Stat = make(chan []byte)

func RouteStart(broker *sse.Broker) *echo.Echo {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("./htmltemplate/*.html")),
	}
	e.Renderer = t
	e.GET("/dashboard", Dashboard)
	e.GET("/eventTest", func(c echo.Context) error {
		go func() {
			if err := sendStatus(); err != nil {
				log.Println("sendStatus error:", err)
			}
		}()
		broker.ServeHTTP(c.Response(), c.Request())

		return nil
	})
	e.POST("/status", StatusHandler)
	e.POST("/shutdown", ShutdownHandler)
	return e
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.Execute(w, Context{})
}

func Dashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", nil)
}

func StatusHandler(ctx echo.Context) error {

	s := Status{}
	err := ParseResponseBody(ctx.Request().Body, &s)
	if err != nil {
		log.Println("StatusHandler parse response body error:", err)
		return err
	}
	StatusBoard[s.PlayerID] = s

	if err = sendStatus(); err != nil {
		log.Println("StatusHandler sendStatus error:", err)
		return err
	}
	return nil
}

func ShutdownHandler(ctx echo.Context) error {
	s := Shutdown{}
	err := ParseResponseBody(ctx.Request().Body, &s)
	if err != nil {
		log.Println("ShutdownHandler parse response body error:", err)
		return err
	}
	s.Delete = true
	body, err := json.Marshal(s)
	if err != nil {
		log.Println("ShutdownHandler marshal error:", err)
		return err
	}
	Stat <- body
	return nil
}

func sendStatus() error {
	body, err := json.Marshal(StatusBoard)
	if err != nil {
		log.Println("sendStatus marshal error:", err)
		return err
	}
	Stat <- body
	return err
}

func ParseResponseBody(rc io.ReadCloser, data interface{}) error {
	var err error
	bData, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Println("[parseResponseBody] Error #1 ", err)
		return err
	}
	if err := json.Unmarshal(bData, &data); err != nil {
		log.Println("[parseResponseBody] Error #2 ", err)
		return err
	}
	return nil
}
