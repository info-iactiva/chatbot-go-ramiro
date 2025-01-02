package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"

	"langtools/config"
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

func HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	var messageHistory []llms.MessageContent
	var results []schema.Document
	var chatID string
	// Simular un ID único por usuario
	userID := uuid.New().String()

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
	// messageHistory = globals.GlobalMemory.GetHistory(userID)
	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, prompt))
	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, "El id del usuario es (userId): "+userID))

	for {
		// Leer mensaje del cliente
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

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

		if chatID == "" {
			// Crear chat en MongoDB
			chatID, err = tools.CreateEmptyChat(ctx, userID)
			if err != nil {
				log.Printf("Error creating chat: %v", err)
				return
			}

			messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, "el chatId de la conversación es (chatId): "+chatID))
			utils.Info("el chatId de la conversación es (chatId): " + chatID)
		}

		// Generar respuesta inicial
		resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableTools))
		if err != nil {
			log.Printf("Error generating content: %v", err)
			break
		}

		log.Println("initial resp: ", resp.Choices[0].Content)
		// Sí se solicitaron herramientas, ejecutarlas, si no regresar el mensaje inicial.
		if len(resp.Choices[0].ToolCalls) == 0 {
			log.Println("No tool calls")
			err := SendResponse(conn, "success", resp.Choices[0].Content, nil)
			if err != nil {
				log.Printf("Error sending response: %v", err)
				break
			}
			continue
		}

		// Actualizar historial y ejecutar herramientas
		messageHistory = message.UpdateHistory(messageHistory, resp)
		messageHistory, results = handlers.ExecuteToolCalls(ctx, messageHistory, resp, pineconeStore)

		// Generar respuesta final
		resp, err = llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableTools))
		if err != nil {
			log.Printf("Error generating final content: %v", err)
			break
		}

		// log.Println("Pinecone Results Cache:" + utils.PrettyPrint(globals.PineconeResultsCache))

		if resp.Choices[0].Content == "" {
			utils.Error("No content to send.", nil)
			// messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, string(parsedMsg.Message)))
			// continue
			err := SendResponse(conn, "success", "Hubo un error, intente más tarde.", nil)
			if err != nil {
				utils.Error("Error sending response: %v", err)
				break
			}
			continue
		}
		if len(results) == 0 {
			log.Println("No pinecone results")
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

		// Actualiza la conversación en mongo con el historial
		err = tools.UpdateChatHistory(chatID, userID, messageHistory)
		if err != nil {
			log.Printf("Error updating chat history: %v", err)
		}

	}
}
