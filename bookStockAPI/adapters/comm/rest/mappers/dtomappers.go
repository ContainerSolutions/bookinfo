package mappers

import (
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/comm/rest/dto"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/domain"
)

func MapBookStock2BookStockResponseDTO(p domain.BookStock) dto.BookStockResponseDTO {
	return dto.BookStockResponseDTO{
		ID:           p.ID,
		CurrentStock: p.CurrentStock,
	}
}
