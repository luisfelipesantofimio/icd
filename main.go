package main

import (
	"fmt"
	"icd_data_fetch/src/config"
	"icd_data_fetch/src/model"
	"icd_data_fetch/src/repository"
)

func main() {
	config := config.GetConfig("./config.json")

	icdEsRepo := repository.IcdRepo{
		ChapterUrl:   "http://id.who.int/icd/release/11/2023-01/mms/334423054",
		ClientId:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Language:     "es",
	}

	icdEnRepo := repository.IcdRepo{
		ChapterUrl:   "http://id.who.int/icd/release/11/2023-01/mms/334423054",
		ClientId:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Language:     "en",
	}

	esIcdBlock, err := icdEsRepo.GetFullIcdChapter()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = model.ExportToJSON(esIcdBlock, icdEsRepo.Language)
	if err != nil {
		fmt.Println(err)
		return
	}

	enIcdBlock, err := icdEnRepo.GetFullIcdChapter()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = model.ExportToJSON(enIcdBlock, icdEnRepo.Language)
	if err != nil {
		fmt.Println(err)
		return
	}
}
