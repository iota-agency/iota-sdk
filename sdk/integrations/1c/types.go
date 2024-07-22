package client_1c

type OdataService struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type OdataServices struct {
	OdataMetadata string         `json:"odata.metadata"`
	Value         []OdataService `json:"value"`
}
