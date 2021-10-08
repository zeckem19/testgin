package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zeckem19/testgin/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// create a validator objectc c
var validate = validator.New()

// patient represents data about a Patient
type Patient struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `json:"name" validate:"required,min=2,max=100"`
	Age           int                `json:"age" validate:"required"`
	Gender        string             `json:"gender" validate:"required,min=1,max=1"`
	Creation_time time.Time          `json:"creation_time"`
	End_time      time.Time          `json:"end_time"`
	Patient_id    string             `json:"patient_id"`
	Kiosk_id      string             `json:"kiosk_id"`
	Doctor_id     string             `json:"doctor_id"`
}

// albums slice to seed record album data.
// var patients = []Patient{
// 	{name: "john", age: 30, gender: "M"},
// 	{name: "jane", age: 20, gender: "M"},
// 	{name: "jule", age: 25, gender: "M"},
// }

//connect to to the database and open a food collection

var patientCollection *mongo.Collection = database.OpenCollection(database.Client, "patient")
var kioskCollection *mongo.Collection = database.OpenCollection(database.Client, "kiosk")

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.New()
	r.Use(gin.Logger())

	r.GET("/patients", getPatients)
	// r.GET("/patients/:id", getPatientsByID)
	r.POST("/register/:kiosk_id", registerPatient)

	r.Run("0.0.0.0:" + port)
}

// getAlbums responds with the list of all patients as JSON.
func getPatients(c *gin.Context) {
	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := patientCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Database read failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var patients []Patient
	if err = cursor.All(ctx, &patients); err != nil {
		fmt.Println("Database fetch failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, patients)
}

// postAlbums adds an album from JSON received in the request body.
func registerPatient(c *gin.Context) {
	kiosk_id := c.Param("kiosk_id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var newPatient Patient

	// Call BindJSON to bind the received JSON to
	// bind the object that comes in with the declared varaible.
	//  thrrow an error if one occurs
	if err := c.BindJSON(&newPatient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("%+v\n", newPatient)
	// use the validation packge to verify that all items coming in meet the requirements of the struct
	validationErr := validate.Struct(newPatient)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	newPatient.Creation_time, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newPatient.End_time = time.Time{}

	//generate new ID for the object to be created
	newPatient.ID = primitive.NewObjectID()

	// assign the the auto generated ID to the primary key attribute
	newPatient.Patient_id = newPatient.ID.Hex()
	newPatient.Kiosk_id = kiosk_id

	// finally, insert to database
	result, insertErr := patientCollection.InsertOne(ctx, newPatient)
	if insertErr != nil {
		msg := fmt.Sprintf("Patient not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	filter := bson.D{{"id", kiosk_id}}
	update := bson.D{
		{"$push", bson.D{
			{"queue", newPatient},
		}},
	}
	if result != nil {
		kioskResult, updateErr := kioskCollection.UpdateOne(ctx, filter, update)
		if updateErr != nil {
			msg := fmt.Sprintf("Patient not queued")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if kioskResult != nil {
			// trigger check_kiosk_empty()
			_, err := http.Post("medlyves-fastapi:8000/kiosk/check/K00001", "application/json", nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			} else {
				c.IndentedJSON(http.StatusCreated, kioskResult)
			}
		}

	}
	c.IndentedJSON(http.StatusInternalServerError, "Unknown error - possibly no such kiosk")

}
