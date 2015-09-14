#!/usr/bin/perl -w

use strict;
use warnings;
use Class::Struct;

struct(Project => [ changelists => '$',
                    delta => '$',
                    added => '$',
                    deleted => '$',
                    changed => '$',
                    tags => '@', ]);

my @filenames = (
  "/Users/bergsieker/Desktop/changes.txt",
);

my %projects = ();

sub newProject {
  return Project->new(changelists=>0,
                      delta=>0,
                      added=>0,
                      deleted=>0,
                      changed=>0);
}

sub finalizeChangelist {
  my $changelist = shift;
  if (!defined($changelist->changelists)) {
    $changelist->changelists(0);
  }
  if (!defined($changelist->tags) || !defined($changelist->tags(0))) {
    $changelist->tags(["#undefined"]);
  }
  if (!defined($changelist->delta)) {
    $changelist->delta(0);
  }
  if (!defined($changelist->added)) {
    $changelist->added(0);
  }
  if (!defined($changelist->deleted)) {
    $changelist->deleted(0);
  }
  if (!defined($changelist->changed)) {
    $changelist->changed(0);
  }
  foreach my $tag (@{$changelist->tags}) {
    my $project = $projects{$tag};
    if (!defined($project)) {
      $project = newProject;
    }
    $project->changelists($project->changelists + 1);
    $project->delta($project->delta + $changelist->delta);
    $project->added($project->added + $changelist->added);
    $project->deleted($project->deleted + $changelist->deleted);
    $project->changed($project->changed + $changelist->changed);
    $projects{$tag} = $project;
  }
}

for my $filename (@filenames) {
  if (-e $filename) {
    open FILE, $filename or die "Could not find file : $!";
    my $current;
    foreach my $line (<FILE>) {
      if ((my $cl) = $line =~ /^Change (\d+) /) {
        if (defined($current)) {
          finalizeChangelist($current);
        }
        $current = Project->new;
        $current->changelists($cl);
      } elsif ($line =~ /^#/ && defined($current)) {
        my @tags = split(" ", $line);
        $current->tags(\@tags);
      } elsif ($line =~ /\s+DELTA=/ && defined($current)) {
        if ((my $delta, my $added, my $deleted, my $changed) = $line =~ /\s+DELTA=(\d+)\s+\((\d+)\s+added,\s+(\d+)\s+deleted,\s+(\d+)\s+changed\)$/) {
          $current->delta($delta);
          $current->added($added);
          $current->deleted($deleted);
          $current->changed($changed);
        } else {
          print "DELTA line does not match regex: $line\n";
        }
      }
    }
    if (defined($current)) {
      finalizeChangelist($current);
    }
  }
}

foreach my $tag (sort keys %projects) {
  my $project = $projects{$tag};
  print "PROJECT: $tag\n";
  print "  CHANGES: " . $project->changelists . "\n";
  print "  DELTA=" . $project->delta . " (" . $project->added . " added, " . $project->deleted . " deleted, " . $project->changed . " changed)\n";
}
