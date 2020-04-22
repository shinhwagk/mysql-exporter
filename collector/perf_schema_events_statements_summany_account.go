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

const perfEventsStatementsSummanyAccountQuery = `
	SELECT
		user,
		host,
		event_name,
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
		sum_no_good_index_used
	FROM performance_schema.events_statements_summary_by_account_by_event_name
		WHERE user is not null
	`

// Metric descriptors.
var (
	performanceSchemaEventsStatementsSummaryAccountCountStarDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_count_star_total"),
		"performance_schema.events_statements_summary_by_digest.count_star",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountTimerWaitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_timer_wait_total"),
		"performance_schema.events_statements_summary_by_digest.sum_timer_wait",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountLockTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_lock_time_total"),
		"performance_schema.events_statements_summary_by_digest.sum_lock_time",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_errors_total"),
		"performance_schema.events_statements_summary_by_digest.sum_errors",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountWarningsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_warnings_total"),
		"performance_schema.events_statements_summary_by_digest.sum_warnings",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountRowsAffectedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_rows_affected_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_affected",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountRowsSentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_rows_sent_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_sent",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountRowsExaminedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_rows_examined_total"),
		"performance_schema.events_statements_summary_by_digest.sum_rows_examined",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountCreatedTmpDiskTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_created_tmp_disk_tables_total"),
		"performance_schema.events_statements_summary_by_digest.sum_created_tmp_disk_tables",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountCreatedTmpTablesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_created_tmp_tables_total"),
		"performance_schema.events_statements_summary_by_digest.sum_created_tmp_tables",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSelectFullJoinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_select_full_join"),
		"performance_schema.events_statements_summary_by_digest.sum_select_full_join",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSelectFullRangeJoinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_select_full_range_join"),
		"performance_schema.events_statements_summary_by_digest.sum_select_full_range_join",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSelectRangeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_select_range"),
		"performance_schema.events_statements_summary_by_digest.sum_select_range",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSelectRangeCheckDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_select_range_check"),
		"performance_schema.events_statements_summary_by_digest.sum_select_range_check",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSelectScanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_select_scan"),
		"performance_schema.events_statements_summary_by_digest.sum_select_scan",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSortMergePassesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_sort_merge_passes_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_merge_passes",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSortRangeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_sort_range_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_range",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSortRowsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_sort_rows_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_rows",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountSortScanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_sort_scan_total"),
		"performance_schema.events_statements_summary_by_digest.sum_sort_scan",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountNoIndexUsedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_no_index_used_total"),
		"performance_schema.events_statements_summary_by_digest.sum_no_index_used",
		[]string{"user", "host", "event"}, nil,
	)
	performanceSchemaEventsStatementsSummaryAccountNoGoodIndexUsedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "events_statements_summary_account_no_good_index_used_total"),
		"performance_schema.events_statements_summary_by_digest.sum_no_good_index_used",
		[]string{"user", "host", "event"}, nil,
	)
)

// ScrapePerfEventsStatementsSummaryAccount collects from `performance_schema.events_statements_summary_by_digest`.
type ScrapePerfEventsStatementsSummaryAccount struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfEventsStatementsSummaryAccount) Name() string {
	return "perf_schema.eventsstatementssummaryaccount"
}

// Help describes the role of the Scraper.
func (ScrapePerfEventsStatementsSummaryAccount) Help() string {
	return "Collect metrics from performance_schema.events_statements_summary_by_digest"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfEventsStatementsSummaryAccount) Version() float64 {
	return 5.6
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfEventsStatementsSummaryAccount) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	// Timers here are returned in picoseconds.
	perfSchemaEventsStatementsRows, err := db.QueryContext(ctx, perfEventsStatementsSummanyAccountQuery)
	if err != nil {
		return err
	}
	defer perfSchemaEventsStatementsRows.Close()

	var (
		user, host, event                                                              string
		countStar, timerWait, lockTime                                                 uint64
		errors, warnings                                                               uint64
		rowsAffected, rowsSent, rowsExamined                                           uint64
		tmpDiskTables, tmpTables                                                       uint64
		selectFullJoin, selectfullRangeJoin, selectRange, selectRangeCheck, selectScan uint64
		sortMergePasses, sortRange, sortRows, sortScan                                 uint64
		noIndexUsed, noGoodIndexUsed                                                   uint64
	)
	for perfSchemaEventsStatementsRows.Next() {
		if err := perfSchemaEventsStatementsRows.Scan(
			&user, &host, &event,
			&countStar, &timerWait, &lockTime,
			&errors, &warnings,
			&rowsAffected, &rowsSent, &rowsExamined,
			&tmpDiskTables, &tmpTables,
			&selectFullJoin, &selectfullRangeJoin, &selectRange, &selectRangeCheck, &selectScan,
			&sortMergePasses, &sortRange, &sortRows, &sortScan,
			&noIndexUsed, &noGoodIndexUsed,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountCountStarDesc, prometheus.CounterValue, float64(countStar),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountTimerWaitDesc, prometheus.CounterValue, float64(timerWait)/picoMilliseconds,
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountLockTimeDesc, prometheus.CounterValue, float64(lockTime)/picoMilliseconds,
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountErrorsDesc, prometheus.CounterValue, float64(errors),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountWarningsDesc, prometheus.CounterValue, float64(warnings),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountRowsAffectedDesc, prometheus.CounterValue, float64(rowsAffected),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountRowsSentDesc, prometheus.CounterValue, float64(rowsSent),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountRowsExaminedDesc, prometheus.CounterValue, float64(rowsExamined),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountCreatedTmpDiskTablesDesc, prometheus.CounterValue, float64(tmpDiskTables),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountCreatedTmpTablesDesc, prometheus.CounterValue, float64(tmpTables),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSelectFullJoinDesc, prometheus.CounterValue, float64(selectFullJoin),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSelectFullRangeJoinDesc, prometheus.CounterValue, float64(selectfullRangeJoin),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSelectRangeDesc, prometheus.CounterValue, float64(selectRange),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSelectRangeCheckDesc, prometheus.CounterValue, float64(selectRangeCheck),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSelectScanDesc, prometheus.CounterValue, float64(selectScan),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSortMergePassesDesc, prometheus.CounterValue, float64(sortMergePasses),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSortRangeDesc, prometheus.CounterValue, float64(sortRange),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSortRowsDesc, prometheus.CounterValue, float64(sortRows),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountSortScanDesc, prometheus.CounterValue, float64(sortScan),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountNoIndexUsedDesc, prometheus.CounterValue, float64(noIndexUsed),
			user, host, event,
		)
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaEventsStatementsSummaryAccountNoGoodIndexUsedDesc, prometheus.CounterValue, float64(noGoodIndexUsed),
			user, host, event,
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfEventsStatementsSummaryAccount{}
