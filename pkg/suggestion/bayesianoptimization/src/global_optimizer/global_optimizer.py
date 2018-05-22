""" module for the global optimizer
DIRECT algorithm is used in this case
"""
import copy
import numpy as np

from pkg.suggestion.bayesianoptimization.src.acquisition_func.acquisition_func import AcquisitionFunc


class RectPack:
    """ class for the rectangular
    including border, center and acquisition function value
    """

    def __init__(self, l, u, division_num, dim, scaler, aq_func):
        self.l = l
        self.u = u
        self.center = (l + u) / 2
        j = np.mod(division_num, dim)
        k = (division_num - j) / dim
        self.d = np.sqrt(j * np.power(3, float(-2 * (k + 1))) + (dim - j) * np.power(3, float(-2 * k))) / 2
        self.division_num = division_num
        self.fc, _, _ = aq_func.compute(scaler.inverse_transform(self.center))
        self.fc = -self.fc


class RectBucket:
    """ class for the rectangular bucket
    rectangular with the same size are put in the same bucket
    the rectangular is sorted by the acquisition function value
    """

    def __init__(self, diff, pack):
        self.diff = diff
        self.array = [pack]

    def insert(self, new_pack):
        """ insert a new rectangular to a bucket """
        for i in range(len(self.array)):
            if new_pack.fc < self.array[i].fc:
                self.array.insert(i, new_pack)
                return
        self.array.append(new_pack)

    def delete(self):
        """ delete the first rectangular"""
        del self.array[0]

    def diff_exist(self, diff):
        """ detect the size difference """
        return abs(self.diff - diff) < 0.00001


class OptimalPoint:
    """ helper class to find potential optimal points"""

    def __init__(self, point, prev, slope):
        self.point = point
        self.prev = prev
        self.slope = slope


class DimPack:
    def __init__(self, dim, fc):
        self.dim = dim
        self.fc = fc


