temps = [22, 18, 25, 31, 29, 17, 24, 28, 20, 26]
total = 0
hot = 0
mild = 0
cold = 0
for t in temps:
    total += t
    if t >= 28:
        hot += 1
    elif t >= 20:
        mild += 1
    else:
        cold += 1
avg = total // len(temps)
if hot + mild + cold != len(temps) or avg != 24:
    raise AssertionError("invalid temperature report")
print("Days measured: " + str(len(temps)))
print("Average: " + str(avg))
print("Hot days (28+): " + str(hot))
print("Mild days: " + str(mild))
print("Cold days (<20): " + str(cold))
