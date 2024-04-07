# OMIP

OMIP - An Eve Online Data Aggregator


## Overview

* omip shows corp main member activity over time: kills, losses, bounties
 * shows loss details and zkill link when clicking on character cell in monthly table
* omip shows structure fuel and structure modules
 * shows fuel warning during update
 * shows attack warning if attacked
 * shows structure state chanegs during update
* on update you get an overview on what has changed on which account.
* shows journal entries aggregated into per-day-activity and detailed list in split view
* shows character and corp industry jobs and contracts



## Getting Started

![image](https://user-images.githubusercontent.com/20628481/190867641-e7166a31-fae0-461e-b4cf-9753ab165754.png)

* to register an ESI key for a character or corporation select "Add Character" in the file menu <2>
  * this will open the Authentification screen in a browser window
  * if you want to add a corporation just choose a character who is a director. In this way the character and the corporation will be registered
  * adding a character or corporation will automatically trigger an update for it
* to fetch latest data from the eve server select "Update all data" in the file menu <3>
  * you can only do this every 5mins to avoid flooding the server.


## Notifications

![image](https://user-images.githubusercontent.com/20628481/190867502-6ba1476f-dcc5-42bb-a1a9-a6d73b4929d3.png)

* if the notifications tab is selected you we see notifications in area <1>
  * everytime changes to the database are performed a summary of the change will be displayed in the notification tab
  * for example if the wallet changed on one of your registered characters / corporations you will see the change here
  * this will also inform you about important events like fuel going low on a structure or a structure being attacked 
  * to display the latest notification automatically on startup add omip_cmd to your autostart
  * * on windows10 press windows+R key and type shell:startup then drag and drop omip_cmd.exe into the window

## My Characters

![image](https://user-images.githubusercontent.com/20628481/190867964-3c041cdc-595c-4cd3-a921-bf4b1e6cd45f.png)

* after adding a character it will apear in the "My characters" tab
* if you have a huge amount of characters you can filter the list to find one <8>
* to get character details click on the blue highlighted character name <7> to open the character tab for it
  * the character tabs contains sub tabs for detected character activities (for example Industry, Contrats, Journal)
  * if a character does not have data for this activity (for example a character does not have any contracts) the tab will not be visible
  
 ### Industry
 
 ![image](https://user-images.githubusercontent.com/20628481/190868246-c03d0d7a-ae36-41c1-a4ef-8c9f9a6c838f.png)

  * The Industry tab <9> shows a list of all active industry jobs for this character
  * the filter entries <13> are just text filters 
    * additionally if you enter a number in the the status filter only jobs with a higher percentage of this value will be shown
  * COPY CSV <10> will copy the currently filtered table into the clipboard
  * Reset Filter <12> will clear all filters
  * the number of totally displaeyed table entries is shown in <11>
  
### Contracts
   
![image](https://user-images.githubusercontent.com/20628481/190868666-acf135d4-d7cd-4547-b7c1-9469af6bbc5e.png)

* Opening the contracts tab <14> will be similar to the industry tab 
* The items field contains "multi" highlighted in blue if a contract contains multiple items
  * in this case clicking on the blue items field a table with all contract items will be shown in a separate pop up window

### Journal
   
![image](https://user-images.githubusercontent.com/20628481/190868652-5cee5d1f-2487-4b3d-8a13-60dec1f959ea.png)
   
* The journal tab <15> is separated into the daily transactions section <16> and the detailed section <17>
  * for each day all transactions from the same type are summed up for example the entry <16> shows the sum of all market_escrow transaction on that day 
  * clicking on the row <16> will show the list of all summed up transactions for this day on the right side <17> 
    * these are the same entries you see in-game
* the time interval which is displayed can be selected with <18>
* you can separate the transaction into income and expenses with <19>
* the sum of all filtered transaction is shown in <20>

### Transactions

* transactions only appear in the notification window <1> if a new transaction has been found during update <3>
* examples for transaction notifications are
  * [OMIP] Wilm0rien sold 1 of units Procurer for 35.00M
  * [OMIP] AresDragon bought 1 of units Heat Sink I for 0.06M

## Corporations

* all characters who are directors of a corp will trigger this corp to be added to the corporation tab
* to open a corporation tab blick on the blue highlighted corporation name <22>
* to find the corporation in a huge list use the filter entries <23>
* the wallet total of all filtered corporations is shown in <22>

### ALTs

![image](https://user-images.githubusercontent.com/20628481/190868959-c2200cfe-29e9-4688-ac1b-4f933f2c1465.png)

* all corporation main characters a listed in <35>
* by default all characters in your corporation are defined as mains. 
* however for the statistic it is important to group together all activities from a main and all his alt characters
  * the alt mapping is optional, so you may skip this chapter
* to remove a character from the main list to re-assign him as an alt, click the x button behind the character name <36>
  * now this character appears in the unassiged character list <38>
  * note that this is an intermediate state, a character should either be a Main or an Alt. after unassigning him you have to perform the next step below.

 ![image](https://user-images.githubusercontent.com/20628481/190869060-8def9c46-25d5-456d-ae3f-65a06ed761e1.png)

* to assign an unassign ALT to a main character select the main in the main list <37>
* then click the alt button for the ALT which you want to assign <38>
* this character will be moved to the ALT list on the right side

![image](https://user-images.githubusercontent.com/20628481/190869198-93fc74f7-c1c9-4f02-a1a6-70596ae63ef2.png)

* when you select a main character <39> you will see all his ALTs in the character specific Alt list <40>
* additionally you can see all assigned alts in the assigned Alt list <42>
* pressing the x on either of this list will put the character back into the un assigned alt list from where he can be set as a main or assigned as an alt to a different main

![image](https://user-images.githubusercontent.com/20628481/190897047-32333676-7634-477a-a8fd-16b6c2510e1a.png)

* finally all Alts should be assigend
  * the main list only contains mains <43>
  * the unassigend alt list is empty <44>
  * the alts of the currently selected main (in <43>) are shown in <45>
  * the list of all alts is filled <46>

* for large corporations all lists can be filtered <47>
* after changing the assignement of your characters press "Safe to DB" <48> to safe the changes and refresh the corp statistic with the new alt mapping
* press "Restore from DB" <48> to restore the mapping to the last safed state
  
* in this state the corp statistic will only show the main activity
* if an alt collects bounties or creates kill mails they will be assigned to its main

### Kill Count

![image](https://user-images.githubusercontent.com/20628481/190897224-aa88e803-aecf-4f55-8153-da6cf548b3b0.png)

* the "Kill Count" tab <23> shows on how many corp kill mails a main character + his alts appeared
* if two characters attack the same victim both of them get the kill count. 
* if a character appears on a killmail from another corp, this kill is not counted because it does not appear in the character corporations kill list

### Bounties

![image](https://user-images.githubusercontent.com/20628481/190898569-fd6e1f87-bbee-44a3-9a6e-a174db409eb0.png)

* all ISK which comes from bounties and ESS (Encounter Surveillance System) is summed up for each main and his assigned alt characters
* The Bounties tab <24> shows how many isk has been collected for each main in this way per month
* you can filter by character name  <25>
* the second filter field in <25> accepts an amount of ISK in millions: 1 = 1 Million or 0.25 = 250000
  * only character who have at least earned this amount of isk will be displayed in this way

### ISK Loss

![image](https://user-images.githubusercontent.com/20628481/190902412-4e8c8841-22ff-461b-9869-db6a26c900ed.png)

* the sum of all lost ships from kill mails is summed up for each main and his assigned alt characters
* click on a number <27> to open the list of kill mails in s separate window
* the sum of all losses is summed up on the bottom <28> 
* you can filter for a character with the "filter char name" entry field <28>
* you can filter all total losses which are higher than a number in millons via the "filter millions" entry field <28>
* if a character loses a ship you will get a notification in the notifications tab <1> on update <3>

### Structures

![image](https://user-images.githubusercontent.com/20628481/190903099-a8494e69-83ba-49a1-83a6-b2b80f318920.png)

* the structures tab <29> shows a list if all sturcture beloning to the corporation 
* clicking on a structure <30> will show all active services and rigs of this structure on the right hand side <31>
* if a structure changes its state by being attach you will get a message in the notifications screen <1> on upate <3>
* you will get a low fuel warning for every structure which is running out of fuel within the next 6 days
  * this warning will be shown in the notification tab <1> on update <3>
  * for example "[FYDYN] Barkrik - Red Dwarf fuel expires in   5d  2h"

## ESI Keys

<img width="875" alt="image" src="https://github.com/Wilm0rien/omip/assets/20628481/17b54771-afed-445c-8eed-c634b55a9e73">


* the "ESI Keys" tab <32> shows the list of registered characters
* unchecking a checkbox <33> for a character  will exclude this data from being refreshed on update <3>
  * ths is useful if you have a lot of activity in some area but you do not want to track it with this tool.
  * this may safe a lot of time and traffic on update <3>
* you can delete the keys via delete button <34>
  * this will not delete the data stored for this character in the database
  * if you delete the last director of a corporation this corporation will not be visible anymore in the corporations tab
* when ever you close the tool <35> all your settings are saved 

## Moon Mining tracker

<img width="1176" alt="image" src="https://github.com/Wilm0rien/omip/assets/20628481/3f8e980c-2ee6-43c2-8677-ef9cc10a1626">

* All mining activities from all corporation structures from the past 12 month are listed in the mining tracker overview <36>
* the character names in <36> are prefixed with alliance ticker and corporation ticker and can be filtered with the filter entry <37>
* in "Group by Char" <38> view mode if the characters belong your own corporation the alt mapping is used to sum up all alts under the main character name
* with filter millions <39> only those table entries which have more than the given million value are listed. if you want to filter with less than millions you can set a floating point number like 0.1 to filter out all entries which have less than 100k
* the ISK/ORE selection <40> allows to switch between ISK and ore values. the Ore values are in volume (m3) and the ISK values are calculated in the same way as https://www.fuzzwork.co.uk/ore/ by reading the latest marked values for the minerals contained inside the ores
* if you click on a cell inside the table (for example <44>) the mining detail for this cell is listed in a separate window. in this way you see the structure name and what has been mined.
* if you click the copy CSV button <41> you get all the details mentioned for <44> copied into the clipboard for all data visible in the overview (filters apply)
* if you want to tax full coperations/alliances rather than single players you can set the selection field <38> to corporation and alliance view. this will summarize all characters from that corp/ally in the over view and detail table underneath.
* if you want to see only a percentage of the full value for taxing purposes you can set a percentage value in field 43 and click the update % button to apply the change. you may also use floating point values like 1.3 here. 


# COPYRIGHT NOTICE
EVE Online and the EVE logo are the registered trademarks of CCP hf. All rights are reserved worldwide. All other trademarks are the property of their respective owners. EVE Online, the EVE logo, EVE and all associated logos and designs are the intellectual property of CCP hf. All artwork, screenshots, characters, vehicles, storylines, world facts or other recognizable features of the intellectual property relating to these trademarks are likewise the intellectual property of CCP hf. CCP hf. has granted permission to OMIP to use EVE Online and all associated logos and designs for promotional and information purposes on its website but does not endorse, and is not in any way affiliated with, OMIP. CCP is in no way responsible for the content on or functioning of this website, nor can it be liable for any damage arising from the use of this website.



