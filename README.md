# OMIP

OMIP - An Eve Online Data Aggregator


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
* if a contract contains multiple items, only the most valuable item will be displayed in the items field which additionally will be highlighted in blue
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

![image](https://user-images.githubusercontent.com/20628481/190869299-94eabdc2-485a-45db-a9f3-682ee0fa14f8.png)

* finally all Alts should be assigend
  * the main list only contains mains <43>
  * the unassigend alt list is empty <44>
  * the alts of the currently selected main (in <43>) are shown in <45>
  * the list of all alts is filled <46>
  
* in this state the corp statistic will only show the main activity
* if an alt collects bounties or creates kill mails they will be assigned to its main





