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

// Scrape `performance_schema.file_summary_by_instance`.

package collector

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const perfFileInstancesQuery = `
	SELECT
		file_name, event_name,
		count_star, sum_timer_wait,
		count_read, sum_timer_read, sum_number_of_bytes_read,
		count_write, sum_timer_write, sum_number_of_bytes_write,
		count_misc, sum_timer_misc
	  FROM performance_schema.file_summary_by_instance
	`

// // Tunable flags.
// var (
// 	performanceSchemaFileInstancesFilter = kingpin.Flag(
// 		"collect.perf_schema.file_instances.filter",
// 		"RegEx file_name filter for performance_schema.file_summary_by_instance",
// 	).Default(".*").String()
// )

// Metric descriptors.
var (
	performanceSchemaFileInstancesRemovePrefix = kingpin.Flag(
		"collect.perf_schema.file_instances.remove_prefix",
		"Remove path prefix in performance_schema.file_summary_by_instance",
	).Default("/var/lib/mysql/").String()

	performanceSchemaSummaryFileInstancesCountStarDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_count_star"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileInstancesTimerWaitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_timer_wait"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileReadCountInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_read_count"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileReadTimerInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_read_timer"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileReadBytesInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_read_bytes"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileWriteCountInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_write_count"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileWriteTimerInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_write_timer"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileWriteBytesInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_write_bytes"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileMiscCountInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_misc_count"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
	performanceSchemaSummaryFileMiscTimerInstancesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "file_summary_instances_misc_timer"),
		"The number of bytes processed by file read/write operations.",
		[]string{"file_name", "event_name"}, nil,
	)
)

// ScrapePerfFileInstances collects from `performance_schema.file_summary_by_instance`.
type ScrapePerfFileInstances struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfFileInstances) Name() string {
	return "perf_schema.file_summary_instances"
}

// Help describes the role of the Scraper.
func (ScrapePerfFileInstances) Help() string {
	return "Collect metrics from performance_schema.file_summary_by_instance"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfFileInstances) Version() float64 {
	return 5.5
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfFileInstances) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	// Timers here are returned in picoseconds.
	perfSchemaFileInstancesRows, err := db.QueryContext(ctx, perfFileInstancesQuery)
	if err != nil {
		return err
	}
	defer perfSchemaFileInstancesRows.Close()

	var (
		fileName, eventName                        string
		countStar, timerWait                       uint64
		countRead, timerRead, memberOfBytesRead    uint64
		countWrite, timerWrite, memberOfBytesWrite uint64
		countMisc, timerMisc                       uint64
	)

	for perfSchemaFileInstancesRows.Next() {
		if err := perfSchemaFileInstancesRows.Scan(
			&fileName, &eventName,
			&countStar, &timerWait,
			&countRead, &timerRead, &memberOfBytesRead,
			&countWrite, &timerWrite, &memberOfBytesWrite,
			&countMisc, &timerMisc,
		); err != nil {
			return err
		}

		fileName = strings.TrimPrefix(fileName, *performanceSchemaFileInstancesRemovePrefix)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileInstancesCountStarDesc, prometheus.CounterValue, float64(countRead), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileInstancesTimerWaitDesc, prometheus.CounterValue, float64(timerWait), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileReadCountInstancesDesc, prometheus.CounterValue, float64(countRead), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileReadTimerInstancesDesc, prometheus.CounterValue, float64(timerRead), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileReadBytesInstancesDesc, prometheus.CounterValue, float64(memberOfBytesRead), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileWriteCountInstancesDesc, prometheus.CounterValue, float64(countWrite), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileWriteTimerInstancesDesc, prometheus.CounterValue, float64(timerWrite), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileWriteBytesInstancesDesc, prometheus.CounterValue, float64(memberOfBytesWrite), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileMiscCountInstancesDesc, prometheus.CounterValue, float64(countMisc), fileName, eventName)
		ch <- prometheus.MustNewConstMetric(performanceSchemaSummaryFileMiscTimerInstancesDesc, prometheus.CounterValue, float64(timerMisc), fileName, eventName)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfFileInstances{}
