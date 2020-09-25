package controllers

import (
	"customer/dbconnection"
	"customer/models"
	"customer/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// InsertCustomer details
func InsertCustomer(c *gin.Context) {

	customer := models.Customer{}

	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "fail",
			"data":   "invalid request body; " + err.Error(),
		})
		return
	}

	customer.CreatedDate = time.Now().UTC()
	db := dbconnection.Get()

	if _, err := db.Model(&customer).Insert(); err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data": gin.H{
				"customer": customer,
			},
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
	}
}

// GetCustomer Details by legalEntityID
func GetCustomer(c *gin.Context) {

	legalEntityID := utils.ParamID(c, "legalEntityID")

	customer := models.Customer{}
	db := dbconnection.Get()

	err := db.Model(&customer).Where("legal_entity_id = ?", int64(legalEntityID)).Select()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"customer": customer,
			},
		})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"status": "fail",
		"data": gin.H{
			"customer": nil,
		},
	})

}

// UpdateCustomer to modify
func UpdateCustomer(c *gin.Context) {

	legalEntityID := utils.ParamID(c, "legalEntityID")

	customer := models.Customer{}
	db := dbconnection.Get()
	err := db.Model(&customer).Where("legal_entity_id = ?", int64(legalEntityID)).Select()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"data":   "customer not found",
		})
		return
	}

	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "fail",
			"data":   "invalid request body; " + err.Error(),
		})
		return
	}

	_, err = db.Model(&customer).Column("bankruptcy_indicator_flag", "company_name", "first_name", "last_name", "legal_entity_stage", "legal_entity_type", "date_of_birth").Where("legal_entity_id=?", legalEntityID).Returning("*").Update()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"customer": customer,
			},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"data":   err.Error(),
	})

}

// RemoveCustomer to delete the customer
func RemoveCustomer(c *gin.Context) {

	legalEntityID := utils.ParamID(c, "legalEntityID")
	customer := models.Customer{}

	db := dbconnection.Get()

	if _, err := db.Model(&customer).Where("legal_entity_id=?", legalEntityID).Delete(); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "customer removed successfully",
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
	}

}

// SearchCustomer to filter the customer
func SearchCustomer(c *gin.Context) {

	customerRequest := models.CustomerRequestBody{}
	db := dbconnection.Get()

	if err := c.BindJSON(&customerRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "fail",
			"data":   "invalid request body; " + err.Error(),
		})
		return
	}

	queryStr := fmt.Sprintf(`SELECT * FROM customers where `)
	if customerRequest.LegalEntityID > 0 {
		legalEntityQuery := fmt.Sprintf("legal_entity_id = '%d'", customerRequest.LegalEntityID)
		queryStr = queryStr + legalEntityQuery
	}

	if customerRequest.CompanyName != "" {
		companyQuery := fmt.Sprintf("AND company_name = '%v'", customerRequest.CompanyName)
		queryStr = queryStr + companyQuery
	}

	if customerRequest.FirstName != "" {
		firstNameQuery := fmt.Sprintf("AND first_name = '%v'", customerRequest.FirstName)
		queryStr = queryStr + firstNameQuery
	}

	if customerRequest.LastName != "" {
		lastNameQuery := fmt.Sprintf("AND last_name = '%v'", customerRequest.LastName)
		queryStr = queryStr + lastNameQuery
	}

	customerIns := []models.Customer{}
	_, err := db.Query(&customerIns, queryStr)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"customer": customerIns,
			},
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
	}
}
