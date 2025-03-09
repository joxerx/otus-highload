import csv
import aiohttp
import aiofiles
import random
import asyncio
import chardet
import requests
import pandas as pd
import aiohttp
import asyncio

# Function to generate a random biography
def generate_biography():
    biographies = [
        "Loves skateboarding",
        "Enjoys reading books",
        "Passionate about music",
        "Avid traveler",
        "Food enthusiast",
        "Fitness freak",
        "Tech geek",
        "Movie buff",
        "Nature lover",
        "Art admirer"
    ]
    return random.choice(biographies)

# Function to generate a random password
def generate_password(length=8):
    characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    password = ''.join(random.choice(characters) for i in range(length))
    return password


# Function to read CSV and send asynchronous requests
async def send_requests(file_path):
    # Read the CSV file
    df = pd.read_csv(file_path, header=None, names=['name', 'birthdate', 'city'])
    
    # Prepare the data
    records = []
    for index, row in df.iterrows():
        name_parts = row['name'].split()
        last_name = name_parts[0]
        first_name = name_parts[1]
        birthdate = row['birthdate']
        city = row['city']
        
        payload = {
            "first_name": first_name,
            "last_name": last_name,
            "birthdate": birthdate,
            "biography": generate_biography(),
            "city": city,
            "password": generate_password()
        }
        
        records.append(payload)
        url = 'http://127.0.0.1:8000/user/register'
        response = requests.post(url, json=payload) 

        if index%100==0:
            print("in preparing records", index)
            if response.status_code == 201:
                print(f"Successfully sent data: {payload}")
            else:
                print(f"Failed to send data: {payload}, Status code: {response.status_code}")
    

# Function to send a POST request
async def send_post_request(session, url, payload):
    print("started task ", payload)
    


if __name__ == "__main__":
    csv_file_path = 'popezzzz.csv'
    asyncio.run(send_requests(csv_file_path))





