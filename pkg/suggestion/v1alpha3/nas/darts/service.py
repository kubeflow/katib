
import torch.nn as nn
import torchvision.datasets as dset
import torchvision.transforms as transforms
import torch
from pkg.suggestion.v1alpha3.nas.darts.darts_model import NetworkCNN
from pkg.suggestion.v1alpha3.nas.darts.architect import Architect
from pkg.suggestion.v1alpha3.nas.darts.utils import AverageMeter, accuracy

w_lr = 0.025
w_lr_min = 0.001
w_momentum = 0.9
w_weight_decay = 3e-4
w_grad_clip = 5.

alpha_lr = 3e-4
alpha_weight_decay = 1e-3

batch_size = 64
num_workers = 4

epochs = 50


def main():
    init_channels = 16
    num_layers = 8

    # Get dataset with meta information
    input_channels, num_classes, train_data = get_dataset()

    criterion = nn.CrossEntropyLoss()

    model = NetworkCNN(init_channels, input_channels,  num_layers, num_classes, criterion)

    # Weights optimizer
    w_optim = torch.optim.SGD(model.weights(), w_lr,  momentum=w_momentum, weight_decay=w_weight_decay)

    # Alphas optimizer
    alpha_optim = torch.optim.Adam(model.alphas(), alpha_lr, betas=(0.5, 0.999), weight_decay=alpha_weight_decay)

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
        epochs,
        eta_min=w_lr_min)

    architect = Architect(model, w_momentum, w_weight_decay)

    # Start training
    best_top1 = 0.

    for epoch in range(epochs):
        lr_scheduler.step()
        lr = lr_scheduler.get_lr()[0]

        model.print_alphas()

        # Training
        train(train_loader, valid_loader, model, architect, w_optim, alpha_optim, lr, epoch)

        # Validation
        cur_step = (epoch + 1) * len(train_loader)
        top1 = validate(valid_loader, model, epoch, cur_step)

        # Print genotype
        genotype = model.genotype()
        print("Model genotype = {}".format(genotype))

        # Modify best top1
        if top1 > best_top1:
            best_top1 = top1
            best_genotype = genotype

    print("Final best Prec@1 = {:.4%}".format(best_top1))
    print("Best Genotype = {}".format(best_genotype))


def train(train_loader, valid_loader, model, architect, w_optim, alpha_optim, lr, epoch):
    top1 = AverageMeter()
    top5 = AverageMeter()
    losses = AverageMeter()

    cur_step = epoch * len(train_loader)

    model.train()

    for step, ((train_x, train_y), (val_x, val_y)) in enumerate(zip(train_loader, valid_loader)):

        train_size = train_x.size(0)

        # Phase 2. Architect step (Alpha)
        alpha_optim.zero_grad()
        architect.unrolled_backward(train_x, train_y, val_x, val_y, lr, w_optim)
        alpha_optim.step()

        # Phase 1. Child network step (W)
        w_optim.zero_grad()
        logits = model(train_x)
        loss = model.criterion(logits, train_y)

        loss.backward()

        # Gradient clipping
        nn.utils.clip_grad_norm_(model.weights(), w_grad_clip)
        w_optim.step()

        prec1, prec5 = accuracy(logits, train_y, topk=(1, 5))

        losses.update(loss.item(), train_size)
        top1.update(prec1.item(), train_size)
        top5.update(prec5.item(), train_size)

        if step % 50 == 0 or step == len(train_loader) - 1:
            print(
                "Train: [{:2d}/{}] Step {:03d}/{:03d} Loss {losses.avg:.3f} "
                "Prec@(1,5) ({top1.avg:.1%}, {top5.avg:.1%})".format(
                    epoch+1, epochs, step, len(train_loader)-1, losses=losses,
                    top1=top1, top5=top5))

        cur_step += 1

    print("Train: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch+1, epochs, top1.avg))


def validate(valid_loader, model, epoch, cur_step):
    top1 = AverageMeter()
    top5 = AverageMeter()
    losses = AverageMeter()

    model.eval()

    with torch.no_grad:
        for step, (valid_x, valid_y) in enumerate(valid_loader):
            valid_size = valid_x.size(0)

            logits = model(valid_x)
            loss = model.criterion(logits, valid_y)

            prec1, prec5 = accuracy(logits, valid_y, topk=(1, 5))
            losses.update(loss.item(), valid_size)
            top1.update(prec1.item(), valid_size)
            top5.update(prec5.item(), valid_size)

            if step % 50 == 0 or step == len(valid_loader) - 1:
                print(
                    "Train: [{:2d}/{}] Step {:03d}/{:03d} Loss {losses.avg:.3f} "
                    "Prec@(1,5) ({top1.avg:.1%}, {top5.avg:.1%})".format(
                        epoch+1, epochs, step, len(valid_loader)-1, losses=losses,
                        top1=top1, top5=top5))

    print("Valid: [{:2d}/{}] Final Prec@1 {:.4%}".format(epoch+1, epochs, top1.avg))

    return top1.avg


def get_dataset():
    dataset_cls = dset.CIFAR10
    num_classes = 10
    input_channels = 3

    # Do preprocessing
    MEAN = [0.49139968, 0.48215827, 0.44653124]
    STD = [0.24703233, 0.24348505, 0.26158768]
    transf = [
        transforms.RandomCrop(32, padding=4),
        transforms.RandomHorizontalFlip()
    ]

    normalize = [
        transforms.ToTensor(),
        transforms.Normalize(MEAN, STD)
    ]

    train_transform = transforms.Compose(transf + normalize)

    train_data = dataset_cls(root="./data", train=True, download=True, transform=train_transform)

    return input_channels, num_classes, train_data


if __name__ == '__main__':
    main()
