player1 = [8, 6, 9, 5, 10, 7]
player2 = [7, 8, 6, 9, 8, 8]
p1_total = 0
p2_total = 0
for score in player1:
    p1_total += score
for score in player2:
    p2_total += score
winner = "Player 2"
if p1_total > p2_total:
    winner = "Player 1"
p1_strong = 0
p2_strong = 0
for score in player1:
    if score >= 8:
        p1_strong += 1
for score in player2:
    if score >= 8:
        p2_strong += 1
print("Player 1 total: " + str(p1_total))
print("Player 2 total: " + str(p2_total))
print("Winner: " + winner)
print("P1 rounds 8+: " + str(p1_strong))
print("P2 rounds 8+: " + str(p2_strong))
