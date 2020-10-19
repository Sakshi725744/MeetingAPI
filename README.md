# MeetingAPI
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/makes-people-smile.svg)](https://forthebadge.com)
<br>
A Meetings API made using GO-lang and mongoDB. Tested using POSTMAN. API hosted locally on URL<br>
```ruby
localhost:8000/
```
Features:<br>
- [X] Creating new meeting 
```ruby
localhost:8000/meeting/
```
Body<br>
JSON Format:<br>
```ruby
{"ID": "ID4",
"Title": "MEET3",
"arrayParticipants": [{"Name": "Kiran","Email": "267","RSVP":"YES"} , {"Name": "nika","Email": "13l","RSVP":"NO"}],
"start":900,
"end" :1800}
```
- [X] Retriving meeting using meeting ID<br>
```ruby
localhost:8000/meeting/?ID=<ID_NAME>
```
- [X] Retriving meeting using start time and end time<br>
```ruby
localhost:8000/meeting/?start=<START_TIME>&end=<END_TIME>
```
- [X] Retriving meeting using participant Email<br>
```ruby
localhost:8000/meeting/?Email=<EMAIL>
```
To check the working of handler refer to main_test.go
