import requests

TARGET_URL = "https://dealzo.onrender.com"

def test_create_listing_valid():
    payload = {
        "title": "Test Item",
        "description": "Test description",
        "price": 1000,
        "city": "Delhi"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    # Should return 201 Created
    assert r.status_code == 201
    data = r.json()
    assert "id" in data
    assert data["title"] == payload["title"]
    assert "created_at" in data
