package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func ChatRoutes() chi.Router {
	r := chi.NewRouter()
	chatHandler := handlers.NewChatHandler(http.DefaultClient, db.DB)

	r.Post("/response", chatHandler.ProcessChatResponse)

	r.Group(func(r chi.Router) {
		r.Use(auth.CombinedAuthContext)

		r.Get("/", chatHandler.GetChat)
		r.Post("/", chatHandler.CreateChat)
		r.Put("/{chat_id}", chatHandler.UpdateChat)
		r.Put("/{chat_id}/archive", chatHandler.ArchiveChat)
		r.Post("/send", chatHandler.SendMessage)
		r.Get("/history/{uuid}", chatHandler.GetChatHistory)
		r.Post("/send/build", chatHandler.SendBuildMessage)
		r.Post("/send/action", chatHandler.SendActionMessage)

		r.Post("/upload", chatHandler.UploadFile)
		r.Get("/file/{id}", chatHandler.GetFile)
		r.Get("/file/all", chatHandler.ListFiles)
		r.Delete("/file/{id}", chatHandler.DeleteFile)

		r.Post("/artefacts", chatHandler.CreateArtefact)
		r.Get("/artefacts/chat/{chatId}", chatHandler.GetArtefactsByChatID)
		r.Get("/artefacts/{artifactId}", chatHandler.GetArtefactByID)
		r.Get("/artefacts/message/{messageId}", chatHandler.GetArtefactsByMessageID)
		r.Put("/artefacts/{artifactId}", chatHandler.UpdateArtefact)
		r.Delete("/artefacts/{artifactId}", chatHandler.DeleteArtefactByID)
		r.Delete("/artefacts/chat/{chatId}", chatHandler.DeleteAllArtefactsByChatID)

		r.Post("/chatworkflow", chatHandler.CreateOrEditChatWorkflow)
		r.Get("/chatworkflow/{workspaceId}", chatHandler.GetChatWorkflow)
		r.Delete("/chatworkflow/{workspaceId}", chatHandler.DeleteChatWorkflow)

		r.Post("/sse/stop", chatHandler.StopSSEClient)
		r.Get("/sse/{chat_id}", chatHandler.GetSSEMessagesByChatID)
	})

	return r
}
