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


def update_python_sdk(src, dest, versions=('v1alpha3', 'v1beta1')):
    # tiny transformers to refine generated codes
    rewrite_rules = [
        lambda l: l.replace('import katib', 'import kubeflow.katib'),
        lambda l: l.replace('from katib', 'from kubeflow.katib'),
    ]

    src_dirs = [
        os.path.join(src, 'katib'),
        os.path.join(src, 'katib', 'models'),
        os.path.join(src, 'test'),
        os.path.join(src, 'docs')
    ]
    dest_dirs = [
        os.path.join(dest, 'kubeflow', 'katib'),
        os.path.join(dest, 'kubeflow', 'katib', 'models'),
        os.path.join(dest, 'test'),
        os.path.join(dest, 'docs')
    ]

    for src_dir, dest_dir in zip(src_dirs, dest_dirs):
        # remove previous generated files explicitly, in case of deprecated instances
        for file in os.listdir(dest_dir):
            path = os.path.join(dest_dir, file)
            if not os.path.isfile(path):
                continue
            for v in versions:
                if v in file.lower():
                    os.remove(path)
                    break
        # fill latest generated files
        for file in os.listdir(src_dir):
            in_file = os.path.join(src_dir, file)
            out_file = os.path.join(dest_dir, file)
            if not os.path.isfile(in_file):
                continue
            _rewrite_helper(in_file, out_file, rewrite_rules)

    # Update doc for Models README.md
    buffer = []
    update_buffer = []

    # Get data from generated doc
    with open(os.path.join(src, 'README.md'), 'r') as src_f:
        anchor = 0
        for line in src_f.readlines():
            if line.startswith('## Documentation For Models'):
                if anchor == 0:
                    anchor = 1
            elif line.startswith('##') and anchor == 1:
                anchor = 2
            if anchor == 0:
                continue
            if anchor == 2:
                break
            # Remove leading space from the list
            if len(line) > 0:
                line = line.lstrip(" ")
            update_buffer.append(line)
    # Remove latest redundant newline
    update_buffer = update_buffer[:-1]

    # Update README with new models
    with open(os.path.join(dest, 'README.md'), 'r') as dest_f:
        anchor = 0
        for line in dest_f.readlines():
            if line.startswith('## Documentation For Models'):
                if anchor == 0:
                    buffer.extend(update_buffer)
                    anchor = 1
            elif line.startswith('##') and anchor == 1:
                anchor = 2
            if anchor == 1:
                continue
            buffer.append(line)
    with open(os.path.join(dest, 'README.md'), 'w') as dest_f:
        dest_f.writelines(buffer)

    # Clear working dictionary
    shutil.rmtree(src)


if __name__ == '__main__':
    update_python_sdk(src=sys.argv[1], dest=sys.argv[2])
