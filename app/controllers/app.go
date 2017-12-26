package controllers

import (
	"errors"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/ajaygh/scuart/app/serial"
	r "github.com/revel/revel"
)

var errInvalidData = errors.New("invalid data")

// App represents application server
type App struct {
	*r.Controller
}

// Index returns home page
func (c App) Index() r.Result {
	return c.Render()
}

// Power serves power on/off request
func (c App) Power(id int) r.Result {
	r.INFO.Println("POWER REQUEST RECEIVED")

	on, err := strconv.ParseBool(c.Params.Get("on"))
	if err != nil {
		c.Response.Status = 400
		return c.RenderError(errInvalidData)
	}
	serial.Send(makePowerData(id, on))
	r.INFO.Println("power data sent")
	return c.Render()
}

// Dim serves dim change request
func (c App) Dim(id int) r.Result {
	r.INFO.Println("DIM REQUEST RECEIVED")

	dim, err := strconv.Atoi(c.Params.Get("dim"))
	if err != nil {
		c.Response.Status = 400
		return c.RenderError(errInvalidData)
	}
	serial.Send(makeDimData(id, dim))
	r.INFO.Println("dim data sent")
	return c.Render()
}

// CCT serves cct change request
func (c App) CCT(id int) r.Result {
	r.INFO.Println("CCT REQUEST RECEIVED")

	cct, err := strconv.Atoi(c.Params.Get("cct"))
	if err != nil {
		c.Response.Status = 400
		return c.RenderError(errInvalidData)
	}
	serial.Send(makeCCTData(id, cct))
	r.INFO.Println("cct data sent")
	return c.Render()
}

// LightStatus sends light status
func (c App) LightStatus(ws *websocket.Conn) r.Result {
	defer ws.Close()

	var msgRcv string
	websocket.Message.Receive(ws, &msgRcv)
	r.INFO.Println("connected: ", msgRcv)

	for {
		status := <-serial.ChanStatus

		if err := websocket.JSON.Send(ws, status); err != nil {
			r.ERROR.Println("error in sending")
			break
		}
		r.INFO.Println("successfully sent status", status)
	}

	return c.Render()
}

func makePowerData(id int, on bool) []byte {
	bOn := byte(0X00)
	if on {
		bOn = 0X64
	}
	bID := getID(id)

	return []byte{0X81, bOn, bID}
}

func makeDimData(id, dim int) []byte {
	bID := getID(id)

	return []byte{0X81, byte(dim), bID}
}

func makeCCTData(id, cct int) []byte {
	bID := getID(id)

	return []byte{0X81, byte(cct), bID}
}

func getID(id int) (bID byte) {
	switch id {
	case 1:
		bID = 0X02
	case 2:
		bID = 0X05
	}
	return
}
