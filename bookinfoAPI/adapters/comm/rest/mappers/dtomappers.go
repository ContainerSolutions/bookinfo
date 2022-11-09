package mappers

import (
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/dto"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
)

func MapBookInfoRequestDTO2BookInfo(prd dto.BookInfoRequestDTO) domain.BookInfo {
	return domain.BookInfo{
		Name:                    prd.Name,
		BookInfoClass:           prd.BookInfoClass,
		Survived:                prd.Survived,
		Sex:                     prd.Sex,
		Age:                     prd.Age,
		SiblingsOrSpousesAboard: prd.SiblingsOrSpousesAboard,
		ParentsOrChildrenAboard: prd.ParentsOrChildrenAboard,
		Fare:                    prd.Fare,
	}
}

func MapBookInfo2BookInfoResponseDTO(p domain.BookInfo) dto.BookInfoResponseDTO {
	return dto.BookInfoResponseDTO{
		ID:                      p.ID,
		Name:                    p.Name,
		BookInfoClass:           p.BookInfoClass,
		Survived:                p.Survived,
		Sex:                     p.Sex,
		Age:                     p.Age,
		SiblingsOrSpousesAboard: p.SiblingsOrSpousesAboard,
		ParentsOrChildrenAboard: p.ParentsOrChildrenAboard,
		Fare:                    p.Fare,
	}
}
