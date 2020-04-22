// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Scrape `SHOW GLOBAL STATUS`.

package collector

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Scrape query.
	statusAccountQuery = `select user, host, variable_name, variable_value FROM performance_schema.status_by_account WHERE user IS NOT NULL`
	// Subsystem.
	status = "status_account"
)

// Metric descriptors.
var (
	StatusAccountDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "status", "account"),
		"Total number of executed MySQL commands.",
		[]string{"user", "host", "name"}, nil,
	)
)

// ScrapeStatusAccount collects from `SHOW GLOBAL STATUS`.
type ScrapeStatusAccount struct{}

// Name of the Scraper. Should be unique.
func (ScrapeStatusAccount) Name() string {
	return status
}

// Help describes the role of the Scraper.
func (ScrapeStatusAccount) Help() string {
	return "Collect from SHOW GLOBAL STATUS"
}

// Version of MySQL from which scraper is available.
func (ScrapeStatusAccount) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeStatusAccount) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	globalStatusRows, err := db.QueryContext(ctx, statusAccountQuery)
	if err != nil {
		return err
	}
	defer globalStatusRows.Close()

	var user, host, key string
	var val uint64

	for globalStatusRows.Next() {
		if err := globalStatusRows.Scan(&user, &host, &key, &val); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(StatusAccountDesc, prometheus.CounterValue, float64(val), user, host, strings.ToLower(key))
	}

	return nil
}

// check interface
var _ Scraper = ScrapeStatusAccount{}
