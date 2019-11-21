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

// Scrape `performance_schema.events_statements_summary_by_digest`.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const perfEventsStatementsQuery = `
	SELECT
		schema_name,
		digest,
		count_star,
		sum_timer_wait,
		sum_lock_time,
		sum_errors,
		sum_warnings,
		sum_rows_affected,
		sum_rows_sent,
		sum_rows_examined,
		sum_created_tmp_disk_tables,
		sum_created_tmp_tables,
		sum_select_full_join,
		sum_select_full_range_join,
		sum_select_range,
		sum_select_range_check,
		sum_select_scan,
		sum_sort_merge_passes,
		sum_sort_range,
		sum_sort_rows,
		sum_sort_scan,
		sum_no_index_used,
		sum_no_good_index_used,
		quantile_95,
		quantile_99,
		quantile_999
	FROM performance_schema.events_statements_summary_by_digest
		WHERE schema_name NOT IN ('mysql', 'performance_schema', 'information_schema')
		AND last_seen >= DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d %H:%i'), INTERVAL 60 SECOND)
	`

// Metric descriptors.
var (
	performanceSchemaEventsStatementsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_total"),
		"The total count of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_seconds_total"),
		"The total time of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_errors_total"),
		"The errors of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsWarningsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_warnings_total"),
		"The warnings of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsRowsAffectedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_affected_total"),
		"The total rows affected of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsRowsSentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_sent_total"),
		"The total rows sent of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsRowsExaminedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_examined_total"),
		"The total rows examined of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsTmpTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_tmp_tables_total"),
		"The total tmp tables of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsTmpDiskTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_tmp_disk_tables_total"),
		"The total tmp disk tables of events statements by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsSortMergePassesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_merge_passes_total"),
		"The total number of merge passes by the sort algorithm performed by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsSortRowsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_rows_total"),
		"The total number of sorted rows by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
	performanceSchemaEventsStatementsNoIndexUsedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_no_index_used_total"),
		"The total number of statements that used full table scans by digest.",
		[]string{"schema", "digest", "digest_text"}, nil,
	)
)

// ScrapePerfEventsStatements collects from `performance_schema.events_statements_summary_by_digest`.
type ScrapePerfEventsStatements struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfEventsStatements) Name() string {
	return "perf_schema.eventsstatements"
}

// Help describes the role of the Scraper.
func (ScrapePerfEventsStatements) Help() string {
	return "Collect metrics from performance_schema.events_statements_summary_by_digest"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfEventsStatements) Version() float64 {
	return 5.6
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfEventsStatements) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	// Timers here are returned in picoseconds.
	perfSchemaEventsStatementsRows, err := db.QueryContext(ctx, perfEventsStatementsQuery)
	if err != nil {
		return err
	}
	defer perfSchemaEventsStatementsRows.Close()

	var (
		schemaName, digest                   string
		count, queryTime, errors, warnings   uint64
		rowsAffected, rowsSent, rowsExamined uint64
		tmpTables, tmpDiskTables             uint64
		sortMergePasses, sortRows            uint64
		noIndexUsed    uint64
	)
	for perfSchemaEventsStatementsRows.Next() {
		if err := perfSchemaEventsStatementsRows.Scan(
			&schemaName, &digest, &count, &queryTime, &errors, &warnings, &rowsAffected, &rowsSent, &rowsExamined, &tmpTables, &tmpDiskTables, &sortMergePasses, &sortRows, &noIndexUsed,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsDesc, prometheus.CounterValue, float64(count),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsTimeDesc, prometheus.CounterValue, float64(queryTime)/picoSeconds,
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsErrorsDesc, prometheus.CounterValue, float64(errors),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsWarningsDesc, prometheus.CounterValue, float64(warnings),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsRowsAffectedDesc, prometheus.CounterValue, float64(rowsAffected),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsRowsSentDesc, prometheus.CounterValue, float64(rowsSent),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsRowsExaminedDesc, prometheus.CounterValue, float64(rowsExamined),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsTmpTablesDesc, prometheus.CounterValue, float64(tmpTables),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsTmpDiskTablesDesc, prometheus.CounterValue, float64(tmpDiskTables),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortMergePassesDesc, prometheus.CounterValue, float64(sortMergePasses),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortRowsDesc, prometheus.CounterValue, float64(sortRows),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsNoIndexUsedDesc, prometheus.CounterValue, float64(noIndexUsed),
			schemaName, digest,
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfEventsStatements{}
