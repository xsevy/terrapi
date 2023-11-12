package helpers

type resourceIDs struct {
	CreateAppSyncDataSource string
	CreateAppSyncAPI        string
}

var ResourceIDs = resourceIDs{
	CreateAppSyncDataSource: "create_app_sync_data_source",
	CreateAppSyncAPI:        "create_app_sync_api",
}

var ResourceNames = map[string]string{
	ResourceIDs.CreateAppSyncDataSource: "Create data source",
	ResourceIDs.CreateAppSyncAPI:        "Create API",
}
