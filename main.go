package main

// NOTE : Cet exercice a été réalisé à partir de la vidéo suivante :
// https://www.youtube.com/watch?v=bj77B59nkTQ

// Le tutoriel propose d'utiliser le framework Gin pour créer une API REST en Go
// Commande BASH : go get -u github.com/gin-gonic/gin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct{
	ID			string 	`json:"id"` 	// Les noms des champs commencent par une majuscule et peuvent ainsi être exportés
	Title		string	`json:"title"`	// En revanche, dans le fichier JSON, on les affiche en minuscules
	Author		string	`json:"author"`
	Quantity	int		`json:"quantity"`
}

var books = []Book{ // books est une slice de la struct "Book"
	{ID: "1", Title: "Harry Potter", Author: "J.K. Rowling", Quantity: 7},
	{ID: "2", Title: "Le Seigneur des Anneaux", Author: "J.R.R. Tolkien", Quantity: 3},
	{ID: "3", Title: "Le Trône de Fer", Author: "George R.R. Martin", Quantity: 5}, 
}

//========================================================
// MÉTHODES -- 👁️ -- GET & GET BY ID
//========================================================

func getBooks(context *gin.Context) {
	// Les méthodes utilisées par le routeur Gin n'ont pas de "return"
	// En revanche, elles peuvent renvoyer une réponse au client au format json
	// Cela se fait de la manière suivante :
	context.IndentedJSON(http.StatusOK, books)
}

// Cette méthode va chercher un livre par son ID
// Elle renvoie un pointeur vers le bon livre (et une erreur)
// Elle ne retourne rien au format JSON

func getBookById(id string) (*Book, error) {
	for index, book := range books {
		if book.ID == id {
			return &books[index], nil
		}
	}
	return nil, errors.New("Book not found")
}

// La méthode suivante sera associée à la route GetBookById
// L'id est passé en paramètre de la route, donc il n'a pas besoin d'être en paramètre de la fonction
func routerGetBookById(context *gin.Context) {
	id := context.Param("id") // récupération de l'ID en paramètre de la route

	book, err := getBookById(id) // utilisation de la méthode créée ci-dessus

	if err != nil { // Cas où on ne trouve pas de livre associé à l'ID
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Book with id %s not found", id)})
		return
	}

	// Envoi du livre au format JSON avec un statut 200
	context.IndentedJSON(http.StatusOK, book)
}

//========================================================
// MÉTHODE -- 🖍️ -- CREATE
//========================================================

func createBook(context *gin.Context) {
	var newBook Book // Création d'une variable (instance de la struct "Book")

	if err := context.BindJSON(&newBook); err != nil {
		return // Dans le cas où il y a une erreur, on s'arrête là
	}

	books = append(books, newBook)	// S'il n'y a pas d'erreur avec le nouveau livre ajouté en body (JSON), on l'ajoute à la slice
	context.IndentedJSON(http.StatusCreated, newBook)	// Puis on renvoie le livre au format JSON avec un statut 201
}

//========================================================
// MÉTHODES -- ➖ / ➕ -- CHECKOUT AND RETURN
//========================================================

func checkoutBookById(context *gin.Context) {
	id := context.Param("id")

	book, err := getBookById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Book with id %s not found", id)})
		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("No more of this book")})
		return
	}

	book.Quantity -= 1
	context.IndentedJSON(http.StatusOK, book)
}

func returnBookById(context *gin.Context) {
	id := context.Param("id")
	
	book, err := getBookById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Book with id %s not found", id)})
		return
	}

	book.Quantity += 1
	
	context.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default() // Routeur de la librairie Gin

	// 👁️📚 GET BOOKS & 👁️📗 GET BOOK BY ID
	router.GET("/books", getBooks)
	router.GET("/books/:id", routerGetBookById)
	
	// 📩📕 POST BOOK
	router.POST("/books", createBook)
	// NOTE : Commande bash pour faire une requête POST à partir des données dans le fichier body.json
	// curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"

	// CHECKOUT and RETURN
	router.PATCH("/checkout/:id", checkoutBookById)
	router.PATCH("/return/:id", returnBookById)

	router.Run("localhost:8080") // Attention à garder cette ligne après toutes les routes
}