package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sofiukl/io-service/api"
	_ "github.com/sofiukl/io-service/docs"
	"github.com/sofiukl/io-service/utils"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App - Application
type App struct {
	Router *mux.Router
	DB     *mongo.Database
	Config utils.Config
}

// Initialize - This function initializes the application
func (a *App) Initialize() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connectionString :=
		fmt.Sprintf("mongodb://%s:%s/%s", config.DBHost, config.DBPort, config.DBName)
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	dbConnectMsg := fmt.Sprintf("Connected to DB %s", connectionString)
	fmt.Println(dbConnectMsg)
	a.DB = client.Database(config.DBName)
	a.Router = mux.NewRouter()
	a.Config = config
	a.initializeRoutes()
}

// Run - This functio funs the application
func (a *App) Run(address string) {
	fmt.Println("Application is running on port", address)
	if err := http.ListenAndServe(address, a.Router); err != nil {
		log.Fatal(err)
	}
}

// @title io Service API Documentation (Golang)
// @version 1.0
// @description Serice for io operations
// @schemes http https
// @host localhost:8082
// @BasePath /
func (a *App) initializeRoutes() {
	a.Router.PathPrefix("/io/api/v1/swagger").Handler(httpSwagger.WrapHandler)
	s := a.Router.PathPrefix("/io/api/v1").Subrouter()
	s.HandleFunc("/upload", a.uploadFile).Methods("POST")
	s.HandleFunc("/download/{fileKey}", a.downloadFile).Methods("GET")
}

// uploadFile godoc
// @Summary Uploaded file
// @Description This api will be used to upload file
// @Tags file-upload
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.GenericResponse
// @Failure 500 {object} models.GenericResponse
// @Param file formData file true "File to be uploaded"
// @Router /io/api/v1/upload [post]
func (a *App) uploadFile(w http.ResponseWriter, r *http.Request) {
	config := a.Config
	DB := a.DB
	api.UploadFileS3(DB, config, w, r)
}

func (a *App) downloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileKey := vars["fileKey"]

	if len(fileKey) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "File Key is mandatory", "")
	}
	config := a.Config
	api.GetDownloadURL(fileKey, config, w, r)
}
