package feedback

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	SchemaVersion = "sdp-trace-feedback-event-v1"
	MaxBodyBytes  = 16 * 1024
)

var safeTokenPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]+$`)
var eventIDPattern = regexp.MustCompile(`^[0-9]{8}T[0-9]{6}Z-[A-Za-z0-9_.:-]+-[0-9a-f]{12}-[0-9a-f]{8}$`)

type Options struct {
	Kind        string
	From        string
	To          string
	SourceRef   string
	Summary     string
	MessageFile string
	Now         time.Time
}

type Event struct {
	SchemaVersion string `json:"schema_version"`
	EventID       string `json:"event_id"`
	EventType     string `json:"event_type"`
	Kind          string `json:"kind"`
	From          string `json:"from"`
	To            string `json:"to"`
	SourceRef     string `json:"source_ref,omitempty"`
	Summary       string `json:"summary"`
	Message       struct {
		Retained bool   `json:"retained"`
		SHA256   string `json:"sha256"`
		Bytes    int    `json:"bytes"`
		Body     string `json:"body"`
	} `json:"message"`
	TrustScope string `json:"trust_scope"`
	ProofState string `json:"proof_state"`
	CreatedAt  string `json:"created_at"`
}

func Record(opts Options) (Event, error) {
	opts, err := normalize(opts)
	if err != nil {
		return Event{}, err
	}
	body, err := os.ReadFile(opts.MessageFile)
	if err != nil {
		return Event{}, err
	}
	if len(body) == 0 {
		return Event{}, errors.New("feedback message file is empty")
	}
	if len(body) > MaxBodyBytes {
		return Event{}, fmt.Errorf("feedback message exceeds %d bytes", MaxBodyBytes)
	}
	sum := sha256.Sum256(body)
	bodyText := string(body)
	event := Event{
		SchemaVersion: SchemaVersion,
		EventType:     "feedback",
		Kind:          opts.Kind,
		From:          opts.From,
		To:            opts.To,
		SourceRef:     opts.SourceRef,
		Summary:       opts.Summary,
		TrustScope:    "local_structural",
		ProofState:    "not_assessed",
		CreatedAt:     opts.Now.Format(time.RFC3339),
	}
	id, err := eventID(opts.Now, opts.Kind, sum[:])
	if err != nil {
		return Event{}, err
	}
	event.EventID = id
	event.Message.Retained = true
	event.Message.SHA256 = hex.EncodeToString(sum[:])
	event.Message.Bytes = len(body)
	event.Message.Body = bodyText
	return event, nil
}

func WriteJSON(path string, event Event) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("feedback observe requires --out")
	}
	clean := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(clean, data, 0o644)
}

func ValidateEvent(event Event) error {
	if event.SchemaVersion != SchemaVersion {
		return errors.New("feedback event has unsupported schema_version")
	}
	if event.EventType != "feedback" {
		return errors.New("feedback event_type must be feedback")
	}
	if !eventIDPattern.MatchString(event.EventID) {
		return errors.New("feedback event_id has invalid shape")
	}
	for label, value := range map[string]string{
		"kind": event.Kind,
		"from": event.From,
		"to":   event.To,
	} {
		if value == "" {
			return fmt.Errorf("feedback event requires %s", label)
		}
		if !safeTokenPattern.MatchString(value) {
			return fmt.Errorf("feedback event %s has unsafe token", label)
		}
	}
	if event.SourceRef != "" && !safeTokenPattern.MatchString(event.SourceRef) {
		return errors.New("feedback event source_ref has unsafe token")
	}
	if strings.TrimSpace(event.Summary) == "" || len(event.Summary) > 240 {
		return errors.New("feedback event summary must be 1-240 characters")
	}
	if event.TrustScope != "local_structural" {
		return errors.New("feedback event trust_scope must be local_structural")
	}
	if event.ProofState != "not_assessed" {
		return errors.New("feedback event proof_state must be not_assessed")
	}
	if _, err := time.Parse(time.RFC3339, event.CreatedAt); err != nil {
		return errors.New("feedback event created_at must be RFC3339")
	}
	if !event.Message.Retained {
		return errors.New("feedback event message must be retained")
	}
	body := []byte(event.Message.Body)
	if len(body) == 0 {
		return errors.New("feedback event message body is empty")
	}
	if len(body) > MaxBodyBytes || event.Message.Bytes != len(body) {
		return errors.New("feedback event message byte count is invalid")
	}
	sum := sha256.Sum256(body)
	if event.Message.SHA256 != hex.EncodeToString(sum[:]) {
		return errors.New("feedback event message sha256 mismatch")
	}
	return nil
}

func normalize(opts Options) (Options, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	opts.Kind = strings.TrimSpace(opts.Kind)
	opts.From = strings.TrimSpace(opts.From)
	opts.To = strings.TrimSpace(opts.To)
	opts.SourceRef = strings.TrimSpace(opts.SourceRef)
	opts.Summary = strings.TrimSpace(opts.Summary)
	opts.MessageFile = strings.TrimSpace(opts.MessageFile)
	if opts.Kind == "" {
		opts.Kind = "corrective_feedback"
	}
	for label, value := range map[string]string{
		"kind": opts.Kind,
		"from": opts.From,
		"to":   opts.To,
	} {
		if value == "" {
			return Options{}, fmt.Errorf("feedback observe requires --%s", label)
		}
		if !safeTokenPattern.MatchString(value) {
			return Options{}, fmt.Errorf("feedback --%s must match [A-Za-z0-9_.:-]+", label)
		}
	}
	if opts.SourceRef != "" && !safeTokenPattern.MatchString(opts.SourceRef) {
		return Options{}, errors.New("feedback --source-ref must match [A-Za-z0-9_.:-]+")
	}
	if opts.Summary == "" {
		return Options{}, errors.New("feedback observe requires --summary")
	}
	if len(opts.Summary) > 240 {
		return Options{}, errors.New("feedback --summary must be 240 characters or fewer")
	}
	if opts.MessageFile == "" {
		return Options{}, errors.New("feedback observe requires --message-file")
	}
	return opts, nil
}

func eventID(now time.Time, kind string, sum []byte) (string, error) {
	nonce := make([]byte, 4)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s-%s-%s", now.UTC().Format("20060102T150405Z"), kind, hex.EncodeToString(sum)[:12], hex.EncodeToString(nonce)), nil
}
