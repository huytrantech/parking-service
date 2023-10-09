package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/sqlscan"
	"parking-service/model/database_model"
	"parking-service/provider/postgres_provider"
	"strings"
)

var ParkingSlotRepositoryCol = []string{"id", "parking_id", "parking_type",
	"price", "total_slot", "current_slot", "status", "metadata"}
type IParkingSlotRepository interface {
	InsertOneParkingSlot(ctx context.Context , model database_model.ParkingSlotModel) (int,error)
	FindManyParkingSlotWithParkingById(ctx context.Context , fields []string  , parkingId int) (
		data []database_model.ParkingSlotModel, err error)
	InsertManyParkingSlot(ctx context.Context , models []database_model.ParkingSlotModel) error
	FindOneParkingSlot(ctx context.Context , queryStr string, dataParams []interface{}) (
		data *database_model.ParkingSlotModel, err error)
	UpdateParkingSlot(ctx context.Context , queryStr string ,  dataParams []interface{}) error
}
type parkingSlotRepository struct {
	Db *sql.DB
}

func NewParkingSlotRepository(IPostgresProvider postgres_provider.IPostgresProvider) IParkingSlotRepository {
	return &parkingSlotRepository{Db: IPostgresProvider.GetDB()}
}

func (rp *parkingSlotRepository) InsertOneParkingSlot(ctx context.Context , model database_model.ParkingSlotModel) (int,error) {

	sqlQuery := `INSERT INTO public.parking_slot(parking_id, parking_type, price, total_slot, current_slot, status, metadata)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	var parkingIdSlot int
	err := rp.Db.QueryRow(sqlQuery, model.ParkingId, model.ParkingType, model.Price,
		model.TotalSlot, model.CurrentSlot, model.Status,model.Metadata).Scan(&parkingIdSlot)
	if err != nil {
		return 0, err
	}
	return parkingIdSlot, nil

}

func (rp *parkingSlotRepository) InsertManyParkingSlot(ctx context.Context , models []database_model.ParkingSlotModel) error {
	arrInsert := make([]string , len(models))
	arrData := make([]interface{} , 0)
	totalCount := 1
	if len(models) == 0 {
		return nil
	}
	for index , v := range models {
		arrInsert[index] = fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d)` ,
			totalCount,totalCount +1 , totalCount +2 , totalCount +3,totalCount +4,totalCount +5)
		arrData = append(arrData , v.ParkingId , v.ParkingType , v.Price ,
			v.TotalSlot , v.CurrentSlot , v.Status)
		totalCount += 6
	}

	sqlQuery := fmt.Sprintf(`INSERT INTO public.parking_slot(parking_id, parking_type, price, total_slot, current_slot, status) VALUES %s ;` , strings.Join(arrInsert , ","))

	_,err := rp.Db.Exec(sqlQuery, arrData...)
	if err != nil {
		return err
	}
	return  nil

}

func (rp *parkingSlotRepository) FindManyParkingSlotWithParkingById(ctx context.Context , fields []string  , parkingId int) (
	data []database_model.ParkingSlotModel, err error) {
	data = make([]database_model.ParkingSlotModel , 0)
	fieldsQuery := strings.Join(ParkingSlotRepositoryCol,",")
	if len(fields) > 0 {
		fieldsQuery = strings.Join(fields,",")
	}
	pgQuery := fmt.Sprintf("select %s from parking_slot where parking_id = $1", fieldsQuery)
	err = sqlscan.Select(ctx , rp.Db , &data , pgQuery , parkingId)
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

func (rp *parkingSlotRepository) FindOneParkingSlot(ctx context.Context , queryStr string, dataParams []interface{}) (
	data *database_model.ParkingSlotModel, err error) {
	data = new(database_model.ParkingSlotModel)
	err = sqlscan.Get(ctx , rp.Db , data , queryStr , dataParams...)
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

func(rp *parkingSlotRepository) UpdateParkingSlot(ctx context.Context , queryStr string ,  dataParams []interface{}) error {

	_,err := rp.Db.Exec(queryStr,dataParams...)
	if err != nil {
		return err
	}

	return nil
}