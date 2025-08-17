import unittest

from kubeflow.katib.utils import utils


class TestUtils(unittest.TestCase):
    def test_format_pip_index_urls(self):
        # Test with a single URL.
        pip_index_urls = ["https://pypi.org/simple"]
        self.assertEqual(
            utils.format_pip_index_urls(pip_index_urls),
            "--index-url https://pypi.org/simple",
        )

        # Test with multiple URLs.
        pip_index_urls = ["https://pypi.org/simple", "https://private-repo.com/simple"]
        self.assertEqual(
            utils.format_pip_index_urls(pip_index_urls),
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple",
        )

        # Test with three URLs.
        pip_index_urls = ["https://pypi.org/simple", "https://private-repo.com/simple", "https://another-repo.com/simple"]
        self.assertEqual(
            utils.format_pip_index_urls(pip_index_urls),
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple --extra-index-url https://another-repo.com/simple",
        )

        # Test with default value.
        self.assertEqual(
            utils.format_pip_index_urls(),
            "--index-url https://pypi.org/simple",
        )


if __name__ == "__main__":
    unittest.main()
