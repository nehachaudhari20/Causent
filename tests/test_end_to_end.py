from fastapi.testclient import TestClient

from api.main import app

client = TestClient(app)


def test_health() -> None:
    r = client.get("/health")
    assert r.status_code == 200
    assert r.json()["status"] == "ok"


def test_create_case_stub_ledger() -> None:
    r = client.post("/v1/cases", json={"policy_version": "1.0"})
    assert r.status_code == 200
    body = r.json()
    assert body["case_id"].startswith("case_")
    assert body["ledger"]["mode"] == "stub"
