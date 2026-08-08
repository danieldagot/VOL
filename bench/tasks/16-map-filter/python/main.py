xs = [1, 2, 3, 4, 5, 6, 7, 8]
mapped = [x * 3 for x in xs]
print(sum(1 for n in mapped if 10 < n < 20))
print(sum(n for n in mapped if n > 10))
