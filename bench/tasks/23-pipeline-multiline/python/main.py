nums = [3, 8, 1, 12, 5, 9]
g = [n for n in nums if n > 5]
print(len(g))
print(sum(g))
print(sum(n * 2 for n in g))
