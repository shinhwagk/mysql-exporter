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
	globalStatus = "global_status"
)

// Regexp to match various groups of status vars.
var globalStatusRE = regexp.MustCompile(`^(com|handler|innodb|performance_schema|mysqlx)_(.*)$`)

// Metric descriptors.
var (
	globalCommandsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "com"),
		"Total number of executed MySQL commands.",
		[]string{"name"}, nil,
	)
	globalHandlerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "handlers"),
		"Total number of executed MySQL handlers.",
		[]string{"name"}, nil,
	)
	globalInnodbDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "innodb"),
		"Innodb.",
		[]string{"name"}, nil,
	)
	globalMysqlxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "mysqlx"),
		"mysqlx.",
		[]string{"name"}, nil,
	)
	globalPerformanceSchemaDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "performance_schema"),
		"Total number of MySQL instrumentations that could not be loaded or created due to memory constraints.",
		[]string{"name"}, nil,
	)
	globalOtherDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, globalStatus, "other"),
		"Total number of MySQL instrumentations that could not be loaded or created due to memory constraints.",
		[]string{"name"}, nil,
	)
)

// ScrapeGlobalStatus collects from `SHOW GLOBAL STATUS`.
type ScrapeGlobalStatus struct{}

// Name of the Scraper. Should be unique.
func (ScrapeGlobalStatus) Name() string {
	return globalStatus
}

// Help describes the role of the Scraper.
func (ScrapeGlobalStatus) Help() string {
	return "Collect from SHOW GLOBAL STATUS"
}

// Version of MySQL from which scraper is available.
func (ScrapeGlobalStatus) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeGlobalStatus) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
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
		if floatVal, ok := parseStatus(val); ok { // Unparsable values are silently skipped.
			key = validPrometheusName(key)
			match := globalStatusRE.FindStringSubmatch(key)
			if match == nil {
				ch <- prometheus.MustNewConstMetric(globalOtherDesc, prometheus.CounterValue, floatVal, key)
				continue
			}
			switch match[1] {
			case "com":
				ch <- prometheus.MustNewConstMetric(globalCommandsDesc, prometheus.CounterValue, floatVal, match[2])
			case "handler":
				ch <- prometheus.MustNewConstMetric(globalHandlerDesc, prometheus.CounterValue, floatVal, match[2])
			case "innodb":
				ch <- prometheus.MustNewConstMetric(globalInnodbDesc, prometheus.CounterValue, floatVal, match[2])
			case "performance_schema":
				ch <- prometheus.MustNewConstMetric(globalPerformanceSchemaDesc, prometheus.CounterValue, floatVal, match[2])
			case "mysqlx":
				ch <- prometheus.MustNewConstMetric(globalMysqlxDesc, prometheus.CounterValue, floatVal, match[2])
			}
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
var _ Scraper = ScrapeGlobalStatus{}
