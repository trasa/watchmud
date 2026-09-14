package player

import "github.com/rs/zerolog/log"

type Recorder struct {
	Sent []any
}

func (r *Recorder) Send(msg any) {
	log.Debug().Msgf("sending message: %v", msg)
	r.Sent = append(r.Sent, msg)
}

func (r *Recorder) Clear() {
	r.Sent = nil
}
