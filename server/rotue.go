package server

import (
	"edServer/sse"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
	// e.GET("/dashboard", Dashboard)
	e.GET("/dashboard/*", echo.WrapHandler(http.StripPrefix("/dashboard", http.FileServer(http.Dir("./htmltemplate/")))))
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
	return c.Render(http.StatusOK, "dashboard", nil)
}

func StatusHandler(ctx echo.Context) error {

	s := Status{}
	err := ParseResponseBody(ctx.Request().Body, &s)
	if err != nil {
		log.Println("StatusHandler parse response body error:", err)
		return err
	}
	parseFlags(&s)
	s.lastUpdate = time.Now()
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

func parseFlags(s *Status) {
	bytePos := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 16, 17, 18, 19, 20, 22, 27, 28, 30}
	for i := range bytePos {
		if (s.Flags>>bytePos[i])&1 == 1 {
			if bytePos[i] == 0 {
				s.Docked = true
			}
			if bytePos[i] == 1 {
				s.Landed = true
			}
			if bytePos[i] == 2 {
				s.LandGear = true
			}
			if bytePos[i] == 3 {
				s.Shields = true
			}
			if bytePos[i] == 4 {
				s.Supercruise = true
			}
			if bytePos[i] == 5 {
				s.FAOff = true
			}
			if bytePos[i] == 6 {
				s.WeaponOn = true
			}
			if bytePos[i] == 7 {
				s.InWing = true
			}
			if bytePos[i] == 8 {
				s.Lights = true
			}
			if bytePos[i] == 9 {
				s.CargoScope = true
			}
			if bytePos[i] == 10 {
				s.Silents = true
			}
			if bytePos[i] == 11 {
				s.FuelScoope = true
			}
			if bytePos[i] == 16 {
				s.FSDMassLock = true
			}
			if bytePos[i] == 17 {
				s.FSDCharge = true
			}
			if bytePos[i] == 18 {
				s.FSDCooldown = true
			}
			if bytePos[i] == 19 {
				s.LowFuel = true
			}
			if bytePos[i] == 20 {
				s.OverHeating = true
			}
			if bytePos[i] == 22 {
				s.IsInDanger = true
			}
			if bytePos[i] == 27 {
				s.AnalysMode = true
			}
			if bytePos[i] == 28 {
				s.NVisin = true
			}
			if bytePos[i] == 30 {
				s.FSDJump = true
			}
		}
	}

}

func CheckAFK() {
	s := Shutdown{}
	// log.Println("CHECK AFK START")
	for key := range StatusBoard {
		if time.Now().After(StatusBoard[key].lastUpdate.Add(1 * time.Second)) {
			log.Println("DELETE ", StatusBoard[key].lastUpdate, StatusBoard[key].lastUpdate.Add(5*time.Minute).After(time.Now()), StatusBoard[key].lastUpdate.Add(5*time.Minute))
			s.PlayerID = key
			s.Delete = true
			body, err := json.Marshal(s)
			if err != nil {
				log.Println("ShutdownHandler marshal error:", err)
				return
			}
			Stat <- body
			delete(StatusBoard, key)
		}
	}
	// log.Println("CHECK AFK END")

}
