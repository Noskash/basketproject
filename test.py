import requests

try:
    url = f"https://line31w.bk6bba-resources.com/events/event?lang=en&version=0&eventId=56607721&scopeMarket=1600"
    response = requests.get(url)
    if response.status_code == 200:
        with open(f"jsons2/31.json" , "w" , encoding="utf-8") as f:
            f.write(response.text)
except:
    ("такого нету")