import torch
import yaml

data=dict()
is_available=torch.cuda.is_available()
data['cuda_available']=is_available

if is_available:
    data['device_count']=torch.cuda.device_count()
    data['device_name']=torch.cuda.get_device_name(0)

print('---')
print('')
print(yaml.dump(data))
