package router

import "time"

const FilepathRoot = 	"."
const maxChirpLength = 	140
const maxTokenLifetime = 1 * time.Hour
const maxRefreshTokenLifetime = 60 * 24 * time.Hour
const templatePath = 	"./admin/index.html"
const port = 			"8080"
var badWords = 			map[string]struct{}{
							"kerfuffle": {},
							"sharbert": {},
							"fornax": {},
}