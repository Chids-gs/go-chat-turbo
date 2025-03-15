package handlers

import (
	"fmt"
	"log"
	"net/http"

	"go-chat/internal/config"
	"go-chat/internal/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var geminiService *services.GeminiService

func init() {
	cfg := config.LoadConfig()
	geminiService = services.NewGeminiService(cfg)
}

func buildChatResponse(content string) string {
	messageContent := fmt.Sprintf(`<div class="markdown">%s</div>`, content)
	feedback := `<div class="feedback">
        <img src="/static/images/thumbs-up.png" class="feedback-icon" alt="Thumbs Up" onclick="giveFeedback('up')"/>
        <img src="/static/images/thumbs-down.jpeg" class="feedback-icon" alt="Thumbs Down" onclick="giveFeedback('down')"/>
    </div>`

	return fmt.Sprintf(`
<turbo-stream action="append" target="box">
    <template>
        <div class="message bot">
            <img src="/static/images/bot-chat.png" class="icon bot-icon" alt="Bot" />
            %s
            %s
        </div>
    </template>
</turbo-stream>`, messageContent, feedback)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// Get response from Gemini
		rawResponse, err := geminiService.GenerateResponse(string(message))
		if err != nil {
			log.Printf("Gemini API error: %v", err)
			errorResponse := buildChatResponse("Sorry, I encountered an error processing your request.")
			if err := conn.WriteMessage(messageType, []byte(errorResponse)); err != nil {
				log.Println("write error:", err)
			}
			continue
		}

		// Build response content based on response type
		var responseContent string
		switch {
		case rawResponse.Text != "":
			responseContent = rawResponse.Text
		case rawResponse.Table != "":
			responseContent = rawResponse.Table
		case rawResponse.Diagram != "":
			responseContent = rawResponse.Diagram
		default:
			responseContent = "No content generated"
		}

		// Send response
		response := buildChatResponse(responseContent)
		if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
