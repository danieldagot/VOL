vals = [22, 18, 25, 31, 29, 17, 24, 28, 20, 26]
total = hot = mild = cold = 0
for v in vals:
    total += v
    if v >= 28:
        hot += 1
    elif v >= 20:
        mild += 1
    else:
        cold += 1
print(len(vals))
print(total // len(vals))
print(hot)
print(mild)
print(cold)
