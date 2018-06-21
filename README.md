The weather app is built in GO and MYSQL_DB
The application folder has necessary libraries and packages, hence no need to download GO and install things.

How to build and Deploy:

1. Run mysql container:
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.5

2. On above success, run below in BluefoxIO folder(where you can find Dockerfile, README ,etc) to first build and then run weather app(on port 3000)
docker build -t weatherapp .
docker run --link mysql:mysql -e  MYSQL_USR="root" -e  MYSQL_PWD="root" -e  MYSQL_DB="weatherapp" -e  MYSQL_TBL="weather" -p 3000:3000 -d weatherapp

On sucess, below are the eligible URL REquests:

 GET: 
 <IP>:3000/pingLB  --For pinging the server, incase Load balanacers are used for future.
 
SAMPLE RESPONSE:
 
{
    "Ping": "Ping Version 1!"
}
------------------------------------------------

POST:

curl -X POST \
  http://192.168.99.100:3000/add_temperature_measurement \
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
  '<IP>:3000/get_temperature_measurements?start_timestamp=2018-05-15T15:04:05.000Z&stop_timestamp=2018-05-20T15:04:05.000Z' \
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
  '<IP>:3000/get_average_temperature?start_timestamp=2018-05-15T15:04:05.000Z&stop_timestamp=2018-04-16T15:04:05.000Z' \
  -H 'Cache-Control: no-cache'
  
SAMPLE RESPONSE:
2.7