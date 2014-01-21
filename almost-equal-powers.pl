#!/usr/bin/perl

# find powers of different bases which are approximately equal
# e.g. 4**64=340282366920938463463374607431768211456 ~~ 11**37=340039485861577398992406882305761986971 (0.0714%)

use strict;
use warnings;

use POSIX qw(floor ceil);
use bignum;

my @bases = map { Math::BigInt->new($_) } (2, 3, 5, 7, 11);
my @powers = map { Math::BigInt->new($_) } (1 .. 100);
my $low_tolerance = 1.01;
my $high_tolerance = 1 / $low_tolerance;

for my $base1 (@bases) {
    for my $power1 (@powers) {
        my $big1 = pow($base1, $power1);
        for my $base2 (@bases) {
            next unless $base2 > $base1 && # don't duplicate work
                        Math::BigInt::bgcd($base1, $base2) == 1; # only do for relative primes
            my $power2 = $power1 * log($base1) / log($base2);
            my $power2_low = floor $power2;
            my $power2_high = ceil $power2;
            my $big2_low = pow($base2, $power2_low);
            my $big2_high = pow($base2, $power2_high);

            my $ratio_low = $big1 / $big2_low;
            if ($ratio_low < $low_tolerance) {
                printf "%d**%d=%s ~~ %d**%d=%s (%.4f%%)\n",
                    $base1, $power1, $big1->bstr(), $base2, $power2_low, $big2_low->bstr(), 100 * ($ratio_low-1);
            }

            my $ratio_high = $big1 / $big2_high;
            if ($ratio_high > $high_tolerance) {
                printf "%d**%d=%s ~~ %d**%d=%s (%.4f%%)\n",
                    $base1, $power1, $big1->bstr(), $base2, $power2_high, $big2_high->bstr(), 100 * ($ratio_high-1);
            }
        }
    }
}

my %cache;
sub pow {
    my ($base, $power) = @_;
    return $cache{"$base**$power"} //= $base ** $power;
}
