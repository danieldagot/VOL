scores = [72, 95, 81, 64]
print("Students: " + str(len(scores)))
scores[3] = 70
for score in scores:
    if score >= 80:
        print("High score: " + str(score))
