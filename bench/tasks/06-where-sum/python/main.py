numbers = [4, 7, 2, 9, 12]
large = [n for n in numbers if n > 5]
total = sum(large)
print("Sum: " + str(total))
assert total == 28, "Unexpected collection total."
