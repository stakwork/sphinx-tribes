package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func SkillRoutes() chi.Router {
	r := chi.NewRouter()
	skillHandler := handlers.NewSkillHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Use(auth.CombinedAuthContext)
		
		r.Post("/", skillHandler.CreateSkill)
		r.Get("/", skillHandler.GetAllSkills)
		r.Get("/{id}", skillHandler.GetSkillByID)
		r.Put("/{id}", skillHandler.UpdateSkillByID)
		r.Delete("/{id}", skillHandler.DeleteSkillByID)
		
		r.Post("/install/{id}", skillHandler.CreateSkillInstall)
		r.Get("/install/{id}", skillHandler.GetAllSkillInstallsBySkillID)
		r.Get("/install/detail/{id}", skillHandler.GetSkillInstallByID)
		r.Put("/install/detail/{id}", skillHandler.UpdateSkillInstallByID)
		r.Delete("/install/detail/{id}", skillHandler.DeleteSkillInstallByID)
	})

	return r
} 