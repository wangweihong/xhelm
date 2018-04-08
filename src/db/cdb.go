package db

var (
	CDB = NewChartDB()
)

type ChartDB interface {
	ListCharts(repo string) (map[string][]byte, error)
	RemoveChart(repo string, chart string) error
}

type etcdChartDB struct {
}

func NewChartDB() ChartDB {
	return &etcdChartDB{}
}

func (cb *etcdChartDB) ListCharts(repo string) (map[string][]byte, error) {
	return nil, nil
}

func (cb *etcdChartDB) RemoveChart(repo string, chart string) error {
	return nil
}
