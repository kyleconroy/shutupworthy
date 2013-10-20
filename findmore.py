import requests
import time
import json
from bs4 import BeautifulSoup


def update():
    titles = json.load(open("titles.json"))
    resp = requests.get("http://www.upworthy.com/random")
    resp.raise_for_status()

    soup = BeautifulSoup(resp.content)
    
    title = soup.find('div', attrs={'id':'pagetitle'})
    text = title.find('header')

    if not text.text:
        return

    titles[text.text] = ""
    print text.text
    print len(titles)

    json.dump(titles, open("titles.json", "w"))


while True:
    try:
        update()
    except:
        pass

    time.sleep(1)
