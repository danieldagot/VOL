scores = [85, 72, 91, 60, 78, 95, 55, 68, 88, 74]
total = a_grades = b_grades = passing = failing = 0
for s in scores:
    total += s
    if s >= 90:
        a_grades += 1
    if 80 <= s < 90:
        b_grades += 1
    if s >= 60:
        passing += 1
    if s < 60:
        failing += 1
avg = total // len(scores)
print("Class average: " + str(avg))
print("A grades: " + str(a_grades))
print("B grades: " + str(b_grades))
print("Passing: " + str(passing))
print("Failing: " + str(failing))
