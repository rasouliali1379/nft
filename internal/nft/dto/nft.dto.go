package nft

import (
	catdto "nft/internal/category/dto"
	userdto "nft/internal/user/dto"
)

type Nft struct {
	ID              string               `json:"id,omitempty"`
	Title           string               `json:"title,omitempty"`
	Description     string               `json:"description,omitempty"`
	Categories      []catdto.CategoryDto `json:"categories,omitempty"`
	User            userdto.User         `json:"user,omitempty"`
	Status          string               `json:"status,omitempty"`
	NftImageUrl     string               `json:"nft_image_url,omitempty"`
	RejectionReason string               `json:"rejection_reason,omitempty"`
}

type NftList struct {
	Nfts []Nft `json:"nfts"`
}
