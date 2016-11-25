<?php

header("Content-Type: text/plain");

$status = opcache_get_status();
$config = opcache_get_configuration();

$stats = array(
    // memory_usage
    'used_memory'   =>   $status['memory_usage']['used_memory'],
    'free_memory'   => $status['memory_usage']['free_memory'],
    'wasted_memory' => $status['memory_usage']['wasted_memory'],
    'current_wasted_percentage' => $status['memory_usage']['current_wasted_percentage'],

    // opcache_statistics
    'num_cached_scripts'   => $status['opcache_statistics']['num_cached_scripts'],
    'num_cached_keys'      => $status['opcache_statistics']['num_cached_keys'],
    'max_cached_keys'      => $status['opcache_statistics']['max_cached_keys'],
    'hits'                 => $status['opcache_statistics']['hits'],
    'oom_restarts'         => $status['opcache_statistics']['oom_restarts'],
    'hash_restarts'        => $status['opcache_statistics']['hash_restarts'],
    'manual_restarts'      => $status['opcache_statistics']['manual_restarts'],
    'misses'               => $status['opcache_statistics']['misses'],
    'blacklist_misses'     => $status['opcache_statistics']['blacklist_misses'],
    'blacklist_miss_ratio' => $status['opcache_statistics']['blacklist_miss_ratio'],
    'opcache_hit_rate'     => $status['opcache_statistics']['opcache_hit_rate'],
);

foreach( $stats as $name => $value ){
    echo sprintf( "%s:%d\n", $name,  $value );
}
