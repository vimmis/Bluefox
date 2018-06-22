The weather app is built in GO and MYSQL_DB
The application folder has necessary libraries and packages, hence no need to download GO and install things.

As the code base is in GO, installation and configuration of Go is required. --Done on the ec2 instance for the project folder path : /home/ubuntu/Bluefox

How to build and Deploy:

1. Clone code from https://github.com/vimmis/Bluefox:

cd /home/ubuntu/
sudo git clone https://github.com/vimmis/Bluefox.git
go get "github.com/codegangsta/negroni"  "github.com/go-sql-driver/mysql" "github.com/gorilla/mux" "github.com/unrolled/render" 

2. Run mysql container:
cd Bluefox
sudo docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.5

3. On above success, run below to first build and then run weather app(on port 3000)
sudo docker build -t weatherapp .
sudo docker run --link mysql:mysql -e  MYSQL_USR="root" -e  MYSQL_PWD="root" -e  MYSQL_DB="weatherapp" -e  MYSQL_TBL="weather" -p 3000:3000 -d weatherapp

4. Run sample Tests:
python3 Test.py

------------------------------

Below are the sample URL REquests:

 GET: 
 http://localhost:3000/pingLB  --For pinging the server, incase Load balanacers are used for future.
 
SAMPLE RESPONSE:
 
{
    "Ping": "Ping Version 1!"
}
------------------------------------------------

POST:

curl -X POST \
  http://localhost:3000/add_temperature_measurement \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"Timestamp":"2018-03-20T15:04:05.000Z",
	"Tempt":"2.2"
}'

SAMPLE RESPONSE:
{
    "Timestamp": "2018-05-12T15:04:05Z",
    "Tempt": 3.2
}
------------------------------------------------

GET:

curl -X GET \
  'http://localhost:3000/get_temperature_measurements?start_timestamp=2018-05-15T15:04:05.000Z&stop_timestamp=2018-05-20T15:04:05.000Z' \
  -H 'Cache-Control: no-cache' 
  
SAMPLE RESPONSE:

[
    {
        "Timestamp": "2018-05-15T15:04:05Z",
        "Tempt": 3.2
    },
    {
        "Timestamp": "2018-05-16T15:04:05Z",
        "Tempt": 2.2
    },
    {
        "Timestamp": "2018-05-18T15:04:05Z",
        "Tempt": 2.2
    }
]
-------------------------------------------------

GET:

curl -X GET \
  'http://localhost:3000/get_average_temperature?start_timestamp=2018-05-15T15:04:05.000Z&stop_timestamp=2018-04-16T15:04:05.000Z' \
  -H 'Cache-Control: no-cache'
  
SAMPLE RESPONSE:
2.7