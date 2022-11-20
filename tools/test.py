#!/usr/bin/env python3
import unittest
import requests
import logging
import binascii
import os
import json
import time


logging.basicConfig(format='%(name)s %(levelname)s %(message)s')
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

HOST = "localhost"
if os.getenv("HOST") != None:
  HOST = os.getenv("HOST")

class TestToDo(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server = f"http://{HOST}:3000"
        cls.repository_server = cls.server + "/repository"
        cls.scan_server = cls.server + "/scan"

    def get_user_id(self):
        return binascii.hexlify(os.urandom(10))

    def test_postgetput(self):
        user_id = self.get_user_id()
        r = requests.post(self.repository_server, params={"user_id": user_id, "repo_name": "react", "repo_url": "https://github.com/facebook/react.git"})
        self.assertEqual(r.status_code, 201)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 1)

        repo_id =  res[0]["id"]
        r = requests.put(self.repository_server, params={"user_id": user_id, "id": repo_id, "repo_name": "new react"})
        self.assertEqual(r.status_code, 200)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        self.assertEqual(len(res), 1)

    def test_delete(self):
        user_id = self.get_user_id()
        r = requests.post(self.repository_server, params={"user_id": user_id, "repo_name": "react", "repo_url": "https://github.com/facebook/react.git"})
        self.assertEqual(r.status_code, 201)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 1)
        
        repo_id = res[0]["id"]
        r = requests.delete(self.repository_server, params={"user_id": user_id, "id": repo_id})
        self.assertEqual(r.status_code, 200)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 0)

    def test_10posts(self):
        user_id = self.get_user_id()
        for i in range(10):
            r = requests.post(self.repository_server, params={"user_id": user_id, "repo_name": "react" + str(i), "repo_url": "https://github.com/facebook/react.git"})
            self.assertEqual(r.status_code, 201)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 10)

    def test_scan(self):
        user_id = self.get_user_id()
        r = requests.post(self.repository_server, params={"user_id": user_id, "repo_name": "techshares", "repo_url": "https://github.com/techsharesteam/techshares"})
        self.assertEqual(r.status_code, 201)

        r = requests.get(self.repository_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 1)
        repo_id =  res[0]["id"]

        r = requests.post(self.scan_server, params={"user_id": user_id, "repo_id": repo_id})
        self.assertEqual(r.status_code, 201)

        #QUEUE
        r = requests.get(self.scan_server, params={"user_id": user_id})
        res = json.loads(r.text)
        self.assertEqual(len(res), 1)
        self.assertEqual(res[0]["findings"], None)
        self.assertEqual(r.status_code, 200)

        #Sleep 5 seconds to scan data
        time.sleep(5)
        r = requests.get(self.scan_server, params={"user_id": user_id})
        res = json.loads(r.text)

        self.assertEqual(len(res), 1)
        self.assertNotEqual(res[0]["findings"], None)
        self.assertEqual(r.status_code, 200)

        
    def test_scan_failed(self):
        user_id_A = self.get_user_id()
        user_id_B = self.get_user_id()

        r = requests.post(self.repository_server, params={"user_id": user_id_A, "repo_name": "react", "repo_url": "https://github.com/facebook/react.git"})
        self.assertEqual(r.status_code, 201)

        r = requests.get(self.repository_server, params={"user_id": user_id_A})
        res = json.loads(r.text)
        self.assertEqual(len(res), 1)
        repo_id =  res[0]["id"]

        r = requests.post(self.scan_server, params={"user_id": user_id_B, "repo_id": repo_id})
        self.assertEqual(r.status_code, 400)

if __name__ == '__main__':
  unittest.main()