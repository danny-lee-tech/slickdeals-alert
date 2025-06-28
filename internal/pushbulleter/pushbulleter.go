package pushbulleter

import (
	"github.com/xconstruct/go-pushbullet"
)

type PushBulleter struct {
	APIKey string
	Tag    string
}

func (p *PushBulleter) PostToChannel(messageString string) error {
	pb := pushbullet.New(p.APIKey)

	err := pb.PushNoteToChannel(p.Tag, "New Slickdeals Alert", messageString)
	if err != nil {
		return err
	}

	return nil
}
