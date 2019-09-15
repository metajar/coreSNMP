package controller

type CoreSNMPResource struct {
	DeviceName string `json:"deviceName" bson:"deviceName"`
	IP         string `json:"ip" bson:"ip"`
	Message    string `json:"message" bson:"message"`
}
