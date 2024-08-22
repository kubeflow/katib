# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


import argparse
import json

import numpy as np
import torch
import torch.nn as nn
import utils
from architect import Architect
from model import NetworkCNN
from search_space import SearchSpace


def main():

    parser = argparse.ArgumentParser(description="TrainingContainer")
    parser.add_argument(
        "--algorithm-settings", type=str, default="", help="algorithm settings"
    )
    parser.add_argument(
        "--search-space",
        type=str,
        default="",
        help="search space for the neural architecture search",
    )
    parser.add_argument(
        "--num-layers",
        type=str,
        default="",
        help="number of layers of the neural network",
    )

    args = parser.parse_args()

    # Get Algorithm Settings
    algorithm_settings = args.algorithm_settings.replace("'", '"')
    algorithm_settings = json.loads(algorithm_settings)
    print(">>> Algorithm settings")
    for key, value in algorithm_settings.items():
        if len(key) > 13:
            print("{}\t{}".format(key, value))
        elif len(key) < 5:
            print("{}\t\t\t{}".format(key, value))
        else:
            print("{}\t\t{}".format(key, value))
    print()

    num_epochs = int(algorithm_settings["num_epochs"])

    w_lr = float(algorithm_settings["w_lr"])
    w_lr_min = float(algorithm_settings["w_lr_min"])
    w_momentum = float(algorithm_settings["w_momentum"])
    w_weight_decay = float(algorithm_settings["w_weight_decay"])
    w_grad_clip = float(algorithm_settings["w_grad_clip"])

    alpha_lr = float(algorithm_settings["alpha_lr"])
    alpha_weight_decay = float(algorithm_settings["alpha_weight_decay"])

    batch_size = int(algorithm_settings["batch_size"])
    num_workers = int(algorithm_settings["num_workers"])

    init_channels = int(algorithm_settings["init_channels"])

    print_step = int(algorithm_settings["print_step"])

    num_nodes = int(algorithm_settings["num_nodes"])
    stem_multiplier = int(algorithm_settings["stem_multiplier"])

    # Get Search Space
    search_space = args.search_space.replace("'", '"')
    search_space = json.loads(search_space)
    search_space = SearchSpace(search_space)

    # Get Num Layers
    num_layers = int(args.num_layers)
    print("Number of layers {}\n".format(num_layers))

    # Set GPU Device
    # Currently use only first available GPU
    # TODO: Add multi GPU support
    # TODO: Add functionality to select GPU
    all_gpus = list(range(torch.cuda.device_count()))
    if len(all_gpus) > 0:
        device = torch.device("cuda")
        torch.cuda.set_device(all_gpus[0])
        np.random.seed(2)
        torch.manual_seed(2)
        torch.cuda.manual_seed_all(2)
        torch.backends.cudnn.benchmark = True
        print(">>> Use GPU for Training <<<")
        print("Device ID: {}".format(torch.cuda.current_device()))
        print("Device name: {}".format(torch.cuda.get_device_name(0)))
        print("Device availability: {}\n".format(torch.cuda.is_available()))
    else:
        device = torch.device("cpu")
        print(">>> Use CPU for Training <<<")

    # Get dataset with meta information
    # TODO: Add support for more dataset
    input_channels, num_classes, train_data = utils.get_dataset()

    criterion = nn.CrossEntropyLoss().to(device)

    model = NetworkCNN(
        init_channels,
        input_channels,
        num_classes,
        num_layers,
        criterion,
        search_space,
        num_nodes,
        stem_multiplier,
    )

    model = model.to(device)

    # Weights optimizer
    w_optim = torch.optim.SGD(
        model.getWeights(), w_lr, momentum=w_momentum, weight_decay=w_weight_decay
    )

    # Alphas optimizer
    alpha_optim = torch.optim.Adam(
        model.getAlphas(), alpha_lr, betas=(0.5, 0.999), weight_decay=alpha_weight_decay
    )

    # Split data to train/validation
    num_train = len(train_data)
    split = num_train // 2
    indices = list(range(num_train))

    train_sampler = torch.utils.data.sampler.SubsetRandomSampler(indices[:split])
    valid_sampler = torch.utils.data.sampler.SubsetRandomSampler(indices[split:])

    train_loader = torch.utils.data.DataLoader(
        train_data,
        batch_size=batch_size,
        sampler=train_sampler,
        num_workers=num_workers,
        pin_memory=True,
    )

    valid_loader = torch.utils.data.DataLoader(
        train_data,
        batch_size=batch_size,
        sampler=valid_sampler,
        num_workers=num_workers,
        pin_memory=True,
    )

    lr_scheduler = torch.optim.lr_scheduler.CosineAnnealingLR(
        w_optim, num_epochs, eta_min=w_lr_min
    )

    architect = Architect(model, w_momentum, w_weight_decay, device)

    # Start training
    best_top1 = 0.0

    for epoch in range(num_epochs):
        lr = lr_scheduler.get_last_lr()

        model.print_alphas()

        # Training
        print(">>> Training")
        train(
            train_loader,
            valid_loader,
            model,
            architect,
            w_optim,
            alpha_optim,
            lr,
            epoch,
            num_epochs,
            device,
            w_grad_clip,
            print_step,
        )
        lr_scheduler.step()

        # Validation
        print("\n>>> Validation")
        cur_step = (epoch + 1) * len(train_loader)
        top1 = validate(
            valid_loader, model, epoch, cur_step, num_epochs, device, print_step
        )

        # Print genotype
        genotype = model.genotype(search_space)
        print("\nModel genotype = {}".format(genotype))

        # Modify best top1
        if top1 > best_top1:
            best_top1 = top1
            best_genotype = genotype

    print("Final best Prec@1 = {:.4%}".format(best_top1))
    print("\nBest-Genotype={}".format(str(best_genotype).replace(" ", "")))


