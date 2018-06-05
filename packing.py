import itertools
import json
import os
import warnings

import binpacking
import numpy as np
import pandas as pd
from keras.models import load_model


def load_data(path):
    if not os.path.isfile(path):
        return []
    with open(path, "r") as json_data:
        statelist = json.load(json_data)
        if isinstance(statelist, list):
            return statelist
        return []


def pack(grids, tasks, maxcores, machines):
    ans = []
    time = 9999999
    tmplist = [i for i in range(1, maxcores + 1)]
    combinations = list(itertools.product(tmplist, repeat=len(tasks)))
    ansbins = []
    for i in combinations:
        rasklad, time, bins = calc(tasks=tasks, ixlst=list(i), summin=time, cores=maxcores, machines=machines,
                                   grids=grids)
        if len(bins) > 0:
            ansbins = bins
        if len(rasklad) > 0:
            ans = rasklad
    return ans, time, ansbins


def calc(tasks, ixlst, summin, cores, machines, grids):
    dfc = {}
    for i in tasks:
        dfc[i] = ixlst[tasks.index(i)]
    bins = binpacking.to_constant_volume(dfc, cores)
    if len(bins) <= machines:
        sumlst = []
        for b in bins:
            summachn = 0
            for key, item in b.items():
                t = grids.iat[tasks.index(key), item - 1]
                summachn = max(summachn, t)
            sumlst.append(summachn)
        tmp = max(sumlst)
        if tmp <= summin:
            return ixlst, tmp, bins
    return [], summin, []


def generate_grid(maxcores, tasks, df):
    grids = pd.DataFrame(columns=[n for n in range(1, maxcores + 1)])

    for i in range(len(tasks)):
        grids.loc[i] = [np.random.randint(-1, 1) for n in range(maxcores)]
    for ix, rowX in df.iterrows():
        grids.iloc[tasks.index(rowX.paramsvector)][rowX.cores] = rowX.time
    return grids


def generate_dataframe(maxcores, tasks, model):
    X_test = pd.DataFrame(columns=["cores", "paramsvector"])

    k = 0
    for i in tasks:
        for j in range(1, maxcores + 1):
            X_test.loc[k] = {"cores": j, "paramsvector": i}
            k += 1

    y_test = model.predict(X_test)
    X_test['time'] = y_test
    return X_test


def save_data(binslist, path="temp/out1.json", ):
    with open(path, "w") as f:
        json.dump(binslist, f)


def preprocess_list(ls1):
    outlst = []
    for mydict in ls1:
        dict1 = {}
        for key in mydict.keys():
            if type(key) is not str:
                dict1[str(key)] = int(mydict[key])
        outlst.append(dict1)
    return outlst


if __name__ == '__main__':
    warnings.filterwarnings("ignore")
    model = load_model('models/sleep.h5')
    intasks = load_data("temp/intasks1.json")
    global inmachines
    with open("temp/machines1.json", "r") as json_data:
        inmachines = json.load(json_data)
    global inmaxcores
    with open("temp/maxcores1.json", "r") as json_data:
        inmaxcores = json.load(json_data)
    df = generate_dataframe(maxcores=inmaxcores, tasks=intasks, model=model)
    ingrids = generate_grid(df=df, tasks=intasks, maxcores=inmaxcores)
    answer = (pack(grids=ingrids, tasks=intasks, maxcores=inmaxcores, machines=inmachines))
    save_data(preprocess_list(answer[2]))
