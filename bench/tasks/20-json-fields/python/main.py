import json

v = json.loads('{"n":3,"name":"vol"}')
print(v["n"])
print(v["name"])
print(json.dumps({"n": 3}, separators=(",", ":")))
