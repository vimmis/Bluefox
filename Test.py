import json
import requests

passed=0
failed =0
total= 0


#Test Ping Request
print("Runnung Test Ping Request...")
total=total+1
response = requests.get('http://localhost:3000/pingLB')
json_data1 = json.loads(response.text)
print("Response: \n")
print(json_data1["Ping"])
if json_data1["Ping"] == "Ping Version 1!" :
	print("Ping Test passed")
	passed=passed+1
else:
	print("Ping Test failed")
	failed=failed+1
print("-----------------------------------------------\n")	

#Test Add Measurement Request-Sucess1
print("Running Test Add Measurement Request-Sucess1..." )
total=total+1
url = "http://localhost:3000/add_temperature_measurement"

payload = "{\n\t\"Timestamp\":\"2018-05-16T15:04:05.000Z\",\n\t\"Tempt\":\"5.2\"\n}"
headers = {
    'Content-Type': "application/json",
    'Cache-Control': "no-cache",
    }

response = requests.request("POST", url, data=payload, headers=headers)
json_data2 = json.loads(response.text)
print("Response: \n")
print(json_data2)
if (json_data2["Timestamp"] == "2018-05-16T15:04:05Z") and (json_data2["Tempt"] == 5.2 ) :
	print("Add Test-Success1 passed")
	passed=passed+1
else:
	print("Add Test-Sucess1 failed")
	failed=failed+1
print("-----------------------------------------------\n")		

#Test Add Measurement Request-Sucess2
print("Running Test Add Measurement Request-Sucess2.. ")
total=total+1
payload = "{\n\t\"Timestamp\":\"2018-04-10T15:04:05.000Z\",\n\t\"Tempt\":\"2.5\"\n}"
response = requests.request("POST", url, data=payload, headers=headers)
json_data22 = json.loads(response.text)
print("Response: \n")
print(json_data22)
if (json_data22["Timestamp"] == "2018-04-10T15:04:05Z") and (json_data22["Tempt"] == 2.5 ) :
	print("Add Test-Success2 passed")
	passed=passed+1
else:
	print("Add Test-Sucess2 failed")
	failed=failed+1
print("-----------------------------------------------\n")		

#Test Add Measurement Request-Invalid timestamp
print("Running Test Add Measurement Request-Invalid timestamp...")
total=total+1
payload = "{\n\t\"Timestamp\":\"abcd\",\n\t\"Tempt\":\"2.2\"\n}"

responseq = requests.request("POST", url, data=payload, headers=headers)
print("Response: \n")
print(responseq.text, responseq.status_code)
if ("Null or Invalid timestamp, expected: 2006-01-02T15:04:05.000Z  format" in responseq.text) and (responseq.status_code == 400 ) :
	print("Add Test-Invalid Timestamp passed")
	passed=passed+1
else:
	print("Add Test-Invalid Timestamp failed")
	failed=failed+1
print("-----------------------------------------------\n")		

#Test Add Measurement Request-Invalid Temperature
print("Running Test Add Measurement Request-Invalid Temperature...")
total=total+1
payload = "{\n\t\"Timestamp\":\"2018-07-16T15:04:05.000Z\",\n\t\"Tempt\":\"sdf\"\n}"

responsea = requests.request("POST", url, data=payload, headers=headers)
print("Response: \n")
print(responsea.text)
if ( "Null or Invalid temperature or out-of-bounds"  in responsea.text) and (responsea.status_code == 400 ) :
	print("Add Test-Invalid Temperature passed")
	passed=passed+1
else:
	print("Add Test-Invalid Temperature failed")
	failed=failed+1
print("-----------------------------------------------\n")	

#Test Add Measurement Request-Missing Body
print("Running Test Add Measurement Request-Missing Body...")
total=total+1
payloadd = ""

responses = requests.request("POST", url, data=payloadd, headers=headers)
print("Response: \n")
print(responses.text)
if ("Null or Invalid timestamp, expected: 2006-01-02T15:04:05.000Z  format" in responses.text) and (responses.status_code == 400 ) :
	print("Add Test-Missing Body passed")
	passed=passed+1
else:
	print("Add Test-Missing Body failed")
	failed=failed+1


