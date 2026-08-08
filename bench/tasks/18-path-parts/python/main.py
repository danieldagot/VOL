import os

path = os.path.join("reports", "2026", "q1.json")
print(path.replace("\\", "/"))
print(os.path.basename("reports/2026/q1.json"))
print(os.path.dirname("reports/2026/q1.json").replace("\\", "/"))
print(os.path.splitext("reports/2026/q1.json")[1])
