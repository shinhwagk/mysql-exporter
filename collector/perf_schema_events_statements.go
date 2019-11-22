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
		AND last_seen >= DATE_SUB(NOW(), INTERVAL 60 SECOND)
	`

// Metric descriptors.
var (
	performanceSchemaEventsStatementsCountStarDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_count_star_total"),
		"performance_schema.events_statements_summary_by_digest.count_star",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsTimerWaitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_timer_wait_total"),
		"performance_schema.events_statements_summary_by_digest.sum_timer_wait",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsLockTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_lock_time_total"),
		"performance_schema.events_statements_summary_by_digest.sum_lock_time",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_errors_total"),
		"performance_schema.events_statements_summary_by_digest.sum_errors",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsWarningsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_warnings_total"),
		"performance_schema.events_statements_summary_by_digest.sum_warnings",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsRowsAffectedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_affected_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_affected",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsRowsSentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_sent_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_sent",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsRowsExaminedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_rows_examined_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_examined",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsCreatedTmpDiskTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_created_tmp_disk_tables_total"),
		"performance_schema.events_statements_summary_by_digest.sum_created_tmp_disk_tables",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsCreatedTmpTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_created_tmp_tables_total"),
		"performance_schema.events_statements_summary_by_digest.sum_created_tmp_tables",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSelectFullJoinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_select_full_join"),
		"performance_schema.events_statements_summary_by_digest.sum_select_full_join",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSelectFullRangeJoinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_select_full_range_join"),
		"performance_schema.events_statements_summary_by_digest.sum_select_full_range_join",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSelectRangeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_select_range"),
		"performance_schema.events_statements_summary_by_digest.sum_select_range",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSelectRangeCheckDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_select_range_check"),
		"performance_schema.events_statements_summary_by_digest.sum_select_range_check",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSelectScanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_select_scan"),
		"performance_schema.events_statements_summary_by_digest.sum_select_scan",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSortMergePassesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_merge_passes_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_merge_passes",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSortRangeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_range_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_range",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSortRowsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_rows_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_rows",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsSortScanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_sort_scan_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_scan",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsNoIndexUsedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_no_index_used_total"),
		"performance_schema.events_statements_summary_by_digest.sum_no_index_used",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsNoGoodIndexUsedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_no_good_index_used_total"),
		"performance_schema.events_statements_summary_by_digest.sum_no_good_index_used",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsQuantile95Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_quantile_95_total"),
		"performance_schema.events_statements_summary_by_digest.quantile_95",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsQuantile99Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_quantile_99_total"),
		"performance_schema.events_statements_summary_by_digest.quantile_99",
		[]string{"schema", "digest"}, nil,
	)
	performanceSchemaEventsStatementsQuantile999Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_quantile_999_total"),
		"performance_schema.events_statements_summary_by_digest.quantile_999",
		[]string{"schema", "digest"}, nil,
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
		schemaName, digest                                                             string
		countStar, timerWait, lockTime                                                 uint64
		errors, warnings                                                               uint64
		rowsAffected, rowsSent, rowsExamined                                           uint64
		tmpTables, tmpDiskTables                                                       uint64
		selectFullJoin, selectfullRangeJoin, selectRange, selectRangeCheck, selectScan uint64
		sortMergePasses, sortRange, sortRows, sortScan                                 uint64
		noIndexUsed, noGoodIndexUsed                                                   uint64
		quantile95, quantile99, quantile999                                            uint64
	)
	for perfSchemaEventsStatementsRows.Next() {
		if err := perfSchemaEventsStatementsRows.Scan(
			&schemaName, &digest,
			&countStar, &timerWait, &lockTime,
			&errors, &warnings, &rowsAffected,
			&rowsExamined, &rowsSent, &rowsExamined,
			&tmpTables, &tmpDiskTables,
			&selectFullJoin, &selectfullRangeJoin, &selectRange, &selectRangeCheck, &selectScan,
			&sortMergePasses, &sortRange, &sortScan,
			&noIndexUsed, &noGoodIndexUsed,
			&quantile95, &quantile99, &quantile999,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsCountStarDesc, prometheus.CounterValue, float64(countStar),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsTimerWaitDesc, prometheus.CounterValue, float64(timerWait)/picoSeconds,
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsLockTimeDesc, prometheus.CounterValue, float64(lockTime)/picoSeconds,
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
			performanceSchemaEventsStatementsCreatedTmpDiskTablesDesc, prometheus.CounterValue, float64(tmpDiskTables),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsCreatedTmpTablesDesc, prometheus.CounterValue, float64(tmpTables),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSelectFullJoinDesc, prometheus.CounterValue, float64(selectFullJoin),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSelectFullRangeJoinDesc, prometheus.CounterValue, float64(selectfullRangeJoin),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSelectRangeDesc, prometheus.CounterValue, float64(selectRange),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSelectRangeCheckDesc, prometheus.CounterValue, float64(selectRangeCheck),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSelectScanDesc, prometheus.CounterValue, float64(selectScan),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortMergePassesDesc, prometheus.CounterValue, float64(sortMergePasses),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortRangeDesc, prometheus.CounterValue, float64(sortRange),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortRowsDesc, prometheus.CounterValue, float64(sortRows),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSortScanDesc, prometheus.CounterValue, float64(sortScan),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsNoIndexUsedDesc, prometheus.CounterValue, float64(noIndexUsed),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsNoGoodIndexUsedDesc, prometheus.CounterValue, float64(noGoodIndexUsed),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsQuantile95Desc, prometheus.CounterValue, float64(quantile95),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsQuantile99Desc, prometheus.CounterValue, float64(quantile99),
			schemaName, digest,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsQuantile999Desc, prometheus.CounterValue, float64(quantile999),
			schemaName, digest,
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfEventsStatements{}
