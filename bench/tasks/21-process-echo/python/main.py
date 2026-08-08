import subprocess

p = subprocess.run(["echo", "vol"], capture_output=True, text=True, check=False)
print(p.returncode)
print(p.stdout.strip())
