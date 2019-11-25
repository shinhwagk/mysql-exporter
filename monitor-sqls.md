## performance_schema

| Name                                                 | Categories      | Description                                            | Used? |
| ---------------------------------------------------- | --------------- | ------------------------------------------------------ | ----- |
| events_statements_summary_by_account_by_event_name   | event,user,host |                                                        | YES   |
| events_statements_summary_by_thread_by_event_name    | event,thread    | join threads                                           | YES   |
| events_statements_summary_by_user_by_event_name      |                 | use events_statements_summary_by_account_by_event_name | NO    |
| events_statements_summary_global_by_event_name       | event           | use events_statements_summary_by_account_by_event_name | NO    |
| events_statements_summary_by_host_by_event_name      |                 | use events_statements_summary_by_account_by_event_name | NO    |
| events_statements_summary_by_digest                  | sql,schema      |                                                        | YES   |
| events_statements_summary_by_program                 |                 |                                                        | ?     |
| events_statements_current                            |                 |                                                        | NO    |
| events_statements_histogram_by_digest                |                 |                                                        | YES   |
| events_statements_histogram_global                   |                 |                                                        |       |
| events_statements_history                            | sql,tread,event |                                                        | NO    |
| events_statements_history_long                       | sql,tread,event |                                                        | NO    |
|                                                      |                 |                                                        |       |
| events_stages_current                                |                 |                                                        | NO    |
| events_stages_history                                |                 |                                                        | NO    |
| events_stages_history_long                           |                 |                                                        | NO    |
| events_stages_summary_by_account_by_event_name       |                 |                                                        | NO    |
| events_stages_summary_by_host_by_event_name          |                 |                                                        | NO    |
| events_stages_summary_by_thread_by_event_name        |                 |                                                        | NO    |
| events_stages_summary_by_user_by_event_name          |                 |                                                        | NO    |
| events_stages_summary_global_by_event_name           |                 |                                                        | NO    |
|                                                      |                 |                                                        |       |
| events_waits_current                                 |                 |                                                        |       |
| events_waits_history                                 |                 |                                                        |       |
| events_waits_history_long                            |                 |                                                        |       |
| events_waits_summary_by_account_by_event_name        |                 |                                                        |       |
| events_waits_summary_by_host_by_event_name           |                 |                                                        |       |
| events_waits_summary_by_instance                     |                 |                                                        |       |
| events_waits_summary_by_thread_by_event_name         |                 |                                                        |       |
| events_waits_summary_by_user_by_event_name           |                 |                                                        |       |
| events_waits_summary_global_by_event_name            |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| memory_summary_by_account_by_event_name              |                 |                                                        |       |
| memory_summary_by_host_by_event_name                 |                 |                                                        |       |
| memory_summary_by_thread_by_event_name               |                 |                                                        |       |
| memory_summary_by_user_by_event_name                 |                 |                                                        |       |
| memory_summary_global_by_event_name                  |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| events_transactions_current                          |                 |                                                        |       |
| events_transactions_history                          |                 |                                                        |       |
| events_transactions_history_long                     |                 |                                                        |       |
| events_transactions_summary_by_account_by_event_name |                 |                                                        |       |
| events_transactions_summary_by_host_by_event_name    |                 |                                                        |       |
| events_transactions_summary_by_thread_by_event_name  |                 |                                                        |       |
| events_transactions_summary_by_user_by_event_name    |                 |                                                        |       |
| events_transactions_summary_global_by_event_name     |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| table_handles                                        |                 |                                                        |       |
| table_io_waits_summary_by_index_usage                |                 |                                                        |       |
| table_io_waits_summary_by_table                      |                 |                                                        |       |
| table_lock_waits_summary_by_table                    |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| global_status                                        |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| file_instances                                       |                 |                                                        |       |
| file_summary_by_event_name                           |                 |                                                        |       |
| file_summary_by_instance                             |                 |                                                        |       |
|                                                      |                 |                                                        |       |
| threads                                              |                 |                                                        |       |

```
-- accounts
-- cond_instances

-- data_lock_waits
-- data_locks


-- events_errors_summary_by_account_by_error
-- events_errors_summary_by_host_by_error
-- events_errors_summary_by_thread_by_error
-- events_errors_summary_by_user_by_error
-- events_errors_summary_global_by_error

 

-- global_variables
-- host_cache
-- -- hosts
-- keyring_keys
-- log_status


-- -- metadata_locks
-- -- mutex_instances
-- -- objects_summary_global_by_type
-- performance_timers
-- persisted_variables
-- prepared_statements_instances
-- replication_applier_configuration
-- replication_applier_filters
-- replication_applier_global_filters
-- replication_applier_status
-- replication_applier_status_by_coordinator
-- replication_applier_status_by_worker
-- replication_connection_configuration
-- replication_connection_status
-- replication_group_member_stats
-- replication_group_members
-- - rwlock_instances
-- - session_account_connect_attrs
-- - session_connect_attrs
-- -- session_status
-- - session_variables
-- setup_actors
-- setup_consumers
-- setup_instruments
-- setup_objects
-- setup_threads
-- socket_instances
-- socket_summary_by_event_name
-- socket_summary_by_instance
-- -- status_by_account
-- -- status_by_host
-- -- status_by_thread
-- -- status_by_user

-- -- 
-- -- user_defined_functions
-- user_variables_by_thread
-- -- users
-- variables_by_thread
-- variables_info
```