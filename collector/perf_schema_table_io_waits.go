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

// Scrape `performance_schema.table_io_waits_summary_by_table`.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const perfTableIOWaitsQuery = `
	SELECT
	    OBJECT_TYPE, OBJECT_SCHEMA, OBJECT_NAME,
	    COUNT_STAR, COUNT_READ, COUNT_WRITE, COUNT_FETCH, COUNT_INSERT, COUNT_UPDATE, COUNT_DELETE,
	    SUM_TIMER_WAIT, SUM_TIMER_READ, SUM_TIMER_WRITE, SUM_TIMER_FETCH, SUM_TIMER_INSERT, SUM_TIMER_UPDATE, SUM_TIMER_DELETE
	  FROM performance_schema.table_io_waits_summary_by_table
	  WHERE OBJECT_SCHEMA NOT IN ('mysql', 'performance_schema')
	`

// Metric descriptors.
var (
	pstiw1 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_star_total"),
		"The total number of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw2 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_read_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw3 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_write_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw4 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_fetch_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw5 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_insert_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw6 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_update_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw7 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_delete_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw8 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw9 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_read_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw10 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_write_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw11 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_fatch_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw12 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_insert_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw13 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_update_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
	pstiw14 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "table_io_waits_delete_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_type", "object_schema", "object_name"}, nil,
	)
)

// ScrapePerfTableIOWaits collects from `performance_schema.table_io_waits_summary_by_table`.
type ScrapePerfTableIOWaits struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfTableIOWaits) Name() string {
	return "perf_schema.tableiowaits"
}

// Help describes the role of the Scraper.
func (ScrapePerfTableIOWaits) Help() string {
	return "Collect metrics from performance_schema.table_io_waits_summary_by_table"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfTableIOWaits) Version() float64 {
	return 5.6
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfTableIOWaits) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	perfSchemaTableWaitsRows, err := db.QueryContext(ctx, perfTableIOWaitsQuery)
	if err != nil {
		return err
	}
	defer perfSchemaTableWaitsRows.Close()

	var (
		objectType, objectSchema, objectName                                                string
		countStar, countRead, countWrite, countFetch, countInsert, countUpdate, countDelete uint64
		timerWait, timeRead, timeWrite, timeFetch, timeInsert, timeUpdate, timeDelete       uint64
	)

	for perfSchemaTableWaitsRows.Next() {
		if err := perfSchemaTableWaitsRows.Scan(
			&objectType, &objectSchema, &objectName,
			&countStar, &countRead, &countWrite, &countFetch, &countInsert, &countUpdate, &countDelete,
			&timerWait, &timeRead, &timeWrite, &timeFetch, &timeInsert, &timeUpdate, &timeDelete,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(pstiw1, prometheus.CounterValue, float64(countStar), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw2, prometheus.CounterValue, float64(countRead), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw3, prometheus.CounterValue, float64(countWrite), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw4, prometheus.CounterValue, float64(countFetch), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw5, prometheus.CounterValue, float64(countInsert), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw6, prometheus.CounterValue, float64(countUpdate), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw7, prometheus.CounterValue, float64(countDelete), objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw8, prometheus.CounterValue, float64(timerWait)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw9, prometheus.CounterValue, float64(timeRead)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw10, prometheus.CounterValue, float64(timeWrite)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw11, prometheus.CounterValue, float64(timeFetch)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw12, prometheus.CounterValue, float64(timeInsert)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw13, prometheus.CounterValue, float64(timeUpdate)/picoMilliseconds, objectType, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(pstiw14, prometheus.CounterValue, float64(timeDelete)/picoMilliseconds, objectType, objectSchema, objectName)

	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfTableIOWaits{}
