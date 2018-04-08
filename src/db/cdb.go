package db

var (
	CDB = NewChartDB()
)

type ChartDB interface {
	ListCharts(repo string) (map[string][]byte, error)
	RemoveChart(repo string, chart string, version *string) error
	GetChart(repo string, chart string, version *string) (map[string][]byte, error)
}

type etcdChartDB struct {
}

func NewChartDB() ChartDB {
	return &etcdChartDB{}
}

func (cb *etcdChartDB) ListCharts(repo string) (map[string][]byte, error) {
	return nil, nil
}

//移除
func (cb *etcdChartDB) RemoveChart(repo string, chart string, version *string) error {
	return nil
}

//这个要有能够通过前缀回去chart的功能
//返回chart的所有版本
func (cb *etcdChartDB) GetChart(repo string, chart string, version *string) (map[string][]byte, error) {
	return nil, nil
}
