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

// Scrape `information_schema.processlist`.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const infoSchemaProcesslistCountQuery = `
		SELECT user,
			SUBSTRING_INDEX(host, ':', 1) AS host,
			db                            as dbname,
			COALESCE(command, '')         AS command,
			COALESCE(state, '')           AS state,
			COUNT(*)                      AS processes,
			SUM(time)                     AS time
		FROM information_schema.processlist
		WHERE ID != connection_id()
		GROUP BY user, SUBSTRING_INDEX(host, ':', 1), db, command, state
		`

// Metric descriptors.
var (
	processlistProcessesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "processes"),
		"The number of threads (connections) split by current state.",
		[]string{"user", "host", "db", "command", "state"}, nil)
	processlistProcessesTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "processes_seconds"),
		"The number of seconds threads (connections) have used split by current state.",
		[]string{"user", "host", "db", "command", "state"}, nil)
)

// ScrapeProcesslist collects from `information_schema.processlist`.
type ScrapeProcesslist struct{}

// Name of the Scraper. Should be unique.
func (ScrapeProcesslist) Name() string {
	return informationSchema + ".processlist"
}

// Help describes the role of the Scraper.
func (ScrapeProcesslist) Help() string {
	return "Collect current thread state counts from the information_schema.processlist"
}

// Version of MySQL from which scraper is available.
func (ScrapeProcesslist) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeProcesslist) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {

	processlistCountRows, err := db.QueryContext(ctx, infoSchemaProcesslistCountQuery)
	if err != nil {
		return err
	}
	defer processlistCountRows.Close()

	var (
		user       string
		host       string
		dbname      string
		command     string
		state        string
		processes       uint32
		processesTime      uint32
	)

	for processlistCountRows.Next() {
		err = processlistCountRows.Scan(&user, &host, &dbname, &command, &state, &processes, &processesTime)
		if err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(processlistProcessesDesc, prometheus.GaugeValue, float64(processes), user, host, dbname, command, state)
		ch <- prometheus.MustNewConstMetric(processlistProcessesTimeDesc, prometheus.GaugeValue, float64(processesTime), user, host, dbname, command, state)
	}

	return nil
}

// check interface
var _ Scraper = ScrapeProcesslist{}
