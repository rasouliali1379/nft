package sale

import (
	"github.com/google/uuid"
	"log"
	"nft/internal/sale/model"
	"testing"
)

func TestMapSaleModelToEntity(t *testing.T) {
	saleEntity := mapSaleModelToEntity(model.Sale{AssetId: uuid.New(), MinPrice: 42})
	log.Println(saleEntity.SaleType, saleEntity.AssetType)
}