class GlobalOptimizer:
    """ class for the global optimizer """

    def __init__(self, N, l, u, scaler, X_train, y_train, current_optimal, mode, trade_off, length_scale,
                 noise, nu, kernel_type, n_estimators, max_features, model_type):
        self.N = N
        self.l = l
        self.u = u
        self.scaler = scaler
        self.buckets = []
        self.dim = None
        self.aq_func = AcquisitionFunc(
            X_train=X_train,
            y_train=y_train,
            current_optimal=current_optimal,
            mode=mode,
            trade_off=trade_off,
            length_scale=length_scale,
            noise=noise,
            nu=nu,
            kernel_type=kernel_type,
            n_estimators=n_estimators,
            max_features=max_features,
            model_type=model_type,
        )

    def potential_opt(self, f_min):
        """ find the potential optimal rectangular """
        b = []
        for i in range(len(self.buckets)):
            b.append(self.buckets[i].array[0])
        b.sort(key=lambda x: x.d)
        index = 0
        min_fc = b[0].fc
        for i in range(len(b)):
            if b[i].fc < min_fc:
                min_fc = b[i].fc
                index = i

        opt_list = [OptimalPoint(b[index], 0, 0)]
        for i in range(index + 1, len(b)):
            prev = len(opt_list) - 1
            diff1 = b[i].d
            diff2 = opt_list[prev].point.d
            current_slope = (b[i].fc - opt_list[prev].point.fc) / (diff1 - diff2)
            prev_slope = opt_list[prev].slope

            while prev >= 0 and current_slope < prev_slope:
                temp = opt_list[prev].prev
                opt_list[prev].prev = -1
                prev = temp
                prev_slope = opt_list[prev].slope
                diff1 = b[i].d
                diff2 = opt_list[prev].point.d
                current_slope = (b[i].fc - opt_list[prev].point.fc) / (diff1 - diff2)

            opt_list.append(OptimalPoint(b[i], prev, current_slope))

        opt_list2 = []
        for i in range(len(opt_list)):
            if opt_list[i].prev != -1:
                opt_list2.append(opt_list[i])

        for i in range(len(opt_list2) - 1):
            c1 = opt_list2[i].point.d
            c2 = opt_list2[i + 1].point.d
            fc1 = opt_list2[i].point.fc
            fc2 = opt_list2[i + 1].point.fc
            if fc1 - c1 * (fc1 - fc2) / (c1 - c2) > (1 - 0.001) * f_min:
                #         if abs(fc1-fc2)<0.0001:
                opt_list2[i] = None
        while None in opt_list2:
            index = opt_list2.index(None)
            del opt_list2[index]
        # for opt in opt_list2:
        #         print(opt.point.fc)
        return opt_list2

    def direct(self):
        """ main algorithm """
        self.dim = self.l.shape[1]
        division_num = 0

        # create the first rectangle and put it in the first bucket
        first_rect = RectPack(self.l, self.u, division_num, self.dim,
                              self.scaler, self.aq_func)
        self.buckets.append(RectBucket(first_rect.d, first_rect))

        ei_min = []
        f_min = first_rect.fc
        x_next = first_rect.center
        ei_min.append(f_min)

        for t in range(self.N):
            opt_set = self.potential_opt(f_min)

            # for bucket in self.buckets:
            #     for i in range(len(bucket.array)):
            #         print(bucket.array[i].fc)
            #         plt.plot(bucket.diff, bucket.array[i].fc, 'b.')
            #
            # for opt in opt_set:
            #     plt.plot(opt.point.d, opt.point.fc, 'r.')
            # plt.show()

            for opt in opt_set:
                f_min, x_next = self.divide_rect(
                    opt.point,
                    f_min,
                    x_next,
                    self.aq_func,
                    self.scaler
                )
                for bucket in self.buckets:
                    if bucket.diff_exist(opt.point.d):
                        bucket.delete()
                        if not bucket.array:
                            index = self.buckets.index(bucket)
                            del self.buckets[index]

            ei_min.append(f_min)
        return f_min, x_next

    def divide_rect(self, opt_rect, f_min, x_next, aq_func, scaler):
        """ divide the rectangular into smaller ones """
        rect = copy.deepcopy(opt_rect)
        division_num = rect.division_num
        j = np.mod(division_num, self.dim)
        k = (division_num - j) / self.dim
        max_side_len = np.power(3, float(-k))
        delta = max_side_len / 3
        dim_set = []
        for i in range(self.dim):
            if abs(max_side_len - (rect.u[0, i] - rect.l[0, i])) < 0.0000001:
                dim_set.append(i)

        dim_list = []
        for i in dim_set:
            e = np.zeros((1, self.dim))
            e[0, i] = 1
            function_value = min(
                aq_func.compute(scaler.inverse_transform(rect.center + delta * e)),
                aq_func.compute(scaler.inverse_transform(rect.center - delta * e))
            )
            dim_list.append(DimPack(i, function_value))
        dim_list.sort(key=lambda x: x.fc)

        for i in range(len(dim_list)):
            division_num = division_num + 1
            temp = np.zeros((1, self.dim))
            temp[0, dim_list[i].dim] = delta
            left_rect = RectPack(
                rect.l,
                rect.u - 2 * temp,
                division_num,
                self.dim,
                self.scaler,
                aq_func
            )
            middle_rect = RectPack(
                rect.l + temp,
                rect.u - temp,
                division_num,
                self.dim,
                self.scaler,
                aq_func
            )
            right_rect = RectPack(
                rect.l + 2 * temp,
                rect.u,
                division_num,
                self.dim,
                self.scaler,
                aq_func
            )
            if left_rect.fc < f_min:
                f_min = left_rect.fc
                x_next = left_rect.center
            if right_rect.fc < f_min:
                f_min = right_rect.fc
                x_next = right_rect.center

            insert = 0
            for bucket in self.buckets:
                if bucket.diff_exist(left_rect.d):
                    bucket.insert(left_rect)
                    bucket.insert(right_rect)
                    if i == len(dim_list) - 1:
                        bucket.insert(middle_rect)
                    insert = 1
                    break
            if insert == 0:
                new_bucket = RectBucket(left_rect.d, left_rect)
                new_bucket.insert(right_rect)
                if i == len(dim_list) - 1:
                    new_bucket.insert(middle_rect)
                self.buckets.append(new_bucket)
            rect = middle_rect
        return f_min, x_next
