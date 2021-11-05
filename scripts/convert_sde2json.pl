use strict;
use utf8;
use JSON;
use Archive::Zip;
use Archive::Zip::MemberRead;

my $zipFile = "sde.zip";
my $zip = Archive::Zip->new($zipFile);

my $file = "sde/fsd/typeIDs.yaml";


open(FH, '>', "typeIDs.json_updated") or die $!;


my $info  = Archive::Zip::MemberRead->new($zip, $file);
my $result_hash;
my $current_type_ID=0;
my $name_ok = 0;
while( my $line = $info->getline({ preserve_line_ending => 1 }))  {   
    if ($line=~/^([0-9]+):$/){
      $current_type_ID = $1;
    }
    if ($line=~/^\s+name:$/){
      $name_ok = 1;
    }
    if ($name_ok == 1) {
      if ($line=~/^\s+en: (.*)/){
        $name_ok = 0;
        $result_hash->{$current_type_ID} = $1;
      }
    }
}


#my $utf8_encoded_json_text = encode_json $result_hash;


my $utf8_encoded_json_text = to_json($result_hash, {utf8 => 1, pretty => 1});

print FH $utf8_encoded_json_text;
close(FH);