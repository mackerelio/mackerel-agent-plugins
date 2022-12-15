#!/usr/bin/env perl

use 5.014;
use strict;
use warnings;
use utf8;

use IO::File;
use JSON::PP qw/decode_json/;

my $PLUGIN_PREFIX = 'mackerel-plugin-';
my $PACKAGE_NAME = 'mackerel-agent-plugins';

# refer Mackerel::ReleaseUtils
sub replace {
    my ($glob, $code) = @_;
    for my $file (glob $glob) {
        my $content = $code->(slurp_utf8($file), $file);
        $content .= "\n" if $content !~ /\n\z/ms;

        # for keeping permission
        append_file($file, $content);
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

    write_file(
        'packaging/deb/debian/source/include-binaries',
        join("\n", map { "debian/$PLUGIN_PREFIX$_" } @plugins) . "\n"
    );
}

sub update_packaging_binaries_list {
    my @plugins = @_;
    write_file(
        'packaging/plugin-lists',
        join("\n", map { "$PLUGIN_PREFIX$_" } @plugins) . "\n"
    );
}

# file utility
sub slurp_utf8 {
    my $filename = shift;
    my $fh = IO::File->new($filename, "<:utf8");
    local $/;
    <$fh>;
}
sub write_file {
    my $filename = shift;
    my $content = shift;
    my $fh = IO::File->new($filename, ">:utf8");
    print $fh $content;
    $fh->close;
}
sub append_file {
    my $filename = shift;
    my $content = shift;
    my $fh = IO::File->new($filename, "+>:utf8");
    print $fh $content;
    $fh->close;
}

sub load_packaging_confg {
    decode_json(slurp_utf8('packaging/config.json'));
}

sub main {
    my @plugins = retrieve_plugins;
    update_readme(@plugins);
    my $config = load_packaging_confg;
    update_packaging_specs(sort @{ $config->{plugins} });
    update_packaging_binaries_list(sort @{ $config->{plugins} });
}

main();
