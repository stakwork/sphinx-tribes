package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
	"net/http"
)

func ChatRoutes() chi.Router {
	r := chi.NewRouter()
	chatHandler := handlers.NewChatHandler(http.DefaultClient, db.DB)

	r.Post("/response", chatHandler.ProcessChatResponse)

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)
		r.Post("/", chatHandler.CreateChat)
		r.Post("/send", chatHandler.SendMessage)
		r.Get("/history/{uuid}", chatHandler.GetChatHistory)
	})

	return r
}
