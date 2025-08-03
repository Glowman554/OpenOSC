package chatbox

import (
	"log"
	"regexp"
	"strings"

	"github.com/hypebeast/go-osc/osc"
)

type ChatBoxLine struct {
	expectedPlaceholders []string
	line                 string
}

func NewChatBoxLine(line string) *ChatBoxLine {
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindAllStringSubmatch(line, -1)

	var placeholders []string
	for _, match := range matches {
		if len(match) > 1 {
			placeholders = append(placeholders, match[1])
		}
	}

	return &ChatBoxLine{
		expectedPlaceholders: placeholders,
		line:                 line,
	}
}

func (c *ChatBoxLine) Applies(parent *ChatBoxBuilder) bool {
	for _, expected := range c.expectedPlaceholders {
		if _, ok := parent.placeholders[expected]; !ok {
			return false
		}
	}
	return true
}

func (c *ChatBoxLine) GetLine(parent *ChatBoxBuilder) string {
	line := c.line
	for _, expected := range c.expectedPlaceholders {
		line = strings.ReplaceAll(line, "{"+expected+"}", parent.placeholders[expected])
	}
	return line
}

type ChatBoxBuilder struct {
	lines        []*ChatBoxLine
	placeholders map[string]string
}

func NewChatBoxBuilder() *ChatBoxBuilder {
	return &ChatBoxBuilder{
		lines:        []*ChatBoxLine{},
		placeholders: map[string]string{},
	}
}

func (c *ChatBoxBuilder) AddLine(line string) {
	c.lines = append(c.lines, NewChatBoxLine(line))
}

func (c *ChatBoxBuilder) BeginTick() {
	for k := range c.placeholders {
		delete(c.placeholders, k)
	}
}

func (c *ChatBoxBuilder) EndTick(client *osc.Client) error {
	chatbox := ""

	for _, line := range c.lines {
		if line.Applies(c) {
			chatbox += line.GetLine(c) + "\n"
		}
	}

	msg := osc.NewMessage("/chatbox/input")
	msg.Append(chatbox)
	msg.Append(true)
	msg.Append(true)

	err := client.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	return nil
}

func (c *ChatBoxBuilder) Placeholder(placeholder string, text string) {
	c.placeholders[placeholder] = text
}
