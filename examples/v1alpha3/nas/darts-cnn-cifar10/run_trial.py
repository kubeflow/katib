
import torch.nn as nn

import torch
import argparse
import json

from model import NetworkCNN
from architect import Architect
import utils
from search_space import SearchSpace


# TODO: Move to the algorithm settings
# w_lr = 0.025
w_lr = 0.001
w_lr_min = 0.001
w_momentum = 0.9
w_weight_decay = 3e-4
w_grad_clip = 5.

alpha_lr = 3e-4
alpha_weight_decay = 1e-3

batch_size = 64
# num_workers = 4
num_workers = 1

init_channels = 16

print_step = 1

# Set GPU Device
# Currently use only first available GPU
# TODO: Add multi GPU support
# TODO: Add functionality to select GPU
use_gpu = False
device = torch.device("cuda")
all_gpus = list(range(torch.cuda.device_count()))
if len(all_gpus) > 0:
    torch.cuda.set_device(all_gpus[0])
    use_gpu = True
    print(">>> Use GPU for Training <<<")
    print(torch.cuda.current_device())
    print(torch.cuda.get_device_name(0))
    print(torch.cuda.is_available())


def main():

    parser = argparse.ArgumentParser(description='TrainingContainer')
    parser.add_argument('--algorithm-settings', type=str, default="", help="algorithm settings")
    parser.add_argument('--search-space', type=str, default="", help="search space for the neural architecture search")
    parser.add_argument('--num-layers', type=str, default="", help="number of layers of the neural network")

    args = parser.parse_args()

    algorithm_settings = args.algorithm_settings.replace("\'", "\"")
    algorithm_settings = json.loads(algorithm_settings)
    print("Algorithm settings")
    print("{}\n".format(algorithm_settings))
    num_epochs = int(algorithm_settings["num_epoch"])

    search_space = args.search_space.replace("\'", "\"")
    search_space = json.loads(search_space)
    search_space = SearchSpace(search_space)

    num_layers = int(args.num_layers)
    print("Number of layers {}\n".format(num_layers))

    # Get dataset with meta information
    # TODO: Add support for more dataset
    input_channels, num_classes, train_data = utils.get_dataset()

    if use_gpu:
        criterion = nn.CrossEntropyLoss().to(device)
    else:
        criterion = nn.CrossEntropyLoss()

    model = NetworkCNN(init_channels, input_channels,  num_classes, num_layers, criterion, search_space)

    if use_gpu:
        model = model.to(device)

    # Weights optimizer
    w_optim = torch.optim.SGD(model.getWeights(), w_lr,  momentum=w_momentum, weight_decay=w_weight_decay)

    # Alphas optimizer
    alpha_optim = torch.optim.Adam(model.getAlphas(), alpha_lr, betas=(0.5, 0.999), weight_decay=alpha_weight_decay)

    # Split data to train/validation
    num_train = len(train_data)
    split = num_train // 2
    indices = list(range(num_train))

    train_sampler = torch.utils.data.sampler.SubsetRandomSampler(indices[:split])
    valid_sampler = torch.utils.data.sampler.SubsetRandomSampler(indices[split:])

    train_loader = torch.utils.data.DataLoader(train_data,
                                               batch_size=batch_size,
                                               sampler=train_sampler,
                                               num_workers=num_workers,
                                               pin_memory=True)

    valid_loader = torch.utils.data.DataLoader(train_data,
                                               batch_size=batch_size,
                                               sampler=valid_sampler,
                                               num_workers=num_workers,
                                               pin_memory=True)

    lr_scheduler = torch.optim.lr_scheduler.CosineAnnealingLR(
        w_optim,
        num_epochs,
        eta_min=w_lr_min)

    architect = Architect(model, w_momentum, w_weight_decay)

    # Start training
    best_top1 = 0.

    for epoch in range(num_epochs):
        lr_scheduler.step()
        lr = lr_scheduler.get_lr()[0]

        model.print_alphas()

        # Training
        print("Training start")
        train(train_loader, valid_loader, model, architect, w_optim, alpha_optim, lr, epoch, num_epochs)

        # Validation
        cur_step = (epoch + 1) * len(train_loader)
        top1 = validate(valid_loader, model, epoch, cur_step, num_epochs)

        # Print genotype
        genotype = model.genotype(search_space)
        print("\nModel genotype = {}".format(genotype))

        # Modify best top1
        if top1 > best_top1:
            best_top1 = top1
            best_genotype = genotype

    print("Final best Prec@1 = {:.4%}".format(best_top1))
    print("\nBest-Genotype={}".format(best_genotype))


def train(train_loader, valid_loader, model, architect, w_optim, alpha_optim, lr, epoch, num_epochs):
    top1 = utils.AverageMeter()
    top5 = utils.AverageMeter()
    losses = utils.AverageMeter()
    cur_step = epoch * len(train_loader)
    model.train()

    for step, ((train_x, train_y), (valid_x, valid_y)) in enumerate(zip(train_loader, valid_loader)):
        if use_gpu:
            train_x, train_y = train_x.to(device), train_y.to(device)
            valid_x, valid_y = valid_x.to(device), valid_y.to(device)

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
                    epoch+1, num_epochs, step, len(train_loader)-1, losses=losses,
                    top1=top1, top5=top5))
        print("STEP IS {}".format(step))
        cur_step += 1

    print("Train: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch+1, num_epochs, top1.avg))


def validate(valid_loader, model, epoch, cur_step, num_epochs):
    top1 = utils.AverageMeter()
    top5 = utils.AverageMeter()
    losses = utils.AverageMeter()

    model.eval()

    with torch.no_grad():
        for step, (valid_x, valid_y) in enumerate(valid_loader):
            if use_gpu:
                valid_x, valid_y = valid_x.to(device), valid_y.to(device)

            valid_size = valid_x.size(0)

            logits = model(valid_x)
            loss = model.criterion(logits, valid_y)

            prec1, prec5 = utils.accuracy(logits, valid_y, topk=(1, 5))
            losses.update(loss.item(), valid_size)
            top1.update(prec1.item(), valid_size)
            top5.update(prec5.item(), valid_size)

            if step % print_step == 0 or step == len(valid_loader) - 1:
                print(
                    "Train: [{:2d}/{}] Step {:03d}/{:03d} Loss {losses.avg:.3f} "
                    "Prec@(1,5) ({top1.avg:.1%}, {top5.avg:.1%})".format(
                        epoch+1, num_epochs, step, len(valid_loader)-1, losses=losses,
                        top1=top1, top5=top5))

    print("Valid: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch+1, num_epochs, top1.avg))

    return top1.avg


if __name__ == "__main__":
    main()
