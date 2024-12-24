package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"

	"langtools/config"
	"langtools/globals"
	"langtools/handlers"
	"langtools/message"
	"langtools/tools"
	"langtools/utils"
)

type Message struct {
	Message string `json:"message"`
}

type FormattedResult struct {
	PageContent string                 `json:"pageContent"`
	Metadata    map[string]interface{} `json:"metadata"`
}

var memory = NewMemory() // Memoria global
func HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	// Simular un ID único por usuario
	userID := conn.RemoteAddr().String()

	// Cargar prompt inicial
	prompt, err := config.LoadPrompt("static/prompt.txt")
	if err != nil {
		utils.Error("Error loading prompt: %w", err)
		return
	}

	// Inicializar OpenAI y Pinecone Store
	llm, err := config.InitOpenAI()
	if err != nil {
		log.Printf("Error initializing OpenAI: %v", err)
		return
	}
	pineconeStore, err := config.GetPineconeStore(llm)
	if err != nil {
		log.Printf("Error initializing Pinecone Store: %v", err)
		return
	}

	availableTools := tools.RegisterTools()
	ctx := context.Background()

	// Crear historial inicial con el prompt
	messageHistory := memory.GetHistory(userID)
	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, prompt))
	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, "El id del usuario es (userId): "+userID))

	log.Println("Message History: ", messageHistory)

	for {
		// Leer mensaje del cliente
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Delete Pinecone Results Cache
		globals.PineconeResultsCache[userID] = nil

		// Parsear el mensaje
		var parsedMsg Message
		err = json.Unmarshal(msg, &parsedMsg)
		if err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}
		utils.Info("Received from " + userID + ": " + parsedMsg.Message)

		// Agregar mensaje del usuario al historial
		messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, string(parsedMsg.Message)))

		// Generar respuesta inicial
		resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableTools))
		if err != nil {
			log.Printf("Error generating content: %v", err)
			break
		}

		// Actualizar historial y ejecutar herramientas
		messageHistory = message.UpdateHistory(messageHistory, resp)
		messageHistory = handlers.ExecuteToolCalls(ctx, messageHistory, resp, pineconeStore)

		// Guardar historial en memoria
		memory.UpdateHistory(userID, messageHistory)

		messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, "Recuerda no incluir imagenes o links en tus mensajes."))

		// Generar respuesta final
		resp, err = llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableTools))
		if err != nil {
			log.Printf("Error generating final content: %v", err)
			break
		}

		// log.Println("Pinecone Results Cache:" + utils.PrettyPrint(globals.PineconeResultsCache))

		if resp.Choices[0].Content == "" {
			utils.Info("No content to send.")
			messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, string(parsedMsg.Message)))
			continue
		}
		utils.Info("Sending to " + userID + ": " + resp.Choices[0].Content)
		results, ok := globals.PineconeResultsCache[userID]
		if !ok || len(results) == 0 {
			utils.Info("No results found or tool not used.")

			// Respuesta genérica con resultados nulos
			err := SendResponse(conn, "success", resp.Choices[0].Content, nil)
			if err != nil {
				log.Printf("Error sending response: %v", err)
				break
			}
		} else {
			// Formatear resultados
			formattedResults := FormatResults(results)

			// Enviar respuesta con los resultados formateados
			err := SendResponse(conn, "success", resp.Choices[0].Content, formattedResults)
			if err != nil {
				log.Printf("Error sending response: %v", err)
				break
			}
		}

		// Agregar mensaje del bot al historial
		messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, resp.Choices[0].Content))
		log.Println("Message History: ", messageHistory)

	}
}
