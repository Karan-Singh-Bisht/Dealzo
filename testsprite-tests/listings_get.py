import requests

TARGET_URL = "https://dealzo.onrender.com"

def test_get_listings():
    r = requests.get(f"{TARGET_URL}/listings")
    assert r.status_code == 200
    # Should return a JSON array
    data = r.json()
    assert isinstance(data, list)
