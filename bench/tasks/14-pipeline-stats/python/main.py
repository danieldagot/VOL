nums = [3, 8, 1, 12, 5, 9, 4, 15, 2, 11]
large = [n for n in nums if n > 5]
print(len(large))
print(sum(large))
print(sum(n * 2 for n in large))
print(sum(1 for n in nums if n < 5))
