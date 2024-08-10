package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"icd_data_fetch/src/model"
	"net/http"
	"net/url"
)

type IcdRepo struct {
	ChapterUrl   string
	ClientId     string
	ClientSecret string
	Language     string
}

// getIcdToken returns access token
func (i IcdRepo) getIcdToken() (string, error) {
	fmt.Println("Fetching access token")
	response := make(map[string]interface{})

	params := url.Values{}
	params.Add("client_id", i.ClientId)
	params.Add("client_secret", i.ClientSecret)
	params.Add("scope", "icdapi_access")
	params.Add("grant_type", "client_credentials")

	payload := bytes.NewBufferString(params.Encode())

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "https://icdaccessmanagement.who.int/connect/token", payload)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	accessToken, ok := response["access_token"].(string)
	if !ok {
		customError := errors.New("access token not found")
		return "", customError
	}

	fmt.Println("Access token fetched")
	return accessToken, nil
}

// getRawIcdBlock returns Icd entity from block to category https://id.who.int/swagger/index.html
func (i IcdRepo) getRawIcdBlock(entity string, token string) (*map[string]interface{}, error) {
	var icdBlock map[string]interface{}

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", entity, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("API-Version", "v2")
	req.Header.Set("Accept-Language", i.Language)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&icdBlock)
	if err != nil {
		return nil, err
	}
	return &icdBlock, nil
}

// GetFullIcdChapter returns Icd entity with all of its children
func (i IcdRepo) GetFullIcdChapter() (*model.IcdBlock, error) {
	token, err := i.getIcdToken()
	if err != nil {
		return nil, err
	}

	icdBlock, err := i.getRawIcdBlock(i.ChapterUrl, token)
	if err != nil {
		return nil, err
	}

	chapterChidren := (*icdBlock)["child"].([]interface{})

	var chapter *model.IcdBlock
	data, err := json.Marshal(icdBlock)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &chapter)

	fmt.Printf("Fetching chapter: %s\n", chapter.Title.Value)

	for _, child := range chapterChidren {
		icdChild, err := i.getRawIcdBlock(child.(string), token)
		if err != nil {
			return nil, err
		}
		childBlock := model.IcdBlock{}
		data, err := json.Marshal(icdChild)
		json.Unmarshal(data, &childBlock)

		fmt.Printf("Fetching block: %s\n", childBlock.Title.Value)

		blockChildren, ok := (*icdChild)["child"].([]interface{})
		if ok {
			for _, blockChild := range blockChildren {
				icdBlockChild, err := i.getRawIcdBlock(blockChild.(string), token)
				if err != nil {
					return nil, err
				}
				categoryChild := model.IcdBlock{}
				data, err := json.Marshal(icdBlockChild)
				json.Unmarshal(data, &categoryChild)

				fmt.Printf("Fetching category: %s\n", categoryChild.Title.Value)

				blockChildChildren, ok := (*icdBlockChild)["child"].([]interface{})
				if ok {
					for _, blockChildChild := range blockChildChildren {
						icdBlockChildChild, err := i.getRawIcdBlock(blockChildChild.(string), token)
						if err != nil {
							return nil, err
						}
						subcategoryChild := model.IcdBlock{}
						data, err := json.Marshal(icdBlockChildChild)
						json.Unmarshal(data, &subcategoryChild)

						fmt.Printf("Fetching subcategory: %s\n", subcategoryChild.Title.Value)
						categoryChild.Children = append(categoryChild.Children, subcategoryChild)
					}
				}
				childBlock.Children = append(childBlock.Children, categoryChild)
			}
		}
		chapter.Children = append(chapter.Children, childBlock)
	}

	return chapter, nil
}
