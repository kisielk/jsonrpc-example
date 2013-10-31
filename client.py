import urllib2
import json

url = 'http://localhost:8080/rpc'

data = json.dumps({
    'id': 1,
    'method': 'Arith.Multiply',
    'params': [{'A': 1, 'B': 2}]
})

response = urllib2.urlopen(urllib2.Request(url, data))
print json.loads(response.read())
