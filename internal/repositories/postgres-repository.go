package repositories

type PostgresRepository struct {
}

type repository struct {
}

func CreatePostgresRepository() ProcessRepository {
	return &repository{}
}

func (r *repository) UpdateStatus() {

}
