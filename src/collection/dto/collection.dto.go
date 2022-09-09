package collection

import (
	catdto "nft/src/category/dto"
	userdto "nft/src/user/dto"
)

type Collection struct {
	ID          string               `json:"id,omitempty"`
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Categories  []catdto.CategoryDto `json:"categories,omitempty"`
	User        userdto.UserDto      `json:"user,omitempty"`
	Status      string               `json:"status,omitempty"`
	HeaderImage string               `json:"header_image,omitempty"`
}

type CollectionList struct {
	Collections []Collection `json:"collections"`
}
