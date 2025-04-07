package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServicebResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type CepInputDTO struct {
	CEP string `json:"cep" binding:"required,min=8"`
}

type CepOutputDTO struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func cepHandler(c *gin.Context) {

	var cepInputDTO CepInputDTO

	if err := c.ShouldBindJSON(&cepInputDTO); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	url := fmt.Sprintf("http://service-b:8081/temperatura?cep=%s", cepInputDTO.CEP)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	var servicebResponse ServicebResponse

	if err := json.NewDecoder(resp.Body).Decode(&servicebResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	cepOutputDTO := CepOutputDTO(servicebResponse)

	c.JSON(http.StatusOK, cepOutputDTO)
}

func main() {

	router := gin.Default()
	router.POST("/", cepHandler)
	router.Run(":8080")

}
