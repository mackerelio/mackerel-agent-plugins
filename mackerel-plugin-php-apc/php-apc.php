<?php

header("Content-Type: text/plain");

$cache      = apc_cache_info();
$cache_user = apc_cache_info('user', 1); 
$mem        = apc_sma_info();

$stats = array(
    "memory_segments"       => (int)$mem['num_seg'],
    "segment_size"          => (int)$mem['seg_size'],
    "total_memory"          => (int)$mem['num_seg'] * $mem['seg_size'],
    "cached_files_count"    => (int)$cache['num_entries'],
    "cached_files_size"     => (int)$cache['mem_size'],
    "cache_hits"            => (int)$cache['num_hits'],
    "cache_misses"          => (int)$cache['num_misses'],
    "cache_full_count"      => (int)$cache['expunges'],
    "user_cache_vars_count" => (int)$cache_user['num_entries'],
    "user_cache_vars_size"  => (int)$cache_user['mem_size'],
    "user_cache_hits"       => (int)$cache_user['num_hits'],
    "user_cache_misses"     => (int)$cache_user['num_misses'],
    "user_cache_full_count" => (int)$cache_user['expunges'],
);

foreach( $stats as $name => $value ){
    echo sprintf( "%s:%d\n", $name,  $value );
}
