-- current lsn range
with clr as (
select start_lsn, start_lsn + total_length end_lsn from (
    select t.*, (end_lsn - start_lsn ) * 2 total_length, (select VARIABLE_VALUE 
            from performance_schema.global_status
            where VARIABLE_NAME='Innodb_redo_log_checkpoint_lsn') checkpoint_lsn   
    FROM performance_schema.innodb_redo_log_files t
) x where x.checkpoint_lsn between x.start_lsn and x.end_lsn)
select 'current_lsn' lsn_point, round((variable_value - clr.start_lsn) / (clr.end_lsn - clr.start_lsn),2) from performance_schema.global_status cl, clr 
where VARIABLE_NAME='Innodb_redo_log_current_lsn'
union all
select 'flushed_lsn' lsn_point, round((variable_value - clr.start_lsn) / (clr.end_lsn - clr.start_lsn),2) from performance_schema.global_status cl, clr 
where VARIABLE_NAME='Innodb_redo_log_flushed_to_disk_lsn';

-- 查看当前整个redo的lsn range
select start_lsn, start_lsn + total_length end_lsn from (
    select t.*, (end_lsn - start_lsn ) * 2 total_length, (select VARIABLE_VALUE 
            from performance_schema.global_status
            where VARIABLE_NAME='Innodb_redo_log_checkpoint_lsn') checkpoint_lsn   
    FROM performance_schema.innodb_redo_log_files t
) x where x.checkpoint_lsn between x.start_lsn and x.end_lsn
