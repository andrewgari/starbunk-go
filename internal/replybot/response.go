package replybot

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

// randomPlaceholderRe matches {random:min-max:char} placeholders.
var randomPlaceholderRe = regexp.MustCompile(`\{random:(\d+)-(\d+):(.)\}`)

// swapPlaceholderRe matches {swap_message:word1:word2} placeholders.
var swapPlaceholderRe = regexp.MustCompile(`\{swap_message:([^:}]+):([^:}]+)\}`)

// MessageData holds the author and content fields from a Discord message
// used for response template expansion. Keeping this separate from discordgo
// types makes the response layer independently testable.
type MessageData struct {
	Content       string // raw message text
	AuthorUsername string // display/username of the message author
	AuthorID      string // Discord user ID of the author
}

// ResponsePool holds a set of response strings and picks one at random,
// applying template substitutions against the triggering message.
type ResponsePool struct {
	responses []string
	rng       func(n int) int // injectable for deterministic testing
}

// NewResponsePool creates a pool from the given responses.
// Returns an error if the slice is empty.
func NewResponsePool(responses []string) (*ResponsePool, error) {
	return NewResponsePoolWithRng(responses, rand.IntN)
}

// NewResponsePoolWithRng creates a pool with a custom rng function.
// This is primarily useful for deterministic testing.
func NewResponsePoolWithRng(responses []string, rng func(int) int) (*ResponsePool, error) {
	if len(responses) == 0 {
		return nil, fmt.Errorf("response pool must contain at least one response")
	}
	return &ResponsePool{
		responses: responses,
		rng:       rng,
	}, nil
}

// Pick selects a random response from the pool and expands template
// placeholders against the provided message data.
//
// Template expansion order:
//  1. {swap_message:word1:word2} — replace word1 with word2 in the full message
//  2. {message}         — full message content
//  3. {start}           — first three words of the message
//  4. {author}          — message author's username
//  5. {author_id}       — message author's Discord user ID
//  6. {author_mention}  — Discord @mention (<@authorID>)
//  7. {random:min-max:char} — char repeated a random number of times
func (p *ResponsePool) Pick(data MessageData) string {
	response := p.responses[p.rng(len(p.responses))]
	return expandTemplates(response, data, p.rng)
}

// expandTemplates applies all template substitutions to a response string.
func expandTemplates(response string, data MessageData, rng func(int) int) string {
	response = expandSwapMessage(response, data.Content)
	response = strings.ReplaceAll(response, "{message}", data.Content)
	response = expandStart(response, data.Content)
	response = strings.ReplaceAll(response, "{author}", data.AuthorUsername)
	response = strings.ReplaceAll(response, "{author_id}", data.AuthorID)
	response = strings.ReplaceAll(response, "{author_mention}", "<@"+data.AuthorID+">")
	response = expandRandom(response, rng)
	return response
}

// expandSwapMessage replaces {swap_message:word1:word2} with the full message
// content after swapping word1 for word2 (case-insensitive).
func expandSwapMessage(response, messageContent string) string {
	return swapPlaceholderRe.ReplaceAllStringFunc(response, func(match string) string {
		parts := swapPlaceholderRe.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		from, to := parts[1], parts[2]
		re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(from))
		if err != nil {
			return match
		}
		return re.ReplaceAllString(messageContent, to)
	})
}

// expandStart replaces {start} with the first three words of the message.
func expandStart(response, messageContent string) string {
	if !strings.Contains(response, "{start}") {
		return response
	}
	words := strings.Fields(messageContent)
	if len(words) > 3 {
		words = words[:3]
	}
	return strings.ReplaceAll(response, "{start}", strings.Join(words, " "))
}

// expandRandom replaces {random:min-max:char} with the char repeated between
// min and max times (inclusive). Capped at 1000 to prevent abuse.
func expandRandom(response string, rng func(int) int) string {
	return randomPlaceholderRe.ReplaceAllStringFunc(response, func(match string) string {
		parts := randomPlaceholderRe.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		min, err1 := strconv.Atoi(parts[1])
		max, err2 := strconv.Atoi(parts[2])
		char := parts[3]
		if err1 != nil || err2 != nil || min > max || min < 0 {
			return match
		}
		if max > 1000 {
			max = 1000
		}
		count := min + rng(max-min+1)
		return strings.Repeat(char, count)
	})
}
