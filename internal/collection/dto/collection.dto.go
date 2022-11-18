package collection

import (
	catdto "nft/internal/category/dto"
	userdto "nft/internal/user/dto"
)

type Collection struct {
	ID          string               `json:"id,omitempty"`
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Categories  []catdto.CategoryDto `json:"categories,omitempty"`
	User        userdto.User         `json:"user,omitempty"`
	Status      string               `json:"status,omitempty"`
	HeaderImage string               `json:"header_image,omitempty"`
}

type CollectionList struct {
	Collections []Collection `json:"collections"`
}