#Test Get Measurement Range Request-Success
print("Running Test Get Measurement Range Request-Success...")
total=total+1
url = "http://localhost:3000/get_temperature_measurements"
querystring = {"start_timestamp":"2018-04-10T15:04:05.000Z","stop_timestamp":"2018-06-10T15:04:05.000Z"} #includes the two added measurement
headers = {
    'Cache-Control': "no-cache",
    }

responsez = requests.request("GET", url, headers=headers, params=querystring)
json_data3 = json.loads(responsez.text)
print("Response: \n")
print(json_data3)
print(len(json_data3))
if (len(json_data3) == 2) and (responsez.status_code == 200 ) :
	print("Get Measurement Range-Sucess passed")
	passed=passed+1
else:
	print("Get Measurement Range-Sucess failed")
	failed=failed+1
print("-----------------------------------------------\n")	

#Test Get Measurement Range Request-Missing param
print("Running Test Get Measurement Range Request-Missing param...")
total=total+1
querystring = {"stop_timestamp":"2018-06-10T15:04:05.000Z"}#missing start_timestamp
responsew = requests.request("GET", url, headers=headers, params=querystring)
print("Response: \n")
print(responsew.text)

if ("start time missing" in responsew.text) and (responsew.status_code == 400 ):
	print("Get Measurement Range-Missing param passed")
	passed=passed+1
else:
	print("Get Measurement Range-Missing param failed")
	failed=failed+1
print("-----------------------------------------------\n")		

#Test Get Measurement Range Request-Incorrect Timestamp
print("Running Test Get Measurement Range Request-Incorrect Timestamp...")
total=total+1
querystring = {"start_timestamp":"abcd","stop_timestamp":"2018-06-10T15:04:05.000Z"}
responsex = requests.request("GET", url, headers=headers, params=querystring)
print("Response: \n")
print(responsex.text)

if ("Invalid start timestamp, expected: 2006-01-02T15:04:05.000Z  format" in responsex.text) and (responsex.status_code == 400 ) :
	print("Get Measurement Range-Incorrect Timestamp passed")
	passed=passed+1
else:
	print("Get Measurement Range-Incorrect Timestamp failed")
	failed=failed+1
print("-----------------------------------------------\n")	

#Test Get Measurement Range Request-No result
print("Running Test Get Measurement Range Request-No result...")
total=total+1
querystring = {"start_timestamp":"2013-06-10T15:04:05.000Z","stop_timestamp":"2013-06-10T15:04:05.000Z"} #this range has no data in db
response = requests.request("GET", url, headers=headers, params=querystring)
print("Response: \n")
print(response.text)

if ( "null" in response.text) and (response.status_code == 200 ) :
	print("Get Measurement Range-No result passed")
	passed=passed+1
else:
	print("Get Measurement Range-No result failed")
	failed=failed+1
print("-----------------------------------------------\n")		

#Test Get Average Temp Range Request-Sucess
print("Running Test Get Average Temp Range Request-Sucess...")
total=total+1	
url = "http://localhost:3000/get_average_temperature"
querystring = {"start_timestamp":"2018-03-15T15:04:05.000Z","stop_timestamp":"2018-10-20T15:04:05.000Z"}#all inclusive
response = requests.request("GET", url, headers=headers, params=querystring)
print("Response: \n")
print(response.text)

if ("3.85" in response.text) and (response.status_code == 200 ) :
	print("Get Average Temp Range Request-Sucess passed")
	passed=passed+1
else:
	print("Get Average Temp Range Request-Sucess failed")
	failed=failed+1
print("-----------------------------------------------\n")		


#Test Get Average Temp Range Request-Out-of-range
print("Running Test Get Average Temp Range Request-Out-of-range...")
total=total+1	
url = "http://localhost:3000/get_average_temperature"
querystring = {"start_timestamp":"2012-05-15T15:04:05.000Z","stop_timestamp":"2012-05-20T15:04:05.000Z"}
response = requests.request("GET", url, headers=headers, params=querystring)
print("Response: \n")
print(response.text)

if ("0" in response.text) and (response.status_code == 200 ) :
	print("Get  Average Temp Range Request-Sucess passed")
	passed=passed+1
else:
	print("Get  Average Temp Range Request-Sucess failed")
	failed=failed+1
print("-----------------------------------------------\n")
		
print("Results: ")
print("	Passed: ",passed)
print("	Failed: ",failed)
print("	Total: ",total)