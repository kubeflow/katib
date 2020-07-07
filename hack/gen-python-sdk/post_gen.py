import os
import shutil
import sys


def _rewrite_helper(input_file, output_file, rewrite_rules):
    rules = rewrite_rules or []
    lines = []
    with open(input_file, 'r') as f:
        while True:
            line = f.readline()
            if not line:
                break
            for rule in rules:
                line = rule(line)
            lines.append(line)
    with open(output_file, 'w') as f:
        f.writelines(lines)


def update_python_sdk(src, dest, version='v1beta1'):
    # tiny transformers to refine generated codes
    rewrite_rules = [
        lambda l: l.replace('import katib', 'from kubeflow import katib'),
        lambda l: l.replace('from katib', 'from kubeflow.katib'),
        lambda l: "" if l.lstrip().startswith('# noqa') else l
    ]

    src_dirs = [
        os.path.join(src, 'katib', 'models'),
        os.path.join(src, 'test'),
        os.path.join(src, 'docs')
    ]
    dest_dirs = [
        os.path.join(dest, 'kubeflow', 'katib', 'models'),
        os.path.join(dest, 'test'),
        os.path.join(dest, 'docs')
    ]

    for src_dir, dest_dir in zip(src_dirs, dest_dirs):
        # remove previous generated files explicitly, in case of deprecated instances
        for file in os.listdir(dest_dir):
            if version in file.lower():
                os.remove(os.path.join(dest_dir, file))
        # fill latest generated files
        for file in os.listdir(src_dir):
            in_file = os.path.join(src_dir, file)
            out_file = os.path.join(dest_dir, file)
            _rewrite_helper(in_file, out_file, rewrite_rules)
    # clear working dictionary
    shutil.rmtree(src)


if __name__ == '__main__':
    update_python_sdk(src=sys.argv[1],
                      dest=sys.argv[2],
                      version='v1beta1' if len(sys.argv) < 4 else sys.argv[3])
