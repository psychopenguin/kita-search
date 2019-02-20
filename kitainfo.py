import requests
from bs4 import BeautifulSoup


def retrieve_kita_list():
    source = 'https://www.berlin.de/sen/jugend/familie-und-kinder/'\
             'kindertagesbetreuung/kitas/verzeichnis/ListeKitas.aspx'
    req = requests.get(source)
    return req.text


def parse_kita_list(s):
    """
    this fuction receive a html document as string
    returns a list of kita ids
    """
    content = BeautifulSoup(s)
    kitalist = content.find(id='DataList_Kitas')
    return kitalist
