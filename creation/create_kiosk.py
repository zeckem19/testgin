import os
import certifi
from datetime import datetime
from pymongo import MongoClient
from dotenv import load_dotenv

load_dotenv()

MURL = os.getenv('MURL')
DBNAME = os.getenv('DBNAME')

client = MongoClient(MURL, tlsCAFile=certifi.where())
db = client[DBNAME]

kiosk_collection = db['kiosk']

kiosk_collection.delete_many({})


patient_collection = db['patient']

patient_collection.delete_many({})

kiosk = {
    'id':'K00000',
    'state':'default',
    'creation_datetime':datetime.now(),
    'name':'TEST KIOSK2',
    'address':'MOCK ADDRESS',
    'image':'',
    'location':'',
    'queue':[],
    'current':None,
}

kiosk_collection.insert_one(kiosk)
