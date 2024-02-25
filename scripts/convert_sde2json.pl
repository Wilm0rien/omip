use strict;
use utf8;
use JSON;
use Archive::Zip;
use Archive::Zip::MemberRead;
use YAML::XS 'LoadFile';

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



# extract only the following ores from the typeMaterials.yaml
my @ores = (
	"Cobaltite", "Copious Cobaltite", "Twinkling Cobaltite",
	"Euxenite", "Copious Euxenite", "Twinkling Euxenite",
	"Scheelite", "Copious Scheelite", "Twinkling Scheelite",
	"Titanite", "Copious Titanite", "Twinkling Titanite",

	# Exceptional
	"Loparite", "Bountiful Loparite", "Shining Loparite",
	"Monazite", "Bountiful Monazite", "Shining Monazite",
	"Xenotime", "Bountiful Xenotime", "Shining Xenotime",
	"Ytterbite", "Bountiful Ytterbite", "Shining Ytterbite",

	# Rare
	"Carnotite", "Glowing Carnotite", "Replete Carnotite",
	"Cinnabar", "Glowing Cinnabar", "Replete Cinnabar",
	"Pollucite", "Glowing Pollucite", "Replete Pollucite",
	"Zircon", "Glowing Zircon", "Replete Zircon",

	# Ubiquitous
	"Bitumens", "Brimful Bitumens", "Glistening Bitumens",
	"Coesite", "Brimful Coesite", "Glistening Coesite",
	"Sylvite", "Brimful Sylvite", "Glistening Sylvite",
	"Zeolites", "Brimful Zeolites", "Glistening Zeolites",

	# uncommon
	"Chromite", "Lavish Chromite", "Shimmering Chromite",
	"Otavite", "Lavish Otavite", "Shimmering Otavite",
	"Sperrylite", "Lavish Sperrylite", "Shimmering Sperrylite",
	"Vanadinite", "Lavish Vanadinite", "Shimmering Vanadinite",

	# normal ores
	"Arkonor", "Crimson Arkonor", "Flawless Arkonor", "Prime Arkonor",
	"Bezdnacine", "Hadal Bezdnacine",
	"Bistot", "Cubic Bistot", "Monoclinic Bistot", "Triclinic Bistot",
	"Crokite", "Crystalline Crokite", "Pellucid Crokite", "Sharp Crokite",
	"Dark Ochre", "Jet Ochre", "Obsidian Ochre", "Onyx Ochre",
	"Ducinium", "Imperial Ducinium", "Noble Ducinium", "Royal Ducinium",
	"Eifyrium", "Augmented Eifyrium", "Boosted Eifyrium", "Doped Eifyrium",
	"Gneiss", "Brilliant Gneiss", "Iridescent Gneiss", "Prismatic Gneiss",
	"Hedbergite", "Glazed Hedbergite", "Lustrous Hedbergite", "Vitric Hedbergite",
	"Hemorphite", "Radiant Hemorphite", "Scintillating Hemorphite", "Vivid Hemorphite",
	"Jaspet", "Luminous Kernite", "Pure Jaspet", "Immaculate Jaspet",
	"Kernite", "Pristine Jaspet", "Fiery Kernite", "Resplendant Kernite",
	"Mercoxit", "Vitreous Mercoxit", "Magma Mercoxit", "Resplendant Kernite",
	"Mordunium", "Plum Mordunium", "Plunder Mordunium", "Prize Mordunium",
	"Omber", "Golden Omber", "Platinoid Omber", "Silvery Omber",
	"Plagioclase", "Rich Plagioclase", "Sparkling Plagioclase", "Azure Plagioclase",
	"Pyroxeres", "Solid Pyroxeres", "Viscous Pyroxeres", "Opulent Pyroxeres",
	"Rakovene", "Nesosilicate Rakovene", "Hadal Rakovene", "Abyssal Rakovene",
	"Scordite", "Massive Scordite", "Glossy Scordite", "Condensed Scordite",
	"Spodumain", "Gleaming Spodumain", "Dazzling Spodumain", "Bright Spodumain",
	"Talassonite", "Hadal Talassonite", "Abyssal Talassonite",
	"Veldspar", "Dense Veldspar", "Stable Veldspar", "Concentrated Veldspar",
);

my $id_map;
foreach my $id (keys %{$result_hash})
{
    my $item_name = $result_hash->{$id};
    $id_map->{$item_name} = $id;
}



my $ore_map;

foreach my $ore (@ores)
{
    $ore_map->{$ore} = 1;
}

my $utf8_encoded_json_text = to_json($result_hash, {utf8 => 1, pretty => 1});

print FH $utf8_encoded_json_text;
close(FH);
my $type_materials_file = "sde/fsd/typeMaterials.yaml";
my $type_materials_zip  = Archive::Zip::MemberRead->new($zip, $type_materials_file);
my $type_materials_content = do { local $/; $type_materials_zip->getline() };
my $type_materials_yaml_data = YAML::XS::Load($type_materials_content);


my $mat_types_omip_style;

# omip style in Go looks like this:
#
# type SdeTypeMatsStruct struct {
#   MaterialTypeID int `json:"MaterialTypeID"`
#   Quantity       int `json:"Quantity"`
# }
# 
# type SdeTypeMatsList map[int][]*SdeTypeMatsStruct


foreach my $mat (keys %{$type_materials_yaml_data})
{
  # check if material is known from "sde/fsd/typeIDs.yaml";
  if (not defined $result_hash->{$mat})
  {
    next;
  }

  # check if it is an ore
  my $mat_name = $result_hash->{$mat};
  if (not defined $ore_map->{$mat_name})
  {
    next;
  }

  # check if it contains a materials list
  if (not defined $type_materials_yaml_data->{$mat}->{materials})
  {
    next;
  }

  # increment verify counter 
  $ore_map->{$mat_name}++;
  my @mat_list;
  foreach my $material (@{$type_materials_yaml_data->{$mat}->{materials}})
  {
    # convert strings to integer
    my $part_mat;
    $part_mat->{MaterialTypeID} = int($material->{materialTypeID});
    $part_mat->{Quantity} = int($material->{quantity});
    push @mat_list, $part_mat;
  }
  $mat_types_omip_style->{int($mat)}=\@mat_list;
}

# check very counter if all ores have been found
foreach my $mat_name (keys %{$ore_map})
{
    if ($ore_map->{$mat_name}!=2)
    {
        printf("warning unexpected ore count %d\n", $ore_map->{$mat_name})
    }
}

# write data to json file
open(FH, '>', "type_materials.json_updated") or die $!;
my $utf8_encoded_json_text2 = to_json($mat_types_omip_style, {utf8 => 1, pretty => 1});
print FH $utf8_encoded_json_text2;
close(FH);

