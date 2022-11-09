package mappers

import (
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/dto"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
)

func MapBookInfo2BookInfoResponseDTO(p domain.BookInfo) dto.BookInfoResponseDTO {
	return dto.BookInfoResponseDTO{
		ID:     p.ID,
		Name:   p.Name,
		Author: p.Author,
	}
}
