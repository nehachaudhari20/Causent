from assurance.report import report_digest


def test_report_digest_is_stable() -> None:
    report = {"case_id": "case_a", "outcome": "CONSISTENT_WITH_EVIDENCE"}
    d1 = report_digest(report)
    d2 = report_digest(report)
    assert d1 == d2
    assert d1.startswith("sha256:")
