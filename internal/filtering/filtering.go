package filtering

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/ArthurMVilela/har-tools/pkg/model"
	"github.com/antchfx/jsonquery"
	"github.com/rs/zerolog"
)

type EntryProcessor struct {
	logger  *zerolog.Logger
	filters []EntryFilter
}

type EntryProcessorOption func(p *EntryProcessor) error

type EntryFilter func(entry model.Entry) (bool, error)

func NewEntryProcessor(options ...EntryProcessorOption) (*EntryProcessor, error) {
	processor := new(EntryProcessor)

	for _, option := range options {
		err := option(processor)
		if err != nil {
			return nil, fmt.Errorf("unable to create new processor: %w", err)
		}
	}

	if processor.logger == nil {
		return nil, errors.New("missing logger in processor")
	}

	return processor, nil
}

func (proc *EntryProcessor) ApplyFilters(entries []model.Entry) ([]model.Entry, error) {
	return slices.DeleteFunc(entries, func(entry model.Entry) bool {
		for _, filter := range proc.filters {
			passes, err := filter(entry)
			if err != nil {
				proc.logger.Debug().Err(err).Msgf("Unable to process filter: %v. Entry: %+v", filter, entry)
				continue
			}
			if passes {
				return false
			}
		}
		return true
	}), nil
}

func WithLogger(logger *zerolog.Logger) EntryProcessorOption {
	return func(p *EntryProcessor) error {
		p.logger = logger
		return nil
	}
}

func WithEntryFilters(filters ...EntryFilter) EntryProcessorOption {
	return func(p *EntryProcessor) error {
		p.filters = filters
		return nil
	}
}

func XPathEntryContentFilter(filter string) EntryFilter {
	return func(entry model.Entry) (bool, error) {
		if entry.Response.Content.Size == 0 || len(entry.Response.Content.Text) == 0 {
			return false, nil
		}

		switch {
		case strings.HasPrefix(entry.Response.Content.MimeType, "application/json"):
			return satisfyJSONXPathFilter(entry.Response.Content.Text, filter)
		default:
			return false, nil
		}
	}
}

func satisfyJSONXPathFilter(text string, filter string) (bool, error) {
	reader := strings.NewReader(text)
	node, err := jsonquery.Parse(reader)
	if err != nil {
		return false, fmt.Errorf("unable to parse JSON content: %w", err)
	}

	node, err = jsonquery.Query(node, filter)
	if err != nil {
		return false, err
	}

	return node != nil, nil
}

func MimeTypeContentFilter(pattern string) EntryFilter {
	return func(entry model.Entry) (bool, error) {
		return regexp.MatchString(pattern, entry.Response.Content.MimeType)
	}
}
