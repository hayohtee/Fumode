package main

import (
	"context"
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/validator"
	"net/http"
	"strconv"
	"time"
)

func (app *application) createFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with a 10 MB max memory limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	priceStr := r.Form.Get("price")
	stockStr := r.Form.Get("stock")
	category := r.Form.Get("category")

	v := validator.New()
	v.Check(name != "", "name", "must be provided")
	v.Check(description != "", "description", "must be provided")
	v.Check(category != "", "category", "must be provided")
	v.Check(priceStr != "", "price", "must be provided")
	v.Check(stockStr != "", "stock", "must be provided")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		v.AddError("price", "must be a valid number")
	}

	stock, err := strconv.ParseInt(stockStr, 10, 32)
	if err != nil {
		v.AddError("stock", "must be a valid number")
	}

	_, bannerHeader, err := r.FormFile("banner")
	if err != nil {
		v.AddError("banner", "must be a valid file")
	}

	images := r.MultipartForm.File["images"]
	if images == nil || len(images) == 0 {
		v.AddError("images", "must contain at least one image")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	furniture := data.Furniture{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       int(stock),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bannerUrl, err := app.s3Uploader.UploadImage(ctx, bannerHeader)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	imageUrls, err := app.s3Uploader.UploadImages(ctx, images)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	furniture.BannerURL = bannerUrl
	furniture.ImageURLs = imageUrls

	err = app.repositories.Furniture.Insert(&furniture)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"furniture": furniture}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
