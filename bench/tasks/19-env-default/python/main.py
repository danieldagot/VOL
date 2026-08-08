import os

os.environ["VOL_DENSITY_PORT"] = "9090"
print(os.environ.get("VOL_DENSITY_PORT", "missing"))
print(os.environ.get("VOL_DENSITY_NO_SUCH_KEY", "fallback"))
