package item

type (
	Repository interface {
		Save(item Item) error
	}

	itemUseCase struct {
		repo Repository
	}
)

func NewItemUseCase(repo Repository) *itemUseCase {
	return &itemUseCase{repo: repo}
}

func (uc *itemUseCase) CreateItem(name string) error {
	i := Item{Name: name}
	return uc.repo.Save(i)
}
