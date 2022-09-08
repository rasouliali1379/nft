package nftasset

import (
	catdto "nft/src/category/dto"
	userdto "nft/src/user/dto"
)

type Nft struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	//CollectionId *uuid.UUID
	Categories      []catdto.CategoryDto `json:"categories,omitempty"`
	User            userdto.UserDto      `json:"user,omitempty"`
	Status          string               `json:"status,omitempty"`
	NftImageUrl     string               `json:"nft_image_url,omitempty"`
	RejectionReason string               `json:"rejection_reason,omitempty"`
}

type NftList struct {
	Nfts []Nft `json:"nfts"`
}
