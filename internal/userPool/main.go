package userPool

type UserPool struct {
	userList         []User
	listTokenRequest [][]byte
}

type User struct {
	Token string
	Name  string
}

func (up *UserPool) initListTokenRequest() {
	for _, tokenData := range fixtureList {
		up.listTokenRequest = append(up.listTokenRequest, []byte(tokenData))
	}
}

var fixtureList = []string{
	`{
	"osVersion": "8.1.1",
	"model": "iPhone 8",
	"platform": "Iphone",
	"pushToken": "string",
	"locale": "ru",

	"applicationPackageName": "com.millcroft.inapp.sandbox",
	"applicationVersion": "1.0.0",
	"idfa": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00B91",
	"installId": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00123",

	"udid": "FF60EE70-1F11-4880-BDE0-F908F2B18F88"
}`,

	`{
	"osVersion": "8.1.1",
	"model": "iPhone 8",
	"platform": "Iphone",
	"pushToken": "string",
	"locale": "ru",

	"applicationPackageName": "com.millcroft.inapp.sandbox",
	"applicationVersion": "1.0.0",
	"idfa": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00B91",
	"installId": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00123",

	"udid": "FF60EE70-1F11-4880-BDE0-F908F2B18A22"
}`,

	`{
	"osVersion": "8.1.1",
	"model": "iPhone 8",
	"platform": "Iphone",
	"pushToken": "string",
	"locale": "ru",

	"applicationPackageName": "com.millcroft.inapp.sandbox",
	"applicationVersion": "1.0.0",
	"idfa": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00B91",
	"installId": "9FF4ACCE-AEBF-4393-A354-E1B1FBF00123",

	"udid": "FF60EE70-1F11-4880-BDE0-F908F2B18A88"
}`,
}
