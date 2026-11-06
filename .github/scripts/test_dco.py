#!/usr/bin/env python3
import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location("dco", Path(__file__).with_name("check-dco.py"))
dco = importlib.util.module_from_spec(spec)
spec.loader.exec_module(dco)

class DCOTests(unittest.TestCase):
    def test_matching_trailer(self):
        self.assertTrue(dco.signed_off("Example Contributor", "123+dev@users.noreply.github.com", "Fix spacing\n\nSigned-off-by: Example Contributor <123+dev@users.noreply.github.com>\n"))
    def test_missing_or_wrong_author(self):
        self.assertFalse(dco.signed_off("Dev", "dev@example.com", "Fix spacing\n"))
        self.assertFalse(dco.signed_off("Dev", "dev@example.com", "Fix\n\nSigned-off-by: Someone Else <else@example.com>\n"))
    def test_body_text_is_not_a_trailer(self):
        self.assertFalse(dco.signed_off("Dev", "dev@example.com", "Signed-off-by: Dev <dev@example.com>\n\nThis is an example, not a sign-off.\n"))
    def test_another_trailer_can_follow(self):
        self.assertTrue(dco.signed_off("Dev", "dev@example.com", "Fix\n\nSigned-off-by: Dev <dev@example.com>\nReviewed-by: Reviewer <reviewer@example.com>\n"))

if __name__ == "__main__":
    unittest.main()
