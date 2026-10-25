from drunix.integration.client import DrunixClient, DrunixSettings
from drunix.integration.models import OpenCaseRequest


def test_drunix_stub_open_case() -> None:
    client = DrunixClient(DrunixSettings(gateway_enabled=False))
    result = client.open_case(OpenCaseRequest(case_id="case_test_001", policy_version="1.0"))
    assert result.success is True
    assert result.mode == "stub"
