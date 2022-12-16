use strict;
use warnings;
use Test::More;

use JSON::PP qw/decode_json/;

my $config = decode_json do {
    local $/;
    open my $fh, '<', 'packaging/config.json' or die $!;
    <$fh>
};
ok $config->{description};

my $plugins_to_be_packaged = $config->{plugins};
isa_ok $plugins_to_be_packaged, 'ARRAY';

my %plugins = map {s/^mackerel-plugin-//; ($_ => 1)} <mackerel-plugin-*>;
for my $plug (@$plugins_to_be_packaged) {
    ok $plugins{$plug}, "$plug ok";
}

my $external_plugins_to_be_packaged = $config->{'external-plugins'};
isa_ok $external_plugins_to_be_packaged, 'ARRAY';

for my $plug (@$external_plugins_to_be_packaged) {
    ok $plug->{name}, "name ok";
    my $name = $plug->{name};
    like $plug->{repository}, qr|^github[.]com/mackerelio/|, "$name repository ok";
}

done_testing;
