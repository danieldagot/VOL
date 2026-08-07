player1 = [8, 6, 9, 5, 10, 7]
player2 = [7, 8, 6, 9, 8, 8]
p1_total = 0
p2_total = 0
for score in player1:
    p1_total += score
for score in player2:
    p2_total += score
print("Player 1 total: " + str(p1_total))
print("Player 2 total: " + str(p2_total))
