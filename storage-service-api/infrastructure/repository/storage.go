package repository

import (
	"encoding/json"
	"github.com/jackc/pgx/v4"
	st "storage.g128.rvlz.io/domain/storage"
)

type WarehouseStorageDTO struct {
	warehouseID        string
	warehouseStorageID string
	capacity           int
	claimedStorage     int
	files              []struct {
		ID   string `json:"id"`
		Size int    `json:"size"`
	}
}

type WarehouseStoragePO struct {
	id               string
	warehouseID      string
	capacity         int
	claimedStorage   int
	fileReservations string
}

type WarehouseStorageRepository struct {
	context DBContext
}

func (repo *WarehouseStorageRepository) getDTO(row pgx.Row) (WarehouseStorageDTO, error) {
	var dto WarehouseStorageDTO
	err := row.Scan(&dto.warehouseStorageID, &dto.warehouseID, &dto.capacity, &dto.claimedStorage, &dto.files)
	if err != nil {
		return WarehouseStorageDTO{}, err
	}
	return dto, nil
}

func (WarehouseStorageRepository) mapDTOToAggregate(dto WarehouseStorageDTO) *st.WarehouseStorage {
	warehouseStorageBuilder := st.StorageBuilder{}.NewWarehouseStorageBuilder()
	warehouseStorageBuilder.
		SetID(dto.warehouseStorageID).
		SetWarehouseID(dto.warehouseID).
		SetCapacity(dto.capacity).
		SetClaimedStorage(dto.claimedStorage)
	for _, file := range dto.files {
		warehouseStorageBuilder.AddFileReservation(file.ID, file.Size)
	}
	return warehouseStorageBuilder.GetWarehouseStorage()
}

func (WarehouseStorageRepository) mapAggregateToPO(aggregate *st.WarehouseStorage) WarehouseStoragePO {
	type FileReservation struct {
		ID   string `json: "id"`
		Size int    `json: "size"`
	}
	var fileReservations []FileReservation
	warehouseStoragePO := WarehouseStoragePO{}
	warehouseStorageProxy := st.NewWarehouseStorageProxy(aggregate)
	warehouseStoragePO.id = warehouseStorageProxy.GetID()
	warehouseStoragePO.warehouseID = warehouseStorageProxy.GetWarehouseID()
	warehouseStoragePO.capacity = warehouseStorageProxy.GetCapacity()
	warehouseStoragePO.claimedStorage = warehouseStorageProxy.GetClaimedStorage()
	for _, fr := range warehouseStorageProxy.GetFileReservations() {
		fileReservations = append(fileReservations, FileReservation{ID: fr.ID, Size: fr.Size})
	}
	fileReservationsData, _ := json.Marshal(fileReservations)
	warehouseStoragePO.fileReservations = string(fileReservationsData)
	return warehouseStoragePO
}

func (repo *WarehouseStorageRepository) GetByWarehouseID(id string) (*st.WarehouseStorage, error) {
	statementName := "get_warehouse_storage_by_id"
	statement := `
	SELECT
	id,
	warehouse_id,
	capacity,
	claimed_storage,
	file_reservations
	FROM warehouse_storage
	WHERE warehouse_id = $1 LIMIT 1
	`
	err := repo.context.Prepare(statementName, statement)
	if err != nil {
		return nil, RepositoryErr
	}
	row := repo.context.QueryOne(statementName, id)
	warehouseDTO, err := repo.getDTO(row)
	if err == pgx.ErrNoRows {
		return nil, NotFoundErr
	} else if err != nil {
		return nil, RepositoryErr
	}
	warehouseStorage := repo.mapDTOToAggregate(warehouseDTO)
	return warehouseStorage, nil
}

func (repo *WarehouseStorageRepository) Save(warehouseStorage *st.WarehouseStorage) error {
	warehouseStoragePO := repo.mapAggregateToPO(warehouseStorage)
	statementName := "update_warehouse_storage"
	statement := `
	UPDATE warehouse_storage SET
	capacity = $1,
	claimed_storage = $2,
	file_reservations = $3
	WHERE warehouse_id = $4
	`
	err := repo.context.Prepare(statementName, statement)
	if err != nil {
		return RepositoryErr
	}
	err = repo.context.Execute(
		statementName,
		warehouseStoragePO.capacity,
		warehouseStoragePO.claimedStorage,
		warehouseStoragePO.fileReservations,
		warehouseStoragePO.warehouseID,
	)
	if err != nil {
		return RepositoryErr
	}
	return nil
}

func NewWarehouseStorageRepository(context DBContext) *WarehouseStorageRepository {
	return &WarehouseStorageRepository{context: context}
}
