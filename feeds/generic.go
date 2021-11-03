package feeds

import (
	"io/ioutil"
	"net/http"
)

func httpget(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return body, nil
}