def train(
    train_loader,
    valid_loader,
    model,
    architect,
    w_optim,
    alpha_optim,
    lr,
    epoch,
    num_epochs,
    device,
    w_grad_clip,
    print_step,
):
    top1 = utils.AverageMeter()
    top5 = utils.AverageMeter()
    losses = utils.AverageMeter()
    cur_step = epoch * len(train_loader)

    model.train()
    for step, ((train_x, train_y), (valid_x, valid_y)) in enumerate(
        zip(train_loader, valid_loader)
    ):

        train_x, train_y = train_x.to(device, non_blocking=True), train_y.to(
            device, non_blocking=True
        )
        valid_x, valid_y = valid_x.to(device, non_blocking=True), valid_y.to(
            device, non_blocking=True
        )

        train_size = train_x.size(0)

        # Phase 1. Architect step (Alpha)
        alpha_optim.zero_grad()
        architect.unrolled_backward(train_x, train_y, valid_x, valid_y, lr, w_optim)
        alpha_optim.step()

        # Phase 2. Child network step (W)
        w_optim.zero_grad()
        logits = model(train_x)
        loss = model.criterion(logits, train_y)
        loss.backward()

        # Gradient clipping
        nn.utils.clip_grad_norm_(model.getWeights(), w_grad_clip)
        w_optim.step()

        prec1, prec5 = utils.accuracy(logits, train_y, topk=(1, 5))

        losses.update(loss.item(), train_size)
        top1.update(prec1.item(), train_size)
        top5.update(prec5.item(), train_size)

        if step % print_step == 0 or step == len(train_loader) - 1:
            print(
                "Train: [{:2d}/{}] Step {:03d}/{:03d} Loss {losses.avg:.3f} "
                "Prec@(1,5) ({top1.avg:.1%}, {top5.avg:.1%})".format(
                    epoch + 1,
                    num_epochs,
                    step,
                    len(train_loader) - 1,
                    losses=losses,
                    top1=top1,
                    top5=top5,
                )
            )

        cur_step += 1

    print(
        "Train: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch + 1, num_epochs, top1.avg)
    )


def validate(valid_loader, model, epoch, cur_step, num_epochs, device, print_step):
    top1 = utils.AverageMeter()
    top5 = utils.AverageMeter()
    losses = utils.AverageMeter()

    model.eval()

    with torch.no_grad():
        for step, (valid_x, valid_y) in enumerate(valid_loader):
            valid_x, valid_y = valid_x.to(device, non_blocking=True), valid_y.to(
                device, non_blocking=True
            )

            valid_size = valid_x.size(0)

            logits = model(valid_x)
            loss = model.criterion(logits, valid_y)

            prec1, prec5 = utils.accuracy(logits, valid_y, topk=(1, 5))
            losses.update(loss.item(), valid_size)
            top1.update(prec1.item(), valid_size)
            top5.update(prec5.item(), valid_size)

            if step % print_step == 0 or step == len(valid_loader) - 1:
                print(
                    "Validation: [{:2d}/{}] Step {:03d}/{:03d} Loss {losses.avg:.3f} "
                    "Prec@(1,5) ({top1.avg:.1%}, {top5.avg:.1%})".format(
                        epoch + 1,
                        num_epochs,
                        step,
                        len(valid_loader) - 1,
                        losses=losses,
                        top1=top1,
                        top5=top5,
                    )
                )

    print(
        "Valid: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch + 1, num_epochs, top1.avg)
    )

    return top1.avg


if __name__ == "__main__":
    main()
