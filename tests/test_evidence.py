from evidence.replay_guard import ReplayGuard


def test_replay_guard_accepts_once() -> None:
    guard = ReplayGuard()
    assert guard.check_and_record("evt-1") is True
    assert guard.check_and_record("evt-1") is False
