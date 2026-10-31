import requests

TARGET_URL = "https://dealzo.onrender.com"

def test_health_check():
    r = requests.get(f"{TARGET_URL}/healthz")
    assert r.status_code == 200
