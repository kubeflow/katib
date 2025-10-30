import pytest
from kubeflow.katib.utils import utils


@pytest.mark.parametrize(
    "pip_index_urls, expected",
    [
        (["https://pypi.org/simple"], "--index-url https://pypi.org/simple"),
        (
            ["https://pypi.org/simple", "https://private-repo.com/simple"],
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple",
        ),
        (
            [
                "https://pypi.org/simple",
                "https://private-repo.com/simple",
                "https://another-repo.com/simple",
            ],
            "--index-url https://pypi.org/simple --extra-index-url https://private-repo.com/simple "
            "--extra-index-url https://another-repo.com/simple",
        ),
    ],
)
def test_format_pip_index_urls(pip_index_urls, expected):
    assert utils.format_pip_index_urls(pip_index_urls) == expected
