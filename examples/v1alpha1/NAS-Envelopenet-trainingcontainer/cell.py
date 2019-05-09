"""Base cell."""

import tensorflow as tf
class Cell:
    """Base cell: All cells derived from this class."""
    def __init__(self):
        pass
    def get_params(self):
        """Get tf params"""
        print([tensor.name for tensor in tf.get_default_graph().as_graph_def().node])
