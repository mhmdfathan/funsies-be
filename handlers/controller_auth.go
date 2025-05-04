package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/mhmdfathan/funsies-be/config"
	dbmodels "github.com/mhmdfathan/funsies-be/models/db-models"
	"github.com/mhmdfathan/funsies-be/models/requests"
	"github.com/mhmdfathan/funsies-be/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest requests.RequestRegister

	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate request
	validate := validator.New()
	validate.RegisterValidation("phoneval", utils.ValidatePhone)
	err = validate.Struct(registerRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed vl",
			"error": err.Error(),
		})
		return
	}

	//encrypt password
	encryptedPassword, err := utils.EncryptPassword(registerRequest.Password)
	if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Registration failed",
			})
			return
	}

	//generate new UUID
	userId, err := uuid.NewV7()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}

	newUserData := dbmodels.User{
		ID: userId.String(),
		Email: registerRequest.Email,
		Username: registerRequest.Username,
		Password: encryptedPassword,
		Phone: registerRequest.Phone,
		FirstName: registerRequest.FirstName,
		LastName:	registerRequest.LastName,
		BirthDate:	registerRequest.BirthDate,
		Gender: registerRequest.Gender,
	}

	insertNewUser := config.DB.Create(&newUserData)
	if insertNewUser.Error != nil || insertNewUser.RowsAffected == 0 {
		//
		if insertNewUser.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_email\" (SQLSTATE 23505)" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Email already exists",
			})
			return
		}

		if insertNewUser.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_username\" (SQLSTATE 23505)" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Username already exists",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}


	//add token to activation token - not yet implemented
	//send email to activate - not yet implemented
		// response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration successful",
		})
}