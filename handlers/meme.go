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
)

func MemeImageUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
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
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to copy saved file", http.StatusInternalServerError)
		return
	}

	challenge := GetMemeChallenge()
	signer := SignChallenge(challenge.Challenge)
	mErr, mToken := GetMemeToken(challenge.Id, signer.Response.Sig)

	if mErr != "" {
		msg := "Could not get meme token"
		fmt.Println(msg, mErr)
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
		fmt.Println(msg)
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

func SignChallenge(challenge string) db.RelaySignerResponse {
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

	formData := url.Values{
		"id":     {id},
		"sig":    {sig},
		"pubkey": {config.RelayNodeKey},
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
	url := fmt.Sprintf("%s/public", config.MemeUrl)
	filePath := path.Join("./uploads", fileName)
	fileW, _ := os.Open(filePath)
	defer file.Close()

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)
	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	io.Copy(part, fileW)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, fileBody)
	req.Header.Set("Authorization", "BEARER "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)

	// Delete image from uploads folder
	DeleteFileFromUploadsFolder(filePath)

	if err != nil {
		fmt.Println("meme request Error ===", err)
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
		fmt.Println("The directory named", dirName, "does not exist")
		os.Mkdir(dirName, 0755)
	}
}
