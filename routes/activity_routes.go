package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func ActivityRoutes() chi.Router {
	r := chi.NewRouter()
	activityHandler := handlers.NewActivityHandler(http.DefaultClient, db.DB)

	r.Group(func(r chi.Router) {
		r.Get("/{id}", activityHandler.GetActivity)
		r.Get("/thread/{thread_id}", activityHandler.GetActivitiesByThread)
		r.Get("/thread/{thread_id}/latest", activityHandler.GetLatestActivityByThread)
		r.Post("/receive", activityHandler.ReceiveActivity)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)
		
		r.Post("/", activityHandler.CreateActivity)                      
		r.Put("/{id}", activityHandler.UpdateActivity)                   
		r.Delete("/{id}", activityHandler.DeleteActivity)             
		
		r.Post("/thread", activityHandler.CreateActivityThread)  
		
		r.Get("/feature/{feature_uuid}", activityHandler.GetActivitiesByFeature)      
		r.Get("/phase/{phase_uuid}", activityHandler.GetActivitiesByPhase)            
		r.Get("/workspace/{workspace}", activityHandler.GetActivitiesByWorkspace)     

		r.Post("/{id}/actions", activityHandler.AddActivityActions)      
		r.Post("/{id}/questions", activityHandler.AddActivityQuestions) 
		r.Delete("/{id}/actions/{action_id}", activityHandler.RemoveActivityAction)    
		r.Delete("/{id}/questions/{question_id}", activityHandler.RemoveActivityQuestion) 
	})

	return r
} 