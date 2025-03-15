# 🚀 Go-Turbo Chat with Gemini LLM

A demonstration project for Conf42 Golang 2025, showcasing how to revolutionize full-stack development by integrating Go with Hotwire Turbo. This project accompanies the session "Turbocharged Golang: Revolutionizing Full Stack Development with Hotwire" and demonstrates building modern, responsive chat applications with minimal JavaScript.

## Project Overview

This chat application demonstrates how to build dynamic, real-time web applications using:
- **Go** for backend processing
- **Hotwire Turbo** for dynamic UI updates
- **Google's Gemini LLM** for intelligent chat responses
- **WebSockets** for real-time communication

### Why Turbo with Go?

Traditional SPAs often require complex JavaScript frameworks and build pipelines. Turbo offers a simpler alternative:

✅ **Server-Side Driven UI**: Update UI directly from Go server
✅ **Zero Build Pipeline**: No npm, webpack, or transpilation needed
✅ **Real-Time Updates**: Native WebSocket support via Turbo Streams
✅ **Progressive Enhancement**: Works with basic JavaScript
✅ **Reduced Complexity**: Minimal client-side code

## Features

### Turbo Stream Implementation
```go
// Example of how we use Turbo Streams for real-time updates
turboStream := fmt.Sprintf(`
<turbo-stream action="append" target="box">
    <template>
        <div class="message bot">
            <img src="/static/images/bot-chat.png" class="icon bot-icon" alt="Bot" />
            %s
        </div>
    </template>
</turbo-stream>`, content)
```

### Key Features
- **Real-time Chat**: WebSocket-based communication
- **Markdown Rendering**: Support for rich text formatting
- **Mermaid Diagrams**: Automatic rendering of diagrams
- **Syntax Highlighting**: For code blocks
- **Responsive Design**: Mobile-friendly interface
- **Feedback System**: Thumbs up/down for responses

### LLM Integration
- Integration with Google's Gemini API


## Technical Implementation

### WebSocket Handler
```go
func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    // Handle real-time messages
    // Send Turbo Stream responses
}
```

### Turbo Stream Response Types
1. **Text Responses**: Regular markdown-formatted text
2. **Table Responses**: Structured data in table format
3. **Diagram Responses**: Mermaid diagram rendering
4. **Error Responses**: Graceful error handling

## Getting Started

1. Clone the repository
```bash
git clone [https://github.com/chids-gs/go-chat-turb.git]
```

2. Configure environment
```bash
# Create .env file in project root
echo "GEMINI_API_KEY=your_api_key" > .env
```

3. Run the application
```bash
go run main.go
```

## TODO Items


- [ ] OIDC Integration for Login/Signup
  - Implement OAuth2 flow
  - Support multiple providers (Google, GitHub)
  
- [ ] Persistence Layer
  - Store chat history
  - Save user feedback
  - Enable context-aware follow-up questions

### Future Enhancements
- [ ] Model Parameter Integration
  - Temperature control
  - Max tokens configuration
  - Top-P settings

- [ ] Enhanced UI Features
  - Chat history navigation
  - Response filtering
  - Theme customization

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE.md for details
