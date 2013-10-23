import requests
import time
import json
from bs4 import BeautifulSoup


def update(page):
    print "http://www.upworthy.com/page/{}".format(page)

    titles = json.load(open("titles.json"))

    resp = requests.get("http://www.upworthy.com/page/{}".format(page))
    resp.raise_for_status()

    soup = BeautifulSoup(resp.content)
    
    shared = soup.find('div', attrs={'id':'recently-shared'})
    stories = shared.find_all('div', attrs={'class':'nugget-info'})

    for div in stories:
        text = div.find('h3')

        if not text.text:
            continue

        titles[text.text.strip()] = ""

    print len(titles)

    json.dump(titles, open("titles.json", "w"))

page = 1

while True:
    update(page)
    page = page + 1
    time.sleep(1)
