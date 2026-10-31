import requests
import uuid

TARGET_URL = "https://dealzo.onrender.com"

def test_create_then_delete_listing():
    # Create a listing
    payload = {
        "title": f"Delete Me {uuid.uuid4()}",
        "description": "temp listing for delete test",
        "price": 500,
        "city": "Mumbai"
    }
    created = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert created.status_code == 201, f"create failed: {created.status_code} {created.text}"
    listing_id = created.json()["id"]

    # Verify it is present in the listing feed
    feed = requests.get(f"{TARGET_URL}/listings")
    assert feed.status_code == 200
    assert any(item["id"] == listing_id for item in feed.json())

    # Delete it
    deleted = requests.delete(f"{TARGET_URL}/listings/{listing_id}")
    assert deleted.status_code == 204, f"delete failed: {deleted.status_code} {deleted.text}"

    # Verify it is gone
    feed_after = requests.get(f"{TARGET_URL}/listings")
    assert feed_after.status_code == 200
    assert all(item["id"] != listing_id for item in feed_after.json())
