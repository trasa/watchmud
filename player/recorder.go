package player

import "github.com/rs/zerolog/log"

type Recorder struct {
	Sent []interface{}
}

func (r *Recorder) Send(msg interface{}) {
	log.Debug().Msgf("sending message: %v", msg)
	r.Sent = append(r.Sent, msg)
}
