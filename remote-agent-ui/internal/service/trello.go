package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TrelloService struct {
	APIKey string
	Token  string
}

func NewTrelloService(apiKey, token string) *TrelloService {
	return &TrelloService{APIKey: apiKey, Token: token}
}

func (t *TrelloService) CreateTicket(listID, name, desc string) error {
	url := fmt.Sprintf("https://api.trello.com/1/cards?key=%s\u0026token=%s", t.APIKey, t.Token)
	
	payload := map[string]string{
		"name":   name,
		"desc":   desc,
		"idList": listID,
	}
	
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("trello api error: %s", resp.Status)
	}
	
	return nil
}
