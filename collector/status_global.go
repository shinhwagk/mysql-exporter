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
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Scrape query.
	globalStatusQuery = `SHOW GLOBAL STATUS`
	// Subsystem.
	globalStatus = "status_global"
)

// Metric descriptors.
var (
	globalCommandsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "status", "global"),
		"Total number of executed MySQL commands.",
		[]string{"name"}, nil,
	)
)

// ScrapeStatusGlobal collects from `SHOW GLOBAL STATUS`.
type ScrapeStatusGlobal struct{}

// Name of the Scraper. Should be unique.
func (ScrapeStatusGlobal) Name() string {
	return globalStatus
}

// Help describes the role of the Scraper.
func (ScrapeStatusGlobal) Help() string {
	return "Collect from SHOW GLOBAL STATUS"
}

// Version of MySQL from which scraper is available.
func (ScrapeStatusGlobal) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeStatusGlobal) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	globalStatusRows, err := db.QueryContext(ctx, globalStatusQuery)
	if err != nil {
		return err
	}
	defer globalStatusRows.Close()

	var key string
	var val sql.RawBytes

	for globalStatusRows.Next() {
		if err := globalStatusRows.Scan(&key, &val); err != nil {
			return err
		}
		if floatVal, ok := parseStatus(val); ok {
			key = validPrometheusName(key)
			ch <- prometheus.MustNewConstMetric(globalCommandsDesc, prometheus.CounterValue, floatVal, key)
		}
	}

	return nil
}
func validPrometheusName(s string) string {
	nameRe := regexp.MustCompile("([^a-zA-Z0-9_])")
	s = nameRe.ReplaceAllString(s, "_")
	s = strings.ToLower(s)
	return s
}

// check interface
var _ Scraper = ScrapeStatusGlobal{}
