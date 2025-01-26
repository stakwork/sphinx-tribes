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
		r.Use(auth.PubKeyContext)
		r.Get("/", chatHandler.GetChat)
		r.Post("/", chatHandler.CreateChat)
		r.Put("/{chat_id}", chatHandler.UpdateChat)
		r.Put("/{chat_id}/archive", chatHandler.ArchiveChat)
		r.Post("/send", chatHandler.SendMessage)
		r.Get("/history/{uuid}", chatHandler.GetChatHistory)

		r.Post("/upload", chatHandler.UploadFile)
		r.Get("/file/{id}", chatHandler.GetFile)
		r.Get("/file/all", chatHandler.ListFiles)
		r.Delete("/file/{id}", chatHandler.DeleteFile)
	})

	return r
}
