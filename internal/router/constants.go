package router

const FilepathRoot = 	"."
const maxChirpLength = 	140
const templatePath = 	"./admin/index.html"
const port = 			"8080"
var badWords = 			map[string]struct{}{
							"kerfuffle": {},
							"sharbert": {},
							"fornax": {},
}