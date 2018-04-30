package engine

import (
	"github.com/venushome/noolite/rx"
	"github.com/venushome/noolite/tx"
)

var rxActions map[byte]string = map[byte]string{
	rx.RX_RESP_TURN_OFF:    "off",
	rx.RX_RESP_DIM_DOWN:    "dim_down",
	rx.RX_RESP_TURN_ON:     "on",
	rx.RX_RESP_DIM_UP:      "dim_up",
	rx.RX_RESP_TURN_CHANGE: "switch",
	rx.RX_RESP_DIM_CHANGE:  "dim_direction",
	rx.RX_RESP_DIM_SET:     "dim_set",
	rx.RX_RESP_SCENE_RUN:   "scene",
	rx.RX_RESP_SCENE_SAVE:  "save_scene",
	rx.RX_RESP_CLEAR_ADDR:  "unbind",
	rx.RX_RESP_DIM_STOP:    "dim_stop",
	rx.RX_RESP_WANT_BIND:   "want_bind",
	rx.RX_RESP_RAINBOW:     "rainbow",
	rx.RX_RESP_SET_COLOR:   "color",
	rx.RX_RESP_WORK_MODE:   "mode",
	rx.RX_RESP_WORK_SPEED:  "speed",
	rx.RX_RESP_LOW_BATERY:  "low_batery",
	rx.RX_RESP_SENSE_INFO:  "sens",
}

var txActions map[string]byte = map[string]byte{
	"off":     tx.TX_CMD_TURN_OFF,
	"on":      tx.TX_CMD_TURN_ON,
	"switch":  tx.TX_CMD_TURN_CHANGE,
	"dim":     tx.TX_CMD_DIM_SET,
	"scene":   tx.TX_CMD_SCENE_RUN,
	"rainbow": tx.TX_CMD_RAINBOW,
	"color":   tx.TX_CMD_CHANGE_COLOR,
	"mode":    tx.TX_CMD_WORK_MODE,
	"speed":   tx.TX_CMD_WORK_SPEED,
}
