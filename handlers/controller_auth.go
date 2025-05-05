package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
		//email or username already exists
		if insertNewUser.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_email\" (SQLSTATE 23505)" || insertNewUser.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_username\" (SQLSTATE 23505)" {	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Email or/and username already exists",
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
	activationToken := utils.GenerateActivationToken(32)

	tokenId, err := uuid.NewV7()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}

	token := dbmodels.ActivationToken{
		ID: tokenId.String(),
		UserID:    userId.String(),
		Token:     activationToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	createActivationToken := config.DB.Create(&token)
	if createActivationToken.Error != nil || createActivationToken.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}
	

	//send email to activate - not yet implemented
	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registration successful",
	})
}

func ActivateAccount (w http.ResponseWriter, r *http.Request) {
	//request query
	activationTokenRequest := r.URL.Query().Get("activation")
	if activationTokenRequest == "" {
		http.Error(w, "Activation token is required", http.StatusBadRequest)
		return
	}

	//find pending user
	var existingPendingUser requests.PendingUser
	result := config.DB.Raw(`
		SELECT users.id, users.is_active, activation_tokens.token, activation_tokens.expires_at
		FROM activation_tokens
		JOIN users ON activation_tokens.user_id = users.id
		WHERE activation_tokens.token = ?
		LIMIT 1
	`, activationTokenRequest).Scan(&existingPendingUser)

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Activation failed",
		})
		return
	}

	//if there's no pending user
	if existingPendingUser.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User doesn't exist",
		})
		return
	}

	//if the token already expired and user's account is not yet activated
	if existingPendingUser.ExpiresAt.Before(time.Now()) && !existingPendingUser.IsActive {
		deletePendingUser := config.DB.Where("id = ?", existingPendingUser.ID).Delete(&dbmodels.User{})

		if deletePendingUser.Error != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Internal server error",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Activation code already expired, please re-register",
		})
		return
	}

	//activate
	activateUser := config.DB.Model(&dbmodels.User{}).Where("id = ?", existingPendingUser.ID).Update("is_active", true)
	if activateUser.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	//delete token
	deleteToken := config.DB.Where("token = ?", existingPendingUser.Token).Delete(dbmodels.ActivationToken{})

	if deleteToken.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User activation successful",
	})
}