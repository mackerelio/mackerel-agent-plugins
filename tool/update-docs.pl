#!/usr/bin/env perl

use 5.014;
use strict;
use warnings;
use utf8;

use File::Copy qw/move/;
use JSON::PP qw/decode_json/;
use Path::Tiny qw/path/;

my $PLUGIN_PREFIX = 'mackerel-plugin-';
my $PACKAGE_NAME = 'mackerel-agent-plugins';

# refer Mackerel::ReleaseUtils
sub replace {
    my ($glob, $code) = @_;
    for my $file (glob $glob) {
        my $content = $code->(path($file)->slurp_utf8, $file);
        $content .= "\n" if $content !~ /\n\z/ms;

        my $f = path($file);
        # for keeping permission
        $f->append_utf8({truncate => 1}, $content);
    }
}

sub retrieve_plugins {
    # exclude plugins which has been moved to other repositories
    sort map {s/^$PLUGIN_PREFIX//; $_} grep { -e "$_/lib" } <$PLUGIN_PREFIX*>;
}

sub update_readme {
    my @plugins = @_;

    my $doc_links = '';
    for my $plug (@plugins) {
        $doc_links .= "* [$PLUGIN_PREFIX$plug](./$PLUGIN_PREFIX$plug/README.md)\n"
    }
    replace 'README.md' => sub {
        my $readme = shift;
        my $plu_reg = qr/$PLUGIN_PREFIX[-0-9a-zA-Z_]+/;
        $readme =~ s!(?:\* \[$plu_reg\]\(\./$plu_reg/README\.md\)\n)+!$doc_links!ms;
        $readme;
    };
}

sub update_packaging_specs {
    my @plugins = @_;
    my $for_in = 'for i in ' . join(' ', @plugins) . '; do';

    my $replace_sub = sub {
        my $content = shift;
        $content =~ s/for i in.*?;\s*do/$for_in/ms;
        $content;
    };
    replace $_, $replace_sub for ("packaging/rpm/$PACKAGE_NAME*.spec", "packaging/deb*/debian/rules");

    path('packaging/deb/debian/source/include-binaries')->spew(join("\n", map { "debian/$PLUGIN_PREFIX$_" } @plugins) . "\n");
}

sub update_packaging_binaries_list {
    my @plugins = @_;
    path('packaging/plugin-lists')->spew(join("\n", map { "$PLUGIN_PREFIX$_" } @plugins) . "\n");
}

sub load_packaging_confg {
    decode_json path('packaging/config.json')->slurp;
}

sub main {
    my @plugins = retrieve_plugins;
    update_readme(@plugins);
    my $config = load_packaging_confg;
    update_packaging_specs(sort @{ $config->{plugins} });
    update_packaging_binaries_list(sort @{ $config->{plugins} });
}

main();
