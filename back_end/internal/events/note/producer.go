package note

import "context"

type Producer interface {
	ProducerReadEvent(ctx context.Context, evt ReadEvent) error
}

type ReadEvent struct {
	Uid int64
	Nid int64
}
