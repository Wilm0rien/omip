# OMIP

OMIP - An Eve Online Data Aggregator with focus on data security and privacy.

## Security Features

* OMIP is a program which runs natively on your computer (Linux or Windows)
* No data is stored outside your computer
* Only the following servers are contacted from this tool
    * esi.evetech.net (for fetching ESI data)
    * github.com (for update checking)
* ESI keys are stored in your local home directory and secured by cryptography
* ESI data is stored separately in a local database file in your home directory
* All source code is on github under GPLv3
* No donations (ISK or otherwise) will be accepted by the maintainers of this
  tool

## Line Member features

If you are a corporation line member you have access to the following tool
features.

* Register any number of characters
* Get notifications about finished industry jobs, marked orders and contracts
* Sort and filter your wallet journal entries

## Corporation Director features

If you are a corporation director you have access to the following additional
tool features.

* Manage any number of alt corps
* Assign alt characters in your corporation to their respective main characters
  to get only main character names listed in all activity tables
* List monthly activity of your line members for Bounties, Kills and Ship Losses
* Get notifications about finished industry jobs, marked orders and contracts
  for your corp
* Get notifications for structure fuel expiry
* Get an overview of all corporation structures and structure services

# FAQ

**Where are my ESI tokens and Data stored exactly?**

* $HOME/omip/ for linux
* %APPDATA%/omip/ for windows

**What do you mean with "tokes are secured by cryptography"?**

You may assume that your local home directory is a save location to store your
private data in general. However, you might want to put this data on an usb 
stick and copy it on another machine because you bought a new computer and want 
to continue filling your monthly statistics from there.

When ever the token data is transferred out of this secure place for whatever 
reason you are in danger that your keys get in the wrong hands.

To protect at least the ESI keys to be only valid on one computer these keys are 
encrypted via AES-256 symmetric encryption. The passphrase for this encryption 
is an SHA-1 hash which is generated from your unique computer id and a salt. 
This means that nobody (not even you) is able to use this key on a different 
computer.

In this way your ESI keys will only be valid on one computer, and you have to
create new keys when moving your database to another machine. 