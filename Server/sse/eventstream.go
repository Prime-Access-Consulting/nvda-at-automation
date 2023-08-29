package sse

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	linePattern = `^(?:(event|id): (.+)|(data):(.*))$`
	eventType   = `event`
	idType      = `id`
	dataType    = `data`
)

type EventHandler func(event Event)

type Parser struct {
	currentEvent *Event
	handler      EventHandler
}

type Event struct {
	Type *string `json:"type"`
	ID   *int    `json:"id"`
	Data *string `json:"data"`
}

func NewParser(handler EventHandler) *Parser {
	return &Parser{handler: handler}
}

func (p *Parser) Process(line string) {
	isValid, _ := regexp.MatchString(linePattern, line)

	if !isValid {
		fmt.Printf("Unexpected line: \"%s\"\n", line)
		return
	}

	re := regexp.MustCompile(linePattern)
	matches := re.FindStringSubmatch(line)

	lineType := matches[1]
	lineData := matches[2]

	if matches[3] == dataType {
		lineType = matches[3]
		lineData = matches[4]
	}

	err := p.handleData(lineType, lineData)

	if err != nil {
		fmt.Printf("Processing error on line %s: %s\n", line, err.Error())
	}

	p.tryDispatchEvent()
}

func (p *Parser) tryDispatchEvent() {
	if p.currentEvent == nil ||
		p.currentEvent.ID == nil ||
		p.currentEvent.Type == nil ||
		p.currentEvent.Data == nil {
		return
	}

	p.handler(*p.currentEvent)
	p.currentEvent = nil
}

func (p *Parser) handleData(someType string, data string) error {
	switch someType {
	case idType:
		return p.handleIdLine(data)
	case eventType:
		return p.handleEventLine(data)
	case dataType:
		return p.handleDataLine(data)
	}

	return nil
}

func (p *Parser) handleIdLine(id string) error {
	intId, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	if p.currentEvent == nil {
		p.currentEvent = &Event{ID: &intId}
		return nil
	}

	if p.currentEvent.ID != nil {
		return fmt.Errorf("CurrentEvent already has ID %d while trying to set %d", *p.currentEvent.ID, intId)
	}

	p.currentEvent.ID = &intId
	return nil
}

func (p *Parser) handleEventLine(eventType string) error {
	if p.currentEvent == nil {
		p.currentEvent = &Event{Type: &eventType}
		return nil
	}

	if p.currentEvent.Type != nil {
		return fmt.Errorf("CurrentEvent already has Type %s while trying to set %s", *p.currentEvent.Type, eventType)
	}

	p.currentEvent.Type = &eventType
	return nil
}

func (p *Parser) handleDataLine(eventData string) error {
	if p.currentEvent == nil {
		p.currentEvent = &Event{Data: &eventData}
		return nil
	}

	if p.currentEvent.Data != nil {
		return fmt.Errorf("CurrentEvent already has Data %s while trying to set %s", *p.currentEvent.Data, eventData)
	}

	p.currentEvent.Data = &eventData
	return nil
}
