package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	infr "storage.g128.rvlz.io/infrastructure"
	repo "storage.g128.rvlz.io/infrastructure/repository"
	// st "storage.g128.rvlz.io/domain/storage"
	// dm "storage.g128.rvlz.io/domain"
)

func main() {
	conn, err := pgx.Connect(
		context.Background(),
		"postgresql://storage_service_api:password@storage-service-db/storage_service_dev",
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())
	tx, _ := conn.Begin(context.Background())
	ctx := infr.NewContext(tx)
	warehouseStorageRepository := repo.NewWarehouseStorageRepository(ctx)
	warehouseStorage, err := warehouseStorageRepository.GetByWarehouseID("1")
	fmt.Println(warehouseStorage.AvailableStorage())
	fmt.Println(err)
	// warehouseStorage.Reserve(st.File{ID: dm.IDConstructor{}.NewFileID(), Size: 40})
	fmt.Println(warehouseStorageRepository.Save(warehouseStorage))
	ctx.SaveChanges()
}
