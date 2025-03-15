// Initialize Marked
function initMarked() {
    // No custom renderer needed; we'll process Mermaid blocks later.
}

// Process Mermaid Diagrams
function renderMermaidDiagrams(container) {
    const mermaidRegex = /```mermaid\s*([\s\S]*?)```/g;
    const originalHTML = container.innerHTML;
    let updatedHTML = originalHTML;
    let match;

    while ((match = mermaidRegex.exec(originalHTML)) !== null) {
        const fullBlock = match[0];
        const code = match[1].trim();
        const id = 'mermaid-' + Date.now() + Math.random().toString(36).substr(2, 5);

        const tempDiv = document.createElement('div');
        tempDiv.className = 'mermaid';
        tempDiv.id = id;
        tempDiv.textContent = code;

        container.insertAdjacentElement('beforeend', tempDiv);

        try {
            mermaid.render(id, code, function (svgCode) {
                updatedHTML = updatedHTML.replace(fullBlock, svgCode);
                container.innerHTML = updatedHTML;
            });
        } catch (error) {
            console.error("Error rendering Mermaid diagram:", error);
        }
    }
}

// Markdown and Mermaid Rendering
function processMarkdownAndMermaid(element) {
    try {
        const content = element.textContent;
        element.innerHTML = marked.parse(content);
        renderMermaidDiagrams(element);
    } catch (error) {
        console.error("Error in processMarkdownAndMermaid:", error);
    }
}

// Initialize Mermaid
function initMermaid() {
    mermaid.initialize({
        startOnLoad: false,
        theme: "default",
        securityLevel: "loose",
        flowchart: { useMaxWidth: true, htmlLabels: true, curve: "basis" },
        sequence: { useMaxWidth: true, showSequenceNumbers: false }
    });
}

// Chat Initialization
function initChat() {
    const ws = new WebSocket("ws://" + window.location.host + "/ws");
    const messagesContainer = document.getElementById("box");
    const input = document.getElementById("line");

    if (!messagesContainer || !input) {
        console.error("Required DOM elements not found");
        return;
    }

    ws.onopen = () => console.log("WebSocket connected");

    ws.onmessage = function (event) {
        try {
            
            
            // Remove spinner
            const spinner = document.getElementById("processing-spinner");
            if (spinner) {
                spinner.remove();
            }

            // Process the message content
            const parser = new DOMParser();
            const doc = parser.parseFromString(event.data, 'text/html');
            const template = doc.querySelector('template');
            
            if (template) {
                const messageContent = template.innerHTML;
                messagesContainer.insertAdjacentHTML('beforeend', messageContent);

                // Process markdown after adding the content
                const markdownElements = messagesContainer.querySelectorAll(".markdown");
                console.log("Found markdown elements:", markdownElements.length);
                markdownElements.forEach(el => {
                    
                    processMarkdownAndMermaid(el);
                });
            }
        } catch (error) {
            console.error("Error processing message:", error);
        }
    };

    // Define onSend and handleKeyPress globally
    window.onSend = function () {
        const text = input.value.trim();
        if (text === "") return;

        // Append user message with updated icon
        messagesContainer.innerHTML += `<div class="message user">
            <img src="/static/images/user-chat.png" class="icon user-icon" alt="User" />
            <span>${text}</span>
        </div>`;
        
        // Updated spinner HTML with proper styling
        const spinnerHTML = `<div id="processing-spinner" class="message bot" style="text-align: center;">
            <img src="/static/images/processing.webp" class="spinner" alt="Processing..." style="width: 48px; height: 48px;"/>
        </div>`;
        messagesContainer.innerHTML += spinnerHTML;

        ws.send(text);
        input.value = "";
    };

    window.handleKeyPress = function(event) {
        if (event.key === "Enter") {
            window.onSend();
        }
    };
}

// Process JSON Response
function renderResponse(response) {
    let content = "";

    try {
        const data = JSON.parse(response);
        if (data.text) content += `<div class="markdown">${marked.parse(data.text)}</div>`;
        if (data.table) content += `<div class="markdown">${marked.parse(data.table)}</div>`;
        if (data.diagram) content += `<div class="mermaid">${data.diagram}</div>`;
    } catch {
        content = `<div class="markdown">${marked.parse(response)}</div>`;
    }

    const messagesContainer = document.getElementById("messages");
    const messageDiv = document.createElement("div");
    messageDiv.className = "message bot";
    messageDiv.innerHTML = content;
    messagesContainer.appendChild(messageDiv);

    // Render Mermaid diagrams
    setTimeout(() => {
        mermaid.init(undefined, ".mermaid");
    }, 100);
}
document.addEventListener("turbo:load", initializeMermaid);
document.addEventListener("turbo:frame-load", initializeMermaid);
document.addEventListener("turbo:render", initializeMermaid);

function initializeMermaid() {
    // Check if Mermaid is available
    if (window.mermaid) {
        // Iterate through all new Mermaid diagrams
        document.querySelectorAll('.mermaid').forEach((el) => {
            // Ensure Mermaid initializes only once per element
            if (!el.dataset.processed) {
                el.dataset.processed = "true"; 
                mermaid.init(undefined, el);  // Initialize Mermaid
            }
        });
    }
}

window.addEventListener("load", () => {
    console.log("Initializing chat components...");
    initMermaid();
    initMarked();
    initChat();
});
