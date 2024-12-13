package dao


type DaoAPI interface {
	CreateN(n interface{})
	GetNById(id string) interface{}
	DeleteNById(id string)
	UpdateN(n interface{})
	Connect()
}

var daoAPI DaoAPI

func SetDaoAPI(d DaoAPI) {
	daoAPI = d
}

func GetDaoAPI() DaoAPI {
	return daoAPI
}

func Init() {
	GetCurrentDaoAPI()
	daoAPI.Connect()
}

// 可通过传参初始化不同的domain
func GetCurrentDaoAPI() {
	if daoAPI == nil {
		daoAPI = &MysqlDao{}
	}
}

func CreateNews(n interface{}) {
	GetCurrentDaoAPI()
	daoAPI.CreateN(n)
}

func UpdateNews(n interface{}) {
	GetCurrentDaoAPI()
	daoAPI.UpdateN(n)
}

func DeleteNewsByid(id string) {
	GetCurrentDaoAPI()
	daoAPI.DeleteNById(id)
}

func getNewsById(id string)  {
	GetCurrentDaoAPI()
	daoAPI.GetNById(id)
}
