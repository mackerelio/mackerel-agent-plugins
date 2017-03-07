use 5.014;
use warnings;
use utf8;
use File::Spec;
use Test::More;

for my $dir (<mackerel-plugin-*>) {
    my $maingo = File::Spec->catfile($dir, 'main.go');
    ok -f -r $maingo or diag "$maingo not found";
    my $readmemd = File::Spec->catfile($dir, 'README.md');
    ok -f -r $readmemd or diag "$readmemd is not available.";

    my $package = $dir;
       $package =~ s/(mackerel-plugin)?-//g;
       $package = "mp$package";
    my $expect = qq[package main

import "github.com/mackerelio/mackerel-agent-plugins/$dir/lib"

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
