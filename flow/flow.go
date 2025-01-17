package flow

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sakajunquality/cloud-pubsub-events/gcrevent"
)

var (
	cfg *Config
)

type Flow struct {
	Env         string
	githubToken string
}

func New(c *Config) (*Flow, error) {
	cfg = c
	f := &Flow{
		githubToken: os.Getenv("FLOW_GITHUB_TOKEN"),
	}

	if f.githubToken == "" {
		return nil, errors.New("You need to specify a non-empty value for FLOW_GITHUB_TOKEN")
	}

	return f, nil
}

func (f *Flow) ProcessGCREvent(ctx context.Context, e gcrevent.Event) error {
	if e.Action != gcrevent.ActionInsert {
		return nil
	}

	parts := strings.Split(*e.Tag, ":")
	image, version := parts[0], parts[1]

	if image == "" || version == "" {
		return fmt.Errorf("Image format invalid: %s", *e.Tag)
	}

	return f.processImage(ctx, image, version)
}
