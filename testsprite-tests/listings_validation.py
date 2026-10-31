import requests

TARGET_URL = "https://dealzo.onrender.com"

def test_create_listing_validation_error():
    # Missing required title -> 422 with the API's error envelope:
    # {"error":{"code":"validation_failed","message":"title: field is required","field":"title"}}
    payload = {"title": "", "description": "", "price": 0, "city": ""}
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 422, f"expected 422, got {r.status_code}: {r.text}"
    body = r.json()
    assert "error" in body
    err = body["error"]
    assert err["code"] == "validation_failed"
    assert err["field"] == "title"
    assert "field is required" in err["message"]
