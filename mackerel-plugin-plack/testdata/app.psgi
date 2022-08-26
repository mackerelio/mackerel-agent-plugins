#!/usr/bin/env perl

use strict;
use warnings;
use utf8;

use Plack::Builder;

builder {
    enable "Plack::Middleware::ServerStatus::Lite",
        path         => '/server-status',
        allow        => [ '0.0.0.0/0' ],
        counter_file => '/tmp/counter_file',
        scoreboard   => '/var/run/server';

    sub {
        [200, ["Content-Type" => "text/plain"], "ok"];
    }
};

