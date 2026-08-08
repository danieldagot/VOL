revenue = [240, 175, 289, 150, 225, 199, 180, 178]
total = high_sum = high_count = budget_count = 0
for r in revenue:
    total += r
    if r >= 200:
        high_count += 1
        high_sum += r
    else:
        budget_count += 1
print("Total revenue: " + str(total))
print("Premium orders (200+): " + str(high_count))
print("Premium revenue: " + str(high_sum))
print("Budget orders: " + str(budget_count))
