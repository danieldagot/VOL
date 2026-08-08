player1 = [8, 6, 9, 5, 10, 7]
player2 = [7, 8, 6, 9, 8, 8]
p1_total = sum(player1)
p2_total = sum(player2)
winner = "Player 1" if p1_total > p2_total else "Player 2"
p1_strong = sum(1 for s in player1 if s >= 8)
p2_strong = sum(1 for s in player2 if s >= 8)
print("Player 1 total: " + str(p1_total))
print("Player 2 total: " + str(p2_total))
print("Winner: " + winner)
print("P1 rounds 8+: " + str(p1_strong))
print("P2 rounds 8+: " + str(p2_strong))
