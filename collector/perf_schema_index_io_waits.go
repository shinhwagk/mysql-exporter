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

// Scrape `performance_schema.table_io_waits_summary_by_index_usage`.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const perfIndexIOWaitsQuery = `
	SELECT OBJECT_SCHEMA, OBJECT_NAME, INDEX_NAME,
		COUNT_STAR, COUNT_READ, COUNT_WRITE, COUNT_FETCH, COUNT_INSERT, COUNT_UPDATE, COUNT_DELETE,
	    SUM_TIMER_WAIT, SUM_TIMER_READ, SUM_TIMER_WRITE, SUM_TIMER_FETCH, SUM_TIMER_INSERT, SUM_TIMER_UPDATE, SUM_TIMER_DELETE
	  FROM performance_schema.table_io_waits_summary_by_index_usage
	  WHERE OBJECT_SCHEMA NOT IN ('mysql', 'performance_schema') AND index_name is not null
	`

// Metric descriptors.
var (
	psiiw1 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_star_total"),
		"The total number of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw2 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_read_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw3 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_write_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw4 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_fetch_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw5 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_insert_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw6 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_update_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw7 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_delete_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw8 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw9 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_read_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw10 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_write_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw11 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_fatch_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw12 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_insert_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw13 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_update_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
	psiiw14 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "index_io_waits_delete_millisecond_total"),
		"The total time of table I/O wait events for each table and operation.",
		[]string{"object_schema", "object_name"}, nil,
	)
)

// ScrapePerfIndexIOWaits collects for `performance_schema.table_io_waits_summary_by_index_usage`.
type ScrapePerfIndexIOWaits struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfIndexIOWaits) Name() string {
	return "perf_schema.indexiowaits"
}

// Help describes the role of the Scraper.
func (ScrapePerfIndexIOWaits) Help() string {
	return "Collect metrics from performance_schema.table_io_waits_summary_by_index_usage"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfIndexIOWaits) Version() float64 {
	return 5.6
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfIndexIOWaits) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	perfSchemaIndexWaitsRows, err := db.QueryContext(ctx, perfIndexIOWaitsQuery)
	if err != nil {
		return err
	}
	defer perfSchemaIndexWaitsRows.Close()

	var (
		objectSchema, objectName, indexName                                                 string
		countStar, countRead, countWrite, countFetch, countInsert, countUpdate, countDelete uint64
		timerWait, timeRead, timeWrite, timeFetch, timeInsert, timeUpdate, timeDelete       uint64
	)

	for perfSchemaIndexWaitsRows.Next() {
		if err := perfSchemaIndexWaitsRows.Scan(
			&objectSchema, &objectName, &indexName,
			&countStar, &countRead, &countWrite, &countFetch, &countInsert, &countUpdate, &countDelete,
			&timerWait, &timeRead, &timeWrite, &timeFetch, &timeInsert, &timeUpdate, &timeDelete,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(psiiw1, prometheus.CounterValue, float64(countStar), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw2, prometheus.CounterValue, float64(countRead), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw3, prometheus.CounterValue, float64(countWrite), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw4, prometheus.CounterValue, float64(countFetch), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw5, prometheus.CounterValue, float64(countInsert), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw6, prometheus.CounterValue, float64(countUpdate), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw7, prometheus.CounterValue, float64(countDelete), objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw8, prometheus.CounterValue, float64(timerWait)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw9, prometheus.CounterValue, float64(timeRead)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw10, prometheus.CounterValue, float64(timeWrite)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw11, prometheus.CounterValue, float64(timeFetch)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw12, prometheus.CounterValue, float64(timeInsert)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw13, prometheus.CounterValue, float64(timeUpdate)/picoMilliseconds, objectSchema, objectName)
		ch <- prometheus.MustNewConstMetric(psiiw14, prometheus.CounterValue, float64(timeDelete)/picoMilliseconds, objectSchema, objectName)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfIndexIOWaits{}
