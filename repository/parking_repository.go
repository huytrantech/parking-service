package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/spf13/cast"
	"parking-service/model/database_model"
	"parking-service/provider/postgres_provider"
	"strings"
)

var ParkingRepositoryCol = []string{"id", "public_id", "owner_name",
	"parking_name", "owner_phone", "parking_phone", "address", "status", "open_at", "close_at", "images"}

const parkingTableName = "parking"

type IParkingRepository interface {
	CreateParking(ctx context.Context, model database_model.ParkingModel) (int, error)
	FindManyParkingWithIds(ctx context.Context, parkingIds []int, fields []string) (
		data []database_model.ParkingModel, err error)
	CountParking(ctx context.Context, model database_model.ParkingQueryModel) (
		data int, err error)
	FindOneParkingById(ctx context.Context, fields []string, parkingId int) (
		data *database_model.ParkingModel, err error)
	FindOne(ctx context.Context, dto database_model.ParkingQueryModel) (
		data *database_model.ParkingModel, err error)
	UpdateOne(ctx context.Context, dto database_model.ParkingQueryModel) error
	FindMany(ctx context.Context, dto database_model.ParkingQueryModel) (
		data []database_model.ParkingModel, err error)
}

type parkingRepository struct {
	Db *sql.DB
}

func NewParkingRepository(IPostgresProvider postgres_provider.IPostgresProvider) IParkingRepository {
	return &parkingRepository{Db: IPostgresProvider.GetDB()}
}

func (rp *parkingRepository) CreateParking(ctx context.Context, model database_model.ParkingModel) (int, error) {
	sqlQuery := `INSERT INTO public.parking(
	public_id, owner_name, parking_name, owner_phone, parking_phone, 
	address, status, created_date, created_name,open_at,close_at,images,parking_types)
	VALUES ($1, $2, $3, $4, $5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id;`
	var parkingId int
	err := rp.Db.QueryRow(sqlQuery, model.PublicId, model.OwnerName, model.ParkingName, model.OwnerPhone,
		model.ParkingPhone, model.Address, model.Status, model.CreatedDate, model.CreatedName,
		model.OpenAt, model.CloseAt, model.Images, model.ParkingTypes).Scan(&parkingId)
	if err != nil {
		return 0, err
	}
	return parkingId, nil
}

func (rp *parkingRepository) CountParking(ctx context.Context, model database_model.ParkingQueryModel) (
	data int, err error) {
	conditions, params := model.ToFilter()
	if len(conditions) == 0 {
		conditions = "1=1"
	}
	sqlQuery := fmt.Sprintf("select count(*) from %s where %s", parkingTableName, conditions)

	err = sqlscan.Get(ctx, rp.Db, &data, sqlQuery, params...)
	if err != nil {
		return
	}

	return
}

func (rp *parkingRepository) FindManyParkingWithIds(ctx context.Context, parkingIds []int, fields []string) (
	data []database_model.ParkingModel, err error) {

	arrStr := make([]string, 0)
	for _, value := range parkingIds {
		arrStr = append(arrStr, cast.ToString(value))
	}
	err = sqlscan.Select(ctx, rp.Db, &data, fmt.Sprintf(`select %s from parking where id in (%s)`, strings.Join(fields, ","),
		strings.Join(arrStr, ",")))
	if err != nil {
		return
	}

	return
}

func (rp *parkingRepository) FindOneParkingById(ctx context.Context, fields []string, parkingId int) (
	data *database_model.ParkingModel, err error) {
	data = new(database_model.ParkingModel)
	fieldsQuery := strings.Join(ParkingRepositoryCol, ",")
	if len(fields) > 0 {
		fieldsQuery = strings.Join(fields, ",")
	}
	pgQuery := fmt.Sprintf("select %s from parking where id = $1", fieldsQuery)
	err = sqlscan.Get(ctx, rp.Db, data, pgQuery, parkingId)
	if sqlscan.NotFound(err) {
		err = nil
		data = nil
		return
	}
	if err != nil {
		return
	}
	return
}

func (rp *parkingRepository) FindOne(ctx context.Context, dto database_model.ParkingQueryModel) (
	data *database_model.ParkingModel, err error) {
	data = new(database_model.ParkingModel)
	conditionStr, params := dto.ToFilter()
	if len(conditionStr) == 0 {
		return
	}
	if len(dto.Fields) == 0 {
		dto.Fields = ParkingRepositoryCol
	}
	sqlString := fmt.Sprintf("select %s from %s where %s limit 1",
		strings.Join(dto.Fields, ","), parkingTableName, conditionStr)
	err = sqlscan.Get(ctx, rp.Db, data, sqlString, params...)
	if sqlscan.NotFound(err) {
		err = nil
		data = nil
		return
	}
	if err != nil {
		return
	}
	return
}

func (rp *parkingRepository) FindMany(ctx context.Context, dto database_model.ParkingQueryModel) (
	data []database_model.ParkingModel, err error) {

	conditionStr, params := dto.ToFilter()
	if len(conditionStr) == 0 {
		conditionStr = "1=1"
	}
	if len(dto.Fields) == 0 {
		dto.Fields = ParkingRepositoryCol
	}
	if dto.Limit <= 0 {
		dto.Limit = 10
	}
	if dto.Offset < 0 {
		dto.Offset = 0
	}
	sqlString := fmt.Sprintf("select %s from %s where %s limit %d offset %d",
		strings.Join(dto.Fields, ","), parkingTableName, conditionStr, dto.Limit, dto.Offset)

	err = sqlscan.Select(ctx, rp.Db, &data, sqlString, params...)
	if err != nil {
		return
	}

	return
}

func (rp *parkingRepository) UpdateOne(ctx context.Context, dto database_model.ParkingQueryModel) error {

	filter, update, params := dto.ToFilterUpdate()
	if len(filter) == 0 || len(update) == 0 {
		return nil
	}
	queryString := fmt.Sprintf("update %s set %s where %s", parkingTableName, update, filter)
	_, err := rp.Db.Exec(queryString, params...)
	if err != nil {
		return err
	}

	return nil
}
