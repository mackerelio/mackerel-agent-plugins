use 5.014;
use warnings;
use utf8;
use File::Spec;
use Test::More;

# The plugins placed in other mackerelio's repositories
my $plugins_in_other_repository = [qw(
    mackerel-plugin-aws-cloudfront
    mackerel-plugin-aws-dynamodb
    mackerel-plugin-aws-ec2
    mackerel-plugin-aws-ec2-cpucredit
    mackerel-plugin-aws-elasticache
    mackerel-plugin-aws-elasticsearch
    mackerel-plugin-aws-elb
    mackerel-plugin-aws-kinesis-firehose
    mackerel-plugin-aws-kinesis-streams
    mackerel-plugin-aws-lambda
    mackerel-plugin-aws-rds
    mackerel-plugin-aws-rekognition
    mackerel-plugin-aws-s3-requests
    mackerel-plugin-aws-ses
    mackerel-plugin-aws-waf
    mackerel-plugin-flume
    mackerel-plugin-gcp-compute-engine
    mackerel-plugin-gearmand
    mackerel-plugin-graphite
    mackerel-plugin-json
    mackerel-plugin-mcrouter
    mackerel-plugin-mongodb
    mackerel-plugin-murmur
    mackerel-plugin-mysql
    mackerel-plugin-nvidia-smi
    mackerel-plugin-rack-stats
    mackerel-plugin-redash
    mackerel-plugin-xentop
)];
my $is_in_other_repository = {
    map { $_ => 1 } @$plugins_in_other_repository,
};

for my $dir (<mackerel-plugin-*>) {
    my $maingo = File::Spec->catfile($dir, 'main.go');
    ok -f -r $maingo or diag "$maingo not found";
    my $readmemd = File::Spec->catfile($dir, 'README.md');
    ok -f -r $readmemd or diag "$readmemd is not available.";

    my $package = $dir;
       $package =~ s/(mackerel-plugin)?-//g;
       $package = "mp$package";
    my $import = sprintf(
        "github.com/mackerelio/%s/lib",
        $is_in_other_repository->{$dir} ? $dir : "mackerel-agent-plugins/$dir",
    );
    my $expect = qq[package main

import "$import"

func main() {
\t$package.Do()
}
];
    my $got = do {
        local $/;
        open my $fh, '<:encoding(UTF-8)', $maingo or die $!;
        <$fh>
    };
    is $got, $expect, "The contents of $maingo does not follow the convention.";
}

done_testing;
