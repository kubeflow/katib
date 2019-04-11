import os


MANAGER_ADDRESS = "vizier-core"
MANAGER_PORT = 6789
SEARCH_ALGORITHM = os.environ.get("SEARCH_ALGORITHM", "random_search")
