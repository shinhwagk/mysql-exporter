package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

type Scraper interface {
	Name() string
	Help() string
	Version() float64
	Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error
}
