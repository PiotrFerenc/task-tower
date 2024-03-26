package repositories

type repository struct {
}

func CreateInMemoryRepository() ProcessRepository {
	return &repository{}
}

func (r *repository) UpdateStatus() {

}
