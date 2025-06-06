package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

// MemeImageUpload godoc
//
//	@Summary		Meme Image Upload
//	@Description	Upload an image for a meme
//	@Tags			Memes
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			file	formData	file	true	"Meme Image File"
//	@Success		200		{string}	string	"Meme image URL"
//	@Router			/meme_upload [post]
func MemeImageUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dirName := "uploads"
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if uploads directory exists or create it
	CreateUploadsDirectory(dirName)

	// Saving the file
	dst, err := os.Create(dirName + "/" + header.Filename)

	if err != nil {
		panic("Unable to create file")
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		panic("Unable to copy saved file")
		return
	}

	challenge := GetMemeChallenge()
	signer := SignChallenge(challenge.Challenge)
	mErr, mToken := GetMemeToken(challenge.Id, signer.Response.Sig)

	if mErr != "" {
		msg := "Could not get meme token"
		logger.Log.Error("%s: %s", msg, mErr)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(msg)
	} else {
		err, memeImgUrl := UploadMemeImage(file, mToken.Token, header.Filename)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(memeImgUrl)
			return
		}

		msg := "Could not get meme image"
		logger.Log.Error("%s", msg)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(msg)
	}
}

func GetMemeChallenge() db.MemeChallenge {
	memeChallenge := db.MemeChallenge{}

	url := fmt.Sprintf("%s/ask", config.MemeUrl)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	// Unmarshal result
	err = json.Unmarshal(body, &memeChallenge)

	if err != nil {
		log.Printf("Reading Invoice body failed: %s", err)
	}

	return memeChallenge
}

type SignChallengeBody struct {
	Message string `json:"message"`
}

type V2SignChallengeResponse struct {
	Sig string `json:"sig"`
}

func SignChallenge(challenge string) db.RelaySignerResponse {

	if config.V2BotUrl != "" {
		url := fmt.Sprintf("%s/sign_base64", config.V2BotUrl)
		client := &http.Client{}

		challengeBody := SignChallengeBody{
			Message: challenge,
		}

		jsonBody, _ := json.Marshal(challengeBody)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("x-admin-token", config.V2BotToken)
		req.Header.Set("Content-Type", "application/json")
		res, _ := client.Do(req)

		if err != nil {
			log.Printf("[Sign Challenge for V2] Request Failed: %s", err)
			return db.RelaySignerResponse{
				Success:  false,
				Response: db.SignerResponse(db.SignerResponse{Sig: ""}),
			}
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Printf("[Sign Challenge for V2] Reading sign challenge response body failed: %s", err)
			return db.RelaySignerResponse{
				Success:  false,
				Response: db.SignerResponse(db.SignerResponse{Sig: ""}),
			}
		}

		v2SignChallengeResponse := V2SignChallengeResponse{}

		// Unmarshal result
		err = json.Unmarshal(body, &v2SignChallengeResponse)

		if err != nil {
			log.Printf("[Sign Challenge for V2] Unmarshalling response body failed: %s", err)
			return db.RelaySignerResponse{
				Success:  false,
				Response: db.SignerResponse(db.SignerResponse{Sig: ""}),
			}
		}

		signerResponse := db.RelaySignerResponse{
			Success:  true,
			Response: db.SignerResponse(v2SignChallengeResponse),
		}

		return signerResponse

	}

	url := fmt.Sprintf("%s/signer/%s", config.RelayUrl, challenge)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	signerResponse := db.RelaySignerResponse{}

	// Unmarshal result
	err = json.Unmarshal(body, &signerResponse)

	if err != nil {
		log.Printf("Reading Challenge body failed: %s", err)
	}

	return signerResponse
}

func GetMemeToken(id string, sig string) (string, db.MemeTokenSuccess) {
	memeUrl := fmt.Sprintf("%s/verify", config.MemeUrl)

	var pubkey string

	if config.V2BotUrl != "" {
		pubkey = config.GetV2ContactKey()
	} else {
		pubkey = config.RelayNodeKey
	}

	formData := url.Values{
		"id":     {id},
		"sig":    {sig},
		"pubkey": {pubkey},
	}

	res, err := http.PostForm(memeUrl, formData)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return "", db.MemeTokenSuccess{}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if res.StatusCode == 200 {
		tokenSuccess := db.MemeTokenSuccess{}

		// Unmarshal result
		err = json.Unmarshal(body, &tokenSuccess)

		if err != nil {
			log.Printf("Reading token success body failed: %s", err)
		}

		return "", tokenSuccess
	} else {
		var tokenError string

		// Unmarshal result
		err = json.Unmarshal(body, &tokenError)

		if err != nil {
			log.Printf("Reading token error body failed: %s %d", err, res.StatusCode)
		}

		return tokenError, db.MemeTokenSuccess{}
	}
}

func UploadMemeImage(file multipart.File, token string, fileName string) (error, string) {

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		logger.Log.Error("Failed to create uploads directory: %v", err)
		return err, ""
	}
	
	filePath := path.Join("./uploads", fileName)
	tempFile, err := os.Create(filePath)
	if err != nil {
		logger.Log.Error("Failed to create temporary file: %v", err)
		return err, ""
	}
	
	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		os.Remove(filePath)
		logger.Log.Error("Failed to write to temporary file: %v", err)
		return err, ""
	}
	tempFile.Close()
	
	url := fmt.Sprintf("%s/public", config.MemeUrl)
	
	fileW, err := os.Open(filePath)
	if err != nil {
		os.Remove(filePath)
		logger.Log.Error("Failed to open temporary file: %v", err)
		return err, ""
	}
	defer fileW.Close()

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		os.Remove(filePath)
		logger.Log.Error("Failed to create form file: %v", err)
		return err, ""
	}
	
	if _, err := io.Copy(part, fileW); err != nil {
		os.Remove(filePath)
		logger.Log.Error("Failed to copy file to form: %v", err)
		return err, ""
	}
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, fileBody)
	if err != nil {
		os.Remove(filePath)
		logger.Log.Error("Failed to create request: %v", err)
		return err, ""
	}
	
	req.Header.Set("Authorization", "BEARER "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)

	// Delete image from uploads folder
	os.Remove(filePath)

	if err != nil {
		logger.Log.Error("Meme request Error: %v", err)
		return err, ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err == nil {
		memeSuccess := db.Meme{}
		// Unmarshal result
		err = json.Unmarshal(body, &memeSuccess)
		if err != nil {
			log.Printf("Reading meme error body failed: %s", err)
		} else {
			return nil, config.MemeUrl + "/public/" + memeSuccess.Muid
		}
	}

	return err, ""
}

func DeleteFileFromUploadsFolder(filePath string) {
	e := os.Remove(filePath)
	if e != nil {
		log.Printf("Could not delete Image %s %s", filePath, e)
	}
}

func CreateUploadsDirectory(dirName string) {
	if _, err := os.Open(dirName); os.IsNotExist(err) {
		logger.Log.Info("The directory named %s does not exist", dirName)
		os.Mkdir(dirName, 0755)
	}
}
