#!/bin/python3
import requests
import json

standings = requests.get('https://algocode.ru/standings_data/bp_fall_2023/')
with open('human_standings.json', 'w', encoding='utf-8') as fl:
    json.dump(standings.json(), fl, indent=2, ensure_ascii=False)
with open('raw_standings.json', 'wb') as fl:
    fl.write(standings.content)