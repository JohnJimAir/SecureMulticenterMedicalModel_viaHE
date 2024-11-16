import numpy as np

# 读取 .npz 文件
data = np.load('syn3-test_data_x_y.npz')

# 获取所有数组的名称
keys = data.files

# 打开一个 .txt 文件用于写入
with open('syn3-test_data_x_y.txt', 'w') as f:
    for key in keys:
        # 获取每个数组的数据
        array_data = data[key]
        # 写入数组名称
        f.write(f'Array name: {key}\n')
        # 写入数组数据
        np.savetxt(f, array_data, fmt='%s')
        f.write('\n')  # 每个数组之间添加一个换行符
