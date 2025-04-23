package domain

type Category struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type CategoryRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

func (cr *CategoryRequest) ToDomain() *Category {
	return &Category{
		Name:        cr.Name,
		Description: cr.Description,
	}
}
