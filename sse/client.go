package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

var ClientRegistry = &Registry{
	clients: make(map[string]*Client),
	mutex:   &sync.RWMutex{},
}

type Registry struct {
	clients map[string]*Client
	mutex   *sync.RWMutex
}

func GenerateClientKey(chatID, sseURL string) string {
	return fmt.Sprintf("%s:%s", chatID, sseURL)
}

func (r *Registry) Register(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	key := GenerateClientKey(client.ChatID, client.URL)
	r.clients[key] = client
}

func (r *Registry) Unregister(sseURL, chatID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	key := GenerateClientKey(chatID, sseURL)
	logger.Log.Info("Attempting to unregister client with key: %s", key)
	
	if client, exists := r.clients[key]; exists {
		client.Stop()
		delete(r.clients, key)
		return true
	}
	
	logger.Log.Info("No client found. Currently registered clients:")
	for k := range r.clients {
		logger.Log.Info("- %s", k)
	}
	
	return false
}

type Client struct {
	URL           string
	ChatID        string
	WebhookURL    string
	LastEventID   string
	RetryInterval time.Duration
	Client        *http.Client
	DB            db.Database
	stopChan      chan struct{}
	firstFailTime time.Time
}

func NewClient(sseURL string, chatID string, webhookURL string, database db.Database) *Client {
	return &Client{
		URL:           sseURL,
		ChatID:        chatID,
		WebhookURL:    webhookURL,
		RetryInterval: 3 * time.Second,
		Client: &http.Client{
			Timeout: 0,
		},
		DB:       database,
		stopChan: make(chan struct{}),
	}
}

func (c *Client) Start() {

	ClientRegistry.Register(c)

	go func() {
		defer func() {

			ClientRegistry.Unregister(c.URL, c.ChatID)
		}()

		for {
			select {
			case <-c.stopChan:
				logger.Log.Info("[ChatID: %s] SSE client stopped", c.ChatID)
				return
			default:
				err := c.connect()
				if err != nil {
					if c.firstFailTime.IsZero() {
						c.firstFailTime = time.Now()
					} else if time.Since(c.firstFailTime) > 60*time.Minute {
						logger.Log.Error("[ChatID: %s] Server unreachable for 60 minutes, stopping client", c.ChatID)
						return
					}

					logger.Log.Error("[ChatID: %s] Connection error: %v. Retrying in %v...", c.ChatID, err, c.RetryInterval)
					time.Sleep(c.RetryInterval)
					continue
				}

				c.firstFailTime = time.Time{}

				time.Sleep(c.RetryInterval)
			}
		}
	}()
}

func (c *Client) Stop() {
	close(c.stopChan)
}

func (c *Client) connect() error {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	if c.LastEventID != "" {
		req.Header.Set("Last-Event-ID", c.LastEventID)
	}

	logger.Log.Info("[ChatID: %s] Connecting to SSE endpoint: %s", c.ChatID, c.URL)
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	logger.Log.Info("[ChatID: %s] Connected successfully, waiting for events...", c.ChatID, c.URL)
	return c.processEvents(resp)
}

func (c *Client) processEvents(resp *http.Response) error {
	scanner := bufio.NewScanner(resp.Body)
	eventData := map[string]string{
		"id":    "",
		"event": "",
		"data":  "",
	}

	for scanner.Scan() {
		select {
		case <-c.stopChan:
			return nil
		default:
			line := scanner.Text()

			if line == "" {
				if eventData["data"] != "" {

					err := c.storeEvent(eventData)
					if err != nil {
						logger.Log.Error("[ChatID: %s] Error storing event: %v", c.ChatID, err)
					}

					if eventData["id"] != "" {
						c.LastEventID = eventData["id"]
					}

					eventData = map[string]string{
						"id":    "",
						"event": "",
						"data":  "",
					}
				}
				continue
			}

			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")

				if len(data) > 0 && data[0] == ' ' {
					data = data[1:]
				}

				if eventData["data"] != "" {
					eventData["data"] += "\n"
				}
				eventData["data"] += data
			} else if strings.HasPrefix(line, "id:") {
				eventData["id"] = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
			} else if strings.HasPrefix(line, "event:") {
				eventData["event"] = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "retry:") {
				retryStr := strings.TrimSpace(strings.TrimPrefix(line, "retry:"))
				if retry, err := time.ParseDuration(retryStr + "ms"); err == nil {
					c.RetryInterval = retry
				}
			}
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("error reading events: %w", scanner.Err())
	}

	return nil
}

func (c *Client) storeEvent(eventData map[string]string) error {
	var parsedEvent db.PropertyMap

	if strings.TrimSpace(eventData["data"])[0] == '{' {
		err := json.Unmarshal([]byte(eventData["data"]), &parsedEvent)
		if err != nil {
			parsedEvent = db.PropertyMap{
				"raw": eventData["data"],
			}
		}
	} else {
		parsedEvent = db.PropertyMap{
			"raw": eventData["data"],
		}
	}

	if eventData["id"] != "" {
		parsedEvent["id"] = eventData["id"]
	}
	if eventData["event"] != "" {
		parsedEvent["event_type"] = eventData["event"]
	}

	messageLog, err := c.DB.CreateSSEMessageLog(parsedEvent, c.ChatID, c.URL, c.WebhookURL)
	if err != nil {
		return fmt.Errorf("failed to create SSE message log: %w", err)
	}

	logger.Log.Info("[ChatID: %s] Stored SSE event with ID: %s", c.ChatID, messageLog.ID)
	return nil
}

func (r *Registry) HasClient(sseURL, chatID string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	key := GenerateClientKey(chatID, sseURL)
	logger.Log.Info("Checking for client with key: %s", key)
	_, exists := r.clients[key]

	return exists
}
